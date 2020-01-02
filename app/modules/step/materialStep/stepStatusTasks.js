import { all, call, put, race, take } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import { orderActions } from '../../order/action';
import { ioDirection, ioTriggerMode } from '../../device/io/constants';
import actions, { MATERIAL_STEP } from './action';
import { CommonLog } from '../../../common/utils';
import dialogActions from '../../dialog/action';
import type { IDevice } from '../../device/IDevice';
import { getDevice } from '../../deviceManager/devices';
import ClsIOModule from '../../device/io/ClsIOModule';

function* enteringState(config) {
  try {
    const { parent } = config;
    const { workcenter: { locations } } = parent.payload;
    const material = this._consumeProduct;
    if (!material) {
      throw new Error(`consumed product not provided`);
    }
    const location = locations.find(l => l.product_code === material);
    if (!location) {
      throw new Error(`location not provided`);
    }
    if (!(location.equipment_sn)) {
      throw new Error(`io module not provided`);
    }
    const ioModule: IDevice = getDevice(location.equipment_sn);
    if (!(ioModule instanceof ClsIOModule)) {
      throw new Error(`io module (${location.equipment_sn}) not found`);
    }
    const input = ioModule.getPort(ioDirection.input, location.io_input);
    const output = ioModule.getPort(ioDirection.input, location.io_output);
    this._io.add(ioModule);
    this._ports.add(input);
    this._items.add({ io: ioModule, in: input, out: output });
    yield all(
      [...this._io].map(io => {
        if (!io?.ioContact) {
          throw new Error(`io invalid ${io?.sn}`);
        }
        return call(io.ioContact);
      })
    );

    const sPayload = this._payload;
    const confirmSN = sPayload?.confirm?.sn;
    if (confirmSN) {
      const confirmIO = getDevice(confirmSN);
      if (confirmIO && confirmIO instanceof ClsIOModule) {
        this._confirm = {
          io: confirmIO,
          port: confirmIO.getPort(ioDirection.input, sPayload.confirm.index)
        };
      }
    }

    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
  } catch (e) {
    CommonLog.lError(e);
    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING, { msg: e.message }));
  }
}

function* doingState() {
  try {
    const listeners = [];
    [...this._items].forEach(i => {
      listeners.push({
        listener: i.io.addListener(
          input => i.in === input.data.port &&
            ioTriggerMode.falling === input.data.triggerMode,
          () => actions.item(i)
        ),
        io: i.io
      });
    });
    yield all(
      [...this._items]
        .map(item => {
          if (item?.io?.openIO && item?.out) {
            return call(item.io.openIO, item.out);
          }
          return null;
        })
        .filter(calls => !!calls)
    );

    let readyListener = null;
    // const confirmPort = io.getPort(ioDirection.input, confirmIdx(sPayload));
    if (this._confirm && this._confirm.io && this._confirm.port) {
      readyListener = this._confirm.io.addListener(
        input =>
          this._confirm.port === input.port &&
          ioTriggerMode.falling === input.triggerMode,
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
}

function* failState(config) {
  try {
    const { msg } = config;
    yield put(
      dialogActions.dialogShow({
        maxWidth: 'md',
        buttons: [
          {
            label: 'Common.Close',
            color: 'danger'
          },
          {
            label: '重试',
            color: 'info',
            action: orderActions.doPreviousStep()
          },
          {
            label: 'Order.Next',
            color: 'warning',
            action: orderActions.finishStep(this)
          }
        ],
        title: `工步失败：${this._name}`,
        content: `${msg || ''}`
      })
    );
    // yield put(orderActions.finishStep(this));
  } catch (e) {
    CommonLog.lError(e);
  }
}

export const materialStepStatusMixin = (superTasks) => ({
  ...superTasks,
  [STEP_STATUS.DOING]: doingState,
  [STEP_STATUS.FAIL]: failState,
  [STEP_STATUS.ENTERING]: enteringState
});

export function* onLeave(): Saga<void> {
  try {
    yield all(
      [...this._io].map(io =>
        call(io.closeIO, [...this._ports].filter(p => io.hasPort(p)))
      )
    );
    this._ports.clear();
    this._io.clear();
    this._items.clear();
    this._confirm = null;
  } catch (e) {
    CommonLog.lError(`MaterialStepMixin onLeave Error: ${e.toString()}`);
  }
}
