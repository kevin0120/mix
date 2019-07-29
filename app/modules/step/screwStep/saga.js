import { put, select } from 'redux-saga/effects';
import { STEP_STATUS } from '../model';
import { stepPayload, processingStep } from '../../order/selector';


export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      const payload = yield select(s => stepPayload(processingStep(s.order)));
      yield put(orderActions.stepData((data) => ({
        points: JSON.parse(JSON.stringify(payload?.points || [])),
        activePoint:0,
        ...data,
      })));
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      console.error(e);
    }
  }


};
