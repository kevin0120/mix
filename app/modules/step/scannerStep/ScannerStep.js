// @flow
import { call, put, select, take, all } from 'redux-saga/effects';
import {STEP_STATUS} from '../constants';
import {
  stepData,
  stepPayload,
  workingOrder,
  workingStep
} from '../../order/selector';
import { SCANNER_STEP, scannerStepAction } from './action';
import { deviceType, getDevicesByType } from '../../external/device';
import { CommonLog } from '../../../common/utils';
import type { IWorkStep } from '../interface/IWorkStep';


const ScannerStepMixin = (ClsBaseStep: Class<IWorkStep>) => class ClsScannerStep extends ClsBaseStep {

  _scanners = [];

  _onLeave = () => {
    this._scanners = [];
    console.log('scanners cleared');
  };

  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        this._scanners = getDevicesByType(deviceType.scanner);
        const scanners = this._scanners;
        yield all(scanners.map((s) => {
          // eslint-disable-next-line no-param-reassign
          s.dispatcher = scannerStepAction.getValue;
          return call(s.Enable);
        }));
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
                ...(d || {}),
                result: {
                  [label]: action?.input?.data
                },
                timeLine: [
                  {
                    title: action?.input?.name,
                    color: 'info',
                    footerTitle: action && action.input && action.input.time.toLocaleString(),
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
        yield all(this._scanners.map((s) => {
          s.dispatcher = null;
          return call(s.Disable);
        }));
      }
    },
    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
        this._scanners = [];
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  };
};
export default ScannerStepMixin;
