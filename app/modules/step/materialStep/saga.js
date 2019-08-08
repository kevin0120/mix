import { put, take, call } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import actions, { MATERIAL_STEP } from './action';
import { defaultIO } from '../../io/saga';
import { ioTriggerMode, ioDirection } from '../../io/type';

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      yield call([defaultIO, defaultIO.ioContact]);
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      const port = defaultIO.getPort(ioDirection.input, 4);
      const ioListener = defaultIO.addListener(
        port,
        ioTriggerMode.falling,
        actions.ioInput
      );
      // while(true){
        const input = yield take(MATERIAL_STEP.IO_INPUT);
        console.log('get input!!',input);
      // }
      defaultIO.removeListener(ioListener);
      yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
    } catch (e) {
      console.error(e);
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
