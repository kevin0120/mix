import { take, call, race, fork, select, put } from 'redux-saga/effects';
import { ORDER_STEP_STATUS } from './model';
import { ORDER, orderActions } from './action';
import stepTypes from '../stepTypes';

export default function* root() {
  try {
    while (true) {
      // take trigger action
      const action = yield take(ORDER.TRIGGER);
      // do order
      yield call(doOrder);
      console.log('order finished');
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
          finish: call(stepTypes[currentProcessingStep.type], ORDER, orderActions),
          next: take(ORDER.STEP.PUSH)
        });
        console.log(next);
        if (next) {
          const newState = yield select();
          const { currentProcessingIndex, currentOrder } = newState.order;
          if (currentProcessingIndex >= currentOrder.steps.length) {
            return;
          }
        }
      }
    }
  } catch (e) {
    console.error(e);
  }
}


