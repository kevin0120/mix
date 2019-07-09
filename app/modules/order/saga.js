import { take, call, race, fork, select, put } from 'redux-saga/effects';
import { ORDER_STEP_STATUS } from './model';
import { ORDER, orderActions } from './action';


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
      const action = yield take(ORDER.STEP.TRY_PUSH);
      const state = yield select();
      const { currentOrder, currentProcessingIndex, currentProcessingStep } = state.order;
      if(currentProcessingStep instanceof Array){
        console.log(action.parallelId);

      }
      switch (currentProcessingStep.status) {
        case ORDER_STEP_STATUS.READY:
          yield put(orderActions.enterStep());
          break;
        case ORDER_STEP_STATUS.ENTERING:
          yield put(orderActions.enteredStep());
          break;
        case ORDER_STEP_STATUS.DOING:
          yield put(orderActions.leaveStep());
          break;
        case ORDER_STEP_STATUS.LEAVING:
          yield put(orderActions.finishStep());
          break;
        case ORDER_STEP_STATUS.FINISHED:
          yield put(orderActions.pushStep());
          break;
        case ORDER_STEP_STATUS.FAIL:

          break;
        default:
          yield put(orderActions.enterStep());
          break;
      }
    }
  } catch (e) {
    console.error(e);
  }
}


