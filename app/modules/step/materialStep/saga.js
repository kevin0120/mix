import { put, take, call, select } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import actions, { MATERIAL_STEP } from './action';
import { defaultIO } from '../../io/saga';
import { ioTriggerMode, ioDirection } from '../../io/type';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { CommonLog } from '../../../common/utils';

const ports = [];

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      yield call(defaultIO.ioContact);
      // yield call(defaultIO.getStatus);
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const { items } = sPayload;
      items.forEach((i) => {
        ports.push(defaultIO.getPort(ioDirection.output, i.index));
      });
      yield call(defaultIO.openIO, ports);
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const inputPort = defaultIO.getPort(ioDirection.input, sPayload?.confirmIO);
      const ioListener = defaultIO.addListener(
        inputPort,
        ioTriggerMode.falling,
        actions.ioInput
      );
      yield take(MATERIAL_STEP.READY);
      defaultIO.removeListener(ioListener);
      yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      yield call(defaultIO.closeIO, ports);
      yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};
