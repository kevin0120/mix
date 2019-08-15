import { put, take, call, select } from 'redux-saga/effects';
import STEP_STATUS from '../model';
import actions, { MATERIAL_STEP } from './action';
import { ioTriggerMode, ioDirection } from '../../external/device/io/type';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { CommonLog } from '../../../common/utils';
import { getDevice } from '../../external/device';

const ports = [];

const ioSN = (payload) => payload?.device?.IO?.sn;
const items = (payload) => payload?.items;
const confirmIdx=(payload)=>payload?.device?.IO?.confirmIdx;

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const io = getDevice(ioSN(sPayload));
      if(io?.ioContact) {
        yield call(io.ioContact);
      } else {
        CommonLog.lError('io not ready', {
          at: 'materialStep entering',
          io
        });
        yield put(orderActions.stepStatus(STEP_STATUS.FAIL));
        return;
      }
      items(sPayload).forEach((i) => {
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
      const io = getDevice(ioSN(sPayload));
      const confirmPort = io.getPort(ioDirection.input, confirmIdx(sPayload));
      const ioListener = io.addListener(
        confirmPort,
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
      const io = getDevice(ioSN(sPayload));
      yield call(io.closeIO, ports);
      yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FAIL](ORDER, orderActions) {
    try {
      const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      const io = getDevice(ioSN(sPayload));
      if (io?.closeIO) {
        yield call(io.closeIO, ports);
      }
      // yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};
