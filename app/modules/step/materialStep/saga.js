import { put, take, call, select } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import actions, { MATERIAL_STEP } from './action';
import { ioTriggerMode, ioDirection } from '../../external/device/io/type';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { CommonLog } from '../../../common/utils';
import { getDevice } from '../../external/device';

const ports = [];

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const { items, ioSN } = sPayload;
      const io = getDevice(ioSN);
      if(io?.ioContact){
        yield call(io.ioContact);
      }else{
        CommonLog.lError('io not ready',{
          at:'materialStep entering',
          io
        });
        yield put(orderActions.stepStatus(STEP_STATUS.FAIL));
        return;
      }
      items.forEach((i) => {
        ports.push(io.getPort(ioDirection.output, i.index));
      });
      yield call(io.openIO, ports);
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const { ioSN } = sPayload;
      const io = getDevice(ioSN);
      const inputPort = io.getPort(ioDirection.input, sPayload?.confirmIO);
      const ioListener = io.addListener(
        inputPort,
        ioTriggerMode.falling,
        actions.ioInput
      );
      yield take(MATERIAL_STEP.READY);
      io.removeListener(ioListener);
      yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const { ioSN } = sPayload;
      const io = getDevice(ioSN);
      yield call(io.closeIO, ports);
      yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FAIL](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const { ioSN } = sPayload;
      const io = getDevice(ioSN);
      if(io?.closeIO){
        yield call(io.closeIO, ports);
      }
      // yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};
