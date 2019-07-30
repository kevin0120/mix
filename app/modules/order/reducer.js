// @flow
import { ORDER } from './action';
import { genReducers } from '../util';
import { ORDER_STATUS } from './model';
import { STEP_STATUS } from '../step/model';
import { demoOrder, demoOrderPending, demoOrderCancel, demoOrderDone } from './demoData';
import {
  viewingStep,
  workingStep,
  workingIndex,
  orderLength,
  viewingIndex,
  viewingOrder,
  workingOrder,
  getStep
} from './selector';
import type { tOrder, tOrderState, tOrderStepIdx, tStep } from './model';

const initState = {
  workingOrder: null,
  viewingOrder: null,
  workingIndex: 0,
  viewingIndex: 0,
  list: [demoOrder, demoOrderCancel, demoOrderPending, demoOrderDone]
};

function reduceStepData(reducer, state: tOrderState, action): tOrderState {
  const newState = { ...state };
  const newStep: ?tStep = workingStep(newState);
  if (newStep) {
    newStep.data = reducer(newStep?.data, action);
  }
  return newState;
}

function setOrderStatus(state: tOrderState, status, selector: (tOrderState)=>?tOrder): tOrderState {
  const newState = { ...state };
  const newOrder = selector(newState);
  if (newOrder) {
    newOrder.status = status;
  }
  return newState;
}

function setStepStatus(state: tOrderState, action): tOrderState {
  const newState = {
    ...state
  };
  const newStep = workingStep(newState);
  if(newStep){
    newStep.status = action.status;
  }
  return newState;
}

function limitIndex(order: ?tOrder, index: tOrderStepIdx): tOrderStepIdx {
  if (index < 0) {
    return 0;
  }
  if (index >= orderLength(order)) {
    return orderLength(order) - 1;
  }
  return index;
}


const orderReducer: { [key: string]: (tOrderState, { type: string, [key: any]: any })=>tOrderState } = {
  [ORDER.VIEW]: (state, action) => {
    const firstIndex = 0;
    const newState = {
      ...state,
      viewingOrder: action?.order || null,
      viewingIndex: firstIndex,
      startTime: new Date()
    };
    const triggered = setOrderStatus(newState, ORDER_STATUS.WIP, workingOrder);
    return {
      ...triggered
    };
  },
  [ORDER.WORK_ON]: (state, action) => {
    const firstIndex = state.workingIndex || 0;
    const newState = {
      ...state,
      workingOrder: action.order || null,
      viewingIndex: firstIndex,
      workingIndex: firstIndex,
      startTime: new Date()
    };
    const triggered = setOrderStatus(newState, ORDER_STATUS.WIP, workingOrder);
    return {
      ...triggered
    };
  },
  [ORDER.FINISH]: (state) => {
    const newState = setOrderStatus(state, ORDER_STATUS.DONE, workingOrder);
    return {
      ...newState
    };
  },
  [ORDER.CANCEL]: (state) => {
    const newState = setOrderStatus(state, ORDER_STATUS.CANCEL, viewingOrder);
    return {
      ...newState,
      workingOrder: null,
      workingIndex:0
    };
  },
  [ORDER.PENDING]: (state) => {
    const newState = setOrderStatus(state, ORDER_STATUS.PENDING, viewingOrder);
    return {
      ...newState,
      workingOrder: null,
      workingIndex:0
    };
  },
  [ORDER.STEP.NEXT]: (state) => {
    const newIndex = limitIndex(viewingOrder(state), viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex,
      workingOrder: null
    };
  },

  [ORDER.STEP.PREVIOUS]: (state) => {
    const newIndex = limitIndex(viewingOrder(state), viewingIndex(state) - 1);
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
  [ORDER.STEP.DO_NEXT]: (state) => {
    const newIndex = workingIndex(state) + 1;
    return {
      ...state,
      workingIndex: newIndex,
      viewingIndex:
        workingStep(state) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: (state) => {
    const newIndex = limitIndex(workingOrder(state), workingIndex(state) - 1);
    const revokedState = {
      ...state,
      workingIndex: newIndex,
      viewingIndex:
        workingStep(state) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
    return setStepStatus(revokedState, STEP_STATUS.READY);
  },
  [ORDER.STEP.DATA]: (state, action) => {
    const { reducer } = action;
    return reduceStepData(reducer, state, action);
  },
  [ORDER.STEP.START_TIME]: (state, action) => {
    const { idx, startTime } = action;
    const newState = { ...state };
    const newStep = getStep(workingOrder(newState), idx);
    if (newStep) {
      newStep.startTime = newStep.startTime || startTime;
    }
    return newState;
  },
  [ORDER.STEP.END_TIME]: (state, action) => {
    const { idx, endTime } = action;
    const newState = { ...state };
    const newStep = getStep(workingOrder(newState), idx);
    if (newStep) {
      newStep.endTime = endTime;
    }
    return newState;
  }
};

export default genReducers(orderReducer, initState);
