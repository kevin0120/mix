import Step from '../Step';
import STEP_STATUS from '../model';
import { call, put, select, take } from 'redux-saga/effects';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { getDevice } from '../../external/device';
import { CommonLog } from '../../../common/utils';
import { ioDirection, ioTriggerMode } from '../../external/device/io/type';
import actions, { MATERIAL_STEP } from './action';

const ioSN = (payload) => payload?.ioSN;
const items = (payload) => payload?.items;
const confirmIdx = (payload) => payload?.confirmIO;


export default class MaterialStep extends Step {
  _ports = [];

  constructor(...args) {
    super(...args);
    function* onLeave() {
      try {
        const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
        const io = getDevice(ioSN(sPayload));
        if (io?.closeIO) {
          yield call(io.closeIO, this._ports);
        }
        this._ports = [];
        console.log('ports cleared');
      } catch (e) {
        CommonLog(e);
      }
    }
    this._onLeave = onLeave.bind(this);
  }

  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        const sPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
        const io = getDevice(ioSN(sPayload));
        if (io?.ioContact) {
          yield call(io.ioContact);
        } else {
          CommonLog.lError('io not ready', {
            at: 'materialStep entering',
            io
          });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
          return;
        }
        items(sPayload).forEach((i) => {
          this._ports.push(io.getPort(ioDirection.output, i.index));
        });
        yield call(io.openIO, this._ports);
        // console.log(this);
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
        yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
      } catch (e) {
        CommonLog.lError(e);
      } finally {

      }
    },
    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    * [STEP_STATUS.FAIL](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  };
}
