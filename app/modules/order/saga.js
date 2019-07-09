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
      if (!currentProcessingStep || !currentProcessingStep.type) {
        throw(new Error('invalid step'));
      }else if(!currentProcessingStep.type){
        throw(new Error('step has no type'));
      }else if(!stepTypes[currentProcessingStep.type]){
        throw(new Error(`step type has no handler: ${currentProcessingStep.type}`));
      }else {
        yield race([call(stepTypes[currentProcessingStep.type], ORDER, orderActions), take(ORDER.STEP.PUSH)]);
      }
    }
  } catch (e) {
    console.error(e);
  }
}


