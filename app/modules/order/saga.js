import { take, call, race, fork, select, put, join } from 'redux-saga/effects';
import { push } from 'connected-react-router';
import { ORDER, orderActions } from './action';
import steps from '../step/saga';
import {
  processingStep,
  processingIndex,
  stepType,
  orderLength
} from './selector';
import dialogActions from '../dialog/action';

const mapping = {
  onOrderFinish: showResult
};

function* showResult() {
  try {
    yield put(
      dialogActions.showDialog({
        hasOk: true,
        closeAction: push('/app')
      })
    );
  } catch (e) {
    console.error('returnHome error', e);
  }
}

export default function* root() {
  try {
    while (true) {
      const {order}=yield take(ORDER.TRIGGER);
      const didSwitchOrder = yield fork(function* () {
        yield take(ORDER.SWITCH);
      });
      yield put(orderActions.switchOrder(order));
      yield join(didSwitchOrder);
      const { finish } = yield race({
        exit: call(doOrder),
        finish: take(ORDER.FINISH),
        fail: take(ORDER.FAIL)
      });
      console.log('order finished');
      if (finish) {
        yield call(mapping.onOrderFinish);
      }
    }
  } catch (e) {
    console.error(e);
  }
}

function* doOrder() {
  try {
    while (true) {
      console.log('doing order');
      const type = yield select(state => stepType(processingStep(state.order)));
      const { next } = yield race({
        exit: call(steps, type, ORDER, orderActions),
        next: take(ORDER.STEP.DO_NEXT),
        previous: take(ORDER.STEP.DO_PREVIOUS)
      });

      if (next) {
        const order = yield select(state => state.order);
        if (processingIndex(order) >= orderLength(order)) {
          yield put(orderActions.finishOrder());
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}
