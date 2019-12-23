import { STEP_STATUS } from '../constants';
import { all, call, put, race, select, take } from 'redux-saga/effects';
import { orderActions } from '../../order/action';
import { ioDirection, ioTriggerMode } from '../../device/io/constants';
import actions, { MATERIAL_STEP } from './action';
import { CommonLog } from '../../../common/utils';
import dialogActions from '../../dialog/action';
import { stepPayload, workingOrder, workingStep } from '../../order/selector';
import type { IDevice } from '../../device/IDevice';
import { getDevice } from '../../deviceManager/devices';
import ClsIOModule from '../../device/io/ClsIOModule';

const items = payload => payload?.items;

function* doingState() {
  try {
    const listeners = [];
    [...this._items].forEach(i => {
      listeners.push({
        listener: i.in.io.addListener(
          input =>
            i.in.port === input.port &&
            ioTriggerMode.falling === input.triggerMode,
          () => actions.item(i)
        ),
        io: i.in.io
      });
    });

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

function* failState(msg) {
  try {
    yield put(
      dialogActions.dialogShow({
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

function* enteringState() {
  try {
    const sPayload = yield select(s =>
      stepPayload(workingStep(workingOrder(s.order)))
    );
    items(sPayload).forEach(i => {
      let item = null;
      ['in', 'out'].forEach(dir => {
        if (!(i && i[dir] && i[dir].sn)) {
          throw new Error(`io module ${i[dir].sn} not provided`);
        }
        const io: IDevice = getDevice(i[dir].sn);
        if (!(io instanceof ClsIOModule)) {
          throw new Error(`io module ${i[dir].sn} not found`);
        }
        this._io.add(io);
        const port = io.getPort(ioDirection.input, i[dir].index);
        this._ports.add(port);
        item = { io, port };
      });
      if (!item) {
        return;
      }
      this._items.add(item);
    });

    yield all(
      [...this._io].map(io => {
        if (!io?.ioContact) {
          throw new Error(`io invalid ${io?.sn}`);
        }
        return call(io.ioContact);
      })
    );

    yield all(
      [...this._items]
        .map(item => {
          if (item?.out?.io?.openIO && item?.out?.port) {
            return call(item.out.io.openIO, item.out.port);
          }
          return null;
        })
        .filter(calls => !!calls)
    );

    const confirmIO = getDevice(sPayload.confirm.sn);
    if (confirmIO && confirmIO instanceof ClsIOModule) {
      this._confirm = {
        io: confirmIO,
        port: confirmIO.getPort(ioDirection.input, sPayload.confirm.index)
      };
    }

    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
  } catch (e) {
    CommonLog.lError(e);
    yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e.message));
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
