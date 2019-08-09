import { take, put, select, call } from 'redux-saga/effects';
import { SCANNER_STEP } from './action';
import {
  workingStep,
  stepData,
  stepPayload,
  workingOrder
} from '../../order/selector';
import STEP_STATUS from '../model';
import { scanner } from '../../scanner/saga';

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      yield call(scanner.Enable);
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
            yield put(
              orderActions.stepData(d => ({
                ...d,
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
                  ...(d.timeLine || [])
                ]
              }))
            );
            break;
          case SCANNER_STEP.SUBMIT:
            if (Object.hasOwnProperty.call(result || {}, label)) {
              yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
            }
            break;
          default:
            break;
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      yield call( scanner.Disable);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      console.error(e);
    }
  }
};
