import { ORDER } from './action';
import { genReducers } from '../indexReducer';
import { STEP_STATUS, ORDER_STATUS } from './model';
import { demoOrder, demoOrder2, demoOrderExcp } from './demoData';
import {
  viewingStep,
  processingStep,
  processingIndex,
  orderLength,
  viewingIndex,
  currentOrder
} from './selector';

const initState = {
  currentOrder: null,
  processingIndex: 0,
  viewingIndex: 0,
  list: [demoOrder, demoOrder2, demoOrderExcp]
};

function reduceStepData(reducer, state, action) {
  const newState = { ...state };
  const newStep = processingStep(newState);
  newStep.data = reducer(newStep.data, action);
  return newState;
}

function setOrderStatus(state, status) {
  const newState = { ...state };
  const newOrder = currentOrder(newState);
  newOrder.status = status;
  return newState;
}

function setStepStatus(state, action) {
  const newState = {
    ...state
  };
  const newStep = processingStep(newState);
  newStep.status = action.status;
  return newState;
}

function limitIndex(state, index) {
  if (index < 0) {
    return 0;
  }
  if (index >= orderLength(state)) {
    return orderLength(state) - 1;
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
  [ORDER.FINISH]: state => {
    const newState = setOrderStatus(state, ORDER_STATUS.DONE);
    return {
      ...newState,
      currentOrder: null
    };
  },
  [ORDER.FAIL]: state => {
    const newState = setOrderStatus(state, ORDER_STATUS.FAIL);
    return {
      ...newState,
      currentOrder: null
    };
  },
  [ORDER.STEP.NEXT]: state => {
    const newIndex = limitIndex(state, viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.PREVIOUS]: state => {
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
  [ORDER.STEP.STATUS]: setStepStatus,
  //
  [ORDER.STEP.DO_NEXT]: state => {
    const newIndex = processingIndex(state) + 1;
    return {
      ...state,
      processingIndex: newIndex,
      viewingIndex:
        processingStep(state) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: state => {
    const newIndex = limitIndex(state, processingIndex(state) - 1);
    const revokedState = {
      ...state,
      processingIndex: newIndex,
      viewingIndex:
        processingStep(state) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
    return setStepStatus(revokedState, STEP_STATUS.READY);
  },
  [ORDER.STEP.DATA]: (state, action) => {
    const { reducer } = action;
    return reduceStepData(reducer, state, action);
  }
};

export default genReducers(orderReducer, initState);
