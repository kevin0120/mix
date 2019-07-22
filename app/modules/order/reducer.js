import { ORDER } from './action';
import { genReducers } from '../indexReducer';
import { ORDER_STEP_STATUS, ORDER_STATUS } from './model';
import { demoOrder, demoOrder2 } from './demoData';
import {
  viewingStep,
  processingStep,
  processingIndex,
  currentOrderLength,
  viewingIndex,
  currentOrder
} from './selector';

const initState = {
  currentOrder: null,
  processingIndex: 0,
  viewingIndex: 0,
  list: [demoOrder, demoOrder2]
};

function setStepData(state, reducer) {
  const newState = { ...state };
  const newStep = processingStep(newState);
  newStep.data = reducer(newStep);
  return newState;
}

function setOrderStatus(state, status){
  const newState = { ...state };
  const newOrder=currentOrder(newState);
  newOrder.status=status;
  return newState;
}


function setStepStatus(state, status) {
  const newState = {
    ...state
  };
  const newStep = processingStep(newState);
  newStep.status = status;
  return newState;
}

function limitIndex(state, index) {
  if (index < 0) {
    return 0;
  }
  if (index >= currentOrderLength(state)) {
    return currentOrderLength(state) - 1;
  }
  return index;
}

const orderReducer = {
  [ORDER.TRIGGER]: (state, action) => {
    const firstIndex = 0;
    return {
      ...state,
      currentOrder: action.order || null,
      viewingIndex: firstIndex,
      processingIndex: firstIndex,
      startTime: new Date()
    };
  },
  [ORDER.FINISH]: (state) => setOrderStatus(state,ORDER_STATUS.DONE),
  [ORDER.FAIL]: (state) => setOrderStatus(state,ORDER_STATUS.FAIL),
  [ORDER.STEP.NEXT]: (state) => {
    const newIndex = limitIndex(state, viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.PREVIOUS]: (state) => {
    const newIndex = limitIndex(state, viewingIndex(state) - 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.JUMP_TO]: (state, action) => {
    const { stepId } = action;
    return {
      ...state,
      viewingIndex: stepId
    };
  },

  // 修改step的状态
  [ORDER.STEP.ENTER]: (state) => setStepStatus(state, ORDER_STEP_STATUS.ENTERING),
  [ORDER.STEP.ENTERED]: (state) => setStepStatus(state, ORDER_STEP_STATUS.DOING),
  [ORDER.STEP.LEAVE]: (state) => setStepStatus(state, ORDER_STEP_STATUS.LEAVING),
  [ORDER.STEP.FINISH]: (state) => setStepStatus(state, ORDER_STEP_STATUS.FINISHED),
  [ORDER.STEP.FAIL]: (state) => setStepStatus(state, ORDER_STEP_STATUS.FAIL),
  [ORDER.STEP.RESET]: (state) => setStepStatus(state, ORDER_STEP_STATUS.READY),
  //
  [ORDER.STEP.DO_NEXT]: (state) => {
    const newIndex = processingIndex(state) + 1;
    return {
      ...state,
      processingIndex: newIndex,
      viewingIndex: processingStep(state) === viewingStep(state) ? newIndex : viewingIndex(state)
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: (state) => {
    const newIndex = limitIndex(state, processingIndex(state) - 1);
    const revokedState = {
      ...state,
      processingIndex: newIndex,
      viewingIndex: processingStep(state) === viewingStep(state) ? newIndex : viewingIndex(state)
    };
    return setStepStatus(revokedState, ORDER_STEP_STATUS.READY);
  },
  [ORDER.STEP.DATA]: (state, action) => {
    const { reducer } = action;
    return setStepData(state, reducer);
  }
};

export default genReducers(orderReducer, initState);
