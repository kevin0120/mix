import Step from '../Step';
import STEP_STATUS from '../model';
import { call, put, select, take } from '@redux-saga/core/effects';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { getDevice } from '../../external/device';
import { CommonLog } from '../../../common/utils';
import { ioDirection, ioTriggerMode } from '../../external/device/io/type';
import actions, { MATERIAL_STEP } from './action';

const ioSN = (payload) => payload?.device?.IO?.sn;
const items = (payload) => payload?.items;
const confirmIdx=(payload)=>payload?.device?.IO?.confirmIdx;


export default class MaterialStep extends Step {
  _ports=[];
  _statusTasks = {
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
          yield put(orderActions.stepStatus(this,STEP_STATUS.FAIL));
          return;
        }
        items(sPayload).forEach((i) => {
          this._ports.push(io.getPort(ioDirection.output, i.index));
        });
        yield call(io.openIO, this._ports);
        console.log(this);
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
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
        yield put(orderActions.stepStatus(this,STEP_STATUS.FINISHED));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
        const io = getDevice(ioSN(sPayload));
        yield call(io.closeIO, this._ports);
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    * [STEP_STATUS.FAIL](ORDER, orderActions) {
      try {
        const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
        const io = getDevice(ioSN(sPayload));
        if (io?.closeIO) {
          yield call(io.closeIO, this._ports);
        }
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  }
}