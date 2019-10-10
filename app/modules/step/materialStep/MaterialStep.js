import { call, put, select, take, all, race } from 'redux-saga/effects';
import STEP_STATUS from '../constents';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import { getDevice } from '../../external/device';
import { CommonLog } from '../../../common/utils';
import { ioDirection, ioTriggerMode } from '../../external/device/io/type';
import actions, { MATERIAL_STEP } from './action';

const items = payload => payload?.items;

const MaterialStepMixin = (ClsBaseStep) => class ClsMaterialStep extends ClsBaseStep {
  _ports = new Set([]);

  _io = new Set([]);

  _items = new Set([]);

  _confirm = null;

  constructor(...args) {
    super(...args);

    function* onLeave() {
      try {
        for (const io of this._io) {
          yield call(io.closeIO, [...this._ports].filter(p => io.hasPort(p)));
        }
        this._ports.clear();
        this._io.clear();
        this._items.clear();
        this._confirm = null;
      } catch (e) {
        CommonLog(e);
      }
    }

    this._onLeave = onLeave.bind(this);
  }

  _statusTasks = {
    *[STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        const sPayload = yield select(s =>
          stepPayload(workingStep(workingOrder(s.order)))
        );
        items(sPayload).forEach(i => {
          const item = {};
          if (i?.in?.sn) {
            const io = getDevice(i.in.sn);
            this._io.add(io);
            const port = io.getPort(ioDirection.input, i.in.index);
            this._ports.add(port);
            item.in = { io, port };
          }
          if (i?.out?.sn) {
            const io = getDevice(i.out.sn);
            this._io.add(io);
            const port = io.getPort(ioDirection.output, i.out.index);
            this._ports.add(port);
            item.out = { io, port };
          }
          this._items.add(item);
        });
        for (const io of this._io) {
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
        }

        for (const item of this._items) {
          if (item?.out?.io?.openIO && item?.out?.port) {
            yield call(item.out.io.openIO, item.out.port);
          }
        }

        const confirmIO = getDevice(sPayload.confirm.sn);
        if (confirmIO) {
          this._confirm = {
            io: confirmIO,
            port: confirmIO.getPort(ioDirection.input, sPayload.confirm.index)
          };
        }

        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    *[STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        const listeners = [];
        this._items.forEach(i => {
          listeners.push({
            listener: i.in.io.addListener(
              i.in.port,
              ioTriggerMode.falling,
              () => actions.item(i)
            ),
            io: i.in.io
          });
        });

        let readyListener = null;
        // const confirmPort = io.getPort(ioDirection.input, confirmIdx(sPayload));
        if (this._confirm && this._confirm.io && this._confirm.port) {
          readyListener = this._confirm.io.addListener(
            this._confirm.port,
            ioTriggerMode.falling,
            actions.ready
          );
        }

        yield race([
          take(MATERIAL_STEP.READY),
          all(
            [...this._items].map(i =>
              take(a => a.type === MATERIAL_STEP.ITEM && a.item === i)
            )
          )
        ]);

        if (readyListener) {
          this._confirm.io.removeListener(readyListener);
        }

        listeners.forEach(l => {
          l.io.removeListener(l.listener);
        });

        yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    *[STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    },
    *[STEP_STATUS.FAIL](ORDER, orderActions) {
      try {
        yield put(orderActions.finishStep(this));
      } catch (e) {
        CommonLog.lError(e);
      }
    }
  };
};
export default MaterialStepMixin;
