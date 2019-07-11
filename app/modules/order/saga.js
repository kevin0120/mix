import { take, call, race, fork, select, put } from 'redux-saga/effects';
import { ORDER_STEP_STATUS } from './model';
import { ORDER, orderActions } from './action';
import stepTypes from '../steps/stepTypes';
import { push } from 'connected-react-router';

const mapping = {
  onOrderFinish:returnHome,
};


function* returnHome(){
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
      const state = yield select();
      const { currentProcessingStep } = state.order;
      if (!currentProcessingStep) {
        throw(new Error('invalid step'));
      } else if (!currentProcessingStep.type) {
        throw(new Error('step has no type'));
      } else if (!stepTypes[currentProcessingStep.type]) {
        throw(new Error(`step type has no handler: ${currentProcessingStep.type}`));
      } else {
        const { next } = yield race({
          exit: call(stepTypes[currentProcessingStep.type], ORDER, orderActions),
          next: take(ORDER.STEP.PUSH)
        });
        if (next) {
          const newState = yield select();
          const { currentProcessingIndex, currentOrder } = newState.order;
          if (currentProcessingIndex >= currentOrder.steps.length) {
            yield put(orderActions.finishOrder());
          }
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}


