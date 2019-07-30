import { ORDER } from './action';
import { genReducers } from '../util';
import { ORDER_STATUS } from './model';
import { STEP_STATUS } from '../step/model';
import { demoOrder, demoOrder2, demoOrderExcp } from './demoData';
import {
  viewingStep,
  processingStep,
  processingIndex,
  orderLength,
  viewingIndex,
  currentOrder, getStep
} from './selector';
import type { tOrderState, tOrderStepIdx } from './model';

const initState = {
  currentOrder: null,
  processingIndex: 0,
  viewingIndex: 0,
  list: [demoOrder, demoOrderExcp]
};

function reduceStepData(reducer, state: tOrderState, action): tOrderState {
  const newState = { ...state };
  const newStep = processingStep(newState);
  newStep.data = reducer(newStep.data, action);
  return newState;
}

function setOrderStatus(state: tOrderState, status): tOrderState {
  const newState = { ...state };
  const newOrder = currentOrder(newState);
  newOrder.status = status;
  return newState;
}

function setStepStatus(state: tOrderState, action): tOrderState {
  const newState = {
    ...state
  };
  const newStep = processingStep(newState);
  newStep.status = action.status;
  return newState;
}

function limitIndex(state: tOrderState, index: tOrderStepIdx): tOrderStepIdx {
  if (index < 0) {
    return 0;
  }
  if (index >= orderLength(state)) {
    return orderLength(state) - 1;
  }
  return index;
}

const orderReducer = {
  [ORDER.SWITCH]: (state: tOrderState, action): tOrderState => {
    const firstIndex = 0;
    return {
      ...state,
      currentOrder: action.order || null,
      viewingIndex: firstIndex,
      processingIndex: firstIndex,
      startTime: new Date()
    };
  },
  [ORDER.FINISH]: (state: tOrderState): tOrderState => {
    const newState = setOrderStatus(state, ORDER_STATUS.DONE);
    return {
      ...newState,
      // currentOrder: null
    };
  },
  [ORDER.CANCEL]: (state: tOrderState): tOrderState => {
    const newState = setOrderStatus(state, ORDER_STATUS.CANCEL);
    return {
      ...newState,
      // currentOrder: null
    };
  },
  [ORDER.PENDING]: (state: tOrderState): tOrderState => {
    const newState = setOrderStatus(state, ORDER_STATUS.PENDING);
    return {
      ...newState,
      // currentOrder: null
    };
  },
  [ORDER.STEP.NEXT]: (state: tOrderState): tOrderState => {
    const newIndex = limitIndex(state, viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.PREVIOUS]: (state: tOrderState): tOrderState => {
    const newIndex = limitIndex(state, viewingIndex(state) - 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.JUMP_TO]: (state: tOrderState, action): tOrderState => {
    const { stepId } = action;
    return {
      ...state,
      viewingIndex: stepId
    };
  },

  // 修改step的状态
  [ORDER.STEP.STATUS]: setStepStatus,
  //
  [ORDER.STEP.DO_NEXT]: (state: tOrderState): tOrderState => {
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
  [ORDER.STEP.DO_PREVIOUS]: (state: tOrderState): tOrderState => {
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
  [ORDER.STEP.DATA]: (state: tOrderState, action): tOrderState => {
    const { reducer } = action;
    return reduceStepData(reducer, state, action);
  },
  [ORDER.STEP.START_TIME]: (state: tOrderState, action): tOrderState => {
    const { idx, startTime } = action;
    const newState = { ...state };
    const newStep = getStep(newState, idx);
    newStep.startTime = startTime;
    return newState;
  },
  [ORDER.STEP.END_TIME]: (state: tOrderState, action): tOrderState => {
    const { idx, endTime } = action;
    const newState = { ...state };
    const newStep = getStep(newState, idx);
    newStep.endTime = endTime;
    return newState;
  }
};

export default genReducers(orderReducer, initState);
