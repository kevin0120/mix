import { take, call, race, fork, select, put } from 'redux-saga/effects';
import { push } from 'connected-react-router';
import { ORDER, orderActions } from './action';
import stepTypes from '../steps/stepTypes';
import { processingStep, processingIndex, stepType, currentOrderLength } from './selector';

const mapping = {
  onOrderFinish: returnHome
};


function* returnHome() {
  try {
    yield put(push('/app'));
  } catch (e) {
    console.error('returnHome error', e);
  }
}

export default function* root() {
  try {
    while (true) {
      // take trigger action
      yield take(ORDER.TRIGGER);
      // do order
      const { exit, finish, fail } = yield race({
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
      const type = yield select((state) => stepType(processingStep(state.order)));
      if (!type) {
        throw new Error('step type not valid:', type);
      }
      const { next } = yield race({
        exit: call(stepTypes[type], ORDER, orderActions),
        next: take(ORDER.STEP.PUSH)
      });
      if (next) {
        const order = yield select((state) => state.order);
        if (processingIndex(order) >= currentOrderLength(order)) {
          yield put(orderActions.finishOrder());
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}


