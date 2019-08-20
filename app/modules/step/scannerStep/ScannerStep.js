import Step from '../Step';
import STEP_STATUS from '../model';
import { call, put, select, take } from 'redux-saga/effects';
import { stepData, stepPayload, workingOrder, workingStep } from '../../order/selector';
import { SCANNER_STEP, scannerStepAction } from './action';
import { deviceType, getDevicesByType } from '../../external/device';
import { CommonLog } from '../../../common/utils';


export default class ScannerStep extends Step {
  _scanners = [];

  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        this._scanners = getDevicesByType(deviceType.scanner);
        for (const scanner of this._scanners) {
          yield call(scanner.Enable);
          scanner.dispatcher = scannerStepAction.getValue;
        }
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    * [STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        while (true) {
          const action = yield take([
            SCANNER_STEP.GET_VALUE,
            SCANNER_STEP.SUBMIT
          ]);
          const result = yield select(
            s => stepData(workingStep(workingOrder(s.order)))?.result
          );
          const label = yield select(
            s => stepPayload(workingStep(workingOrder(s.order)))?.label
          );
          switch (action.type) {
            case SCANNER_STEP.GET_VALUE:
              yield call(this.updateData, d => ({
                ...d || {},
                result: {
                  [label]: action?.input?.data
                },
                timeLine: [
                  {
                    title: action?.input?.name,
                    color: 'info',
                    footerTitle: action?.input?.time.toLocaleString(),
                    body: action?.input?.data
                  },
                  ...(d?.timeLine || [])
                ]
              }));
              break;
            case SCANNER_STEP.SUBMIT:
              // if(this._steps.length>0){
              //   yield call(this.runSubStep,this._steps[0])
              // }
              if (Object.hasOwnProperty.call(result || {}, label)) {
                yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
              }
              break;
            default:
              break;
          }
        }
      } catch (e) {
        CommonLog.lError(e);
      } finally {
        for (const scanner of this._scanners) {
          yield call(scanner.Disable);
          scanner.dispatcher = null;
        }
      }
    },
    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  };
}
