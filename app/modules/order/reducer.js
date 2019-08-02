// @flow
import { ORDER } from './action';
import { genReducers } from '../util';
import type { tOrder, tOrderState, tOrderStepIdx, tStep } from './model';
import { ORDER_STATUS } from './model';
import STEP_STATUS from '../step/model';
import { demoOrder, demoOrderCancel, demoOrderDone, demoOrderPending } from './demoData';
import {
  getStep,
  orderLength,
  viewingIndex,
  viewingOrder,
  viewingStep,
  workingIndex,
  workingOrder,
  workingStep
} from './selector';

const initState = {
  workingOrder: null,
  viewingOrder: null,
  viewingIndex: 0,
  list: [demoOrder, demoOrderCancel, demoOrderPending, demoOrderDone]
};

function reduceStepData(reducer, state: tOrderState, action): tOrderState {
  const newState = { ...state };
  const newStep: ?tStep = workingStep(workingOrder(newState));
  if (newStep) {
    newStep.data = reducer(newStep?.data, action);
  }
  return newState;
}

function setStepStatus(state: tOrderState, action): tOrderState {
  const newState = {
    ...state
  };
  const newStep = workingStep(workingOrder(newState));
  if (newStep) {
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
    const { order } = action;
    const firstIndex = 0;
    return {
      ...state,
      viewingOrder: order || null,
      viewingIndex: firstIndex,
      startTime: new Date()
    };
  },
  [ORDER.WORK_ON]: (state, action) => {
    const wOrder = workingOrder(state);
    if (wOrder) {
      return state;
    }
    const { order } = action;
    const startIndex = workingIndex(order);
    order.status = ORDER_STATUS.WIP;
    return {
      ...state,
      workingOrder: order || null,
      viewingIndex: startIndex
    };
  },
  [ORDER.FINISH]: (state) => {
    const wOrder = workingOrder(state);
    if (wOrder) {
      wOrder.status = ORDER_STATUS.DONE;
    }
    return {
      ...state,
      workingOrder: null
    };
  },
  [ORDER.CANCEL]: (state) => {
    const vOrder = viewingOrder(state);
    if (vOrder) {
      vOrder.status = ORDER_STATUS.CANCEL;
    }
    const wStep = workingStep(viewingOrder(state));
    if (wStep && wStep.times && wStep.times.length % 2 === 1) {
      wStep.times.push(new Date());
    }
    const wOrder = workingOrder(state);
    const newWOrder = wOrder === vOrder ? null : wOrder;
    return {
      ...state,
      workingOrder: newWOrder
    };
  },
  [ORDER.PENDING]: (state, action) => {
    const { order } = action;
    order.status = ORDER_STATUS.PENDING;
    const wStep = workingStep(order);
    if (wStep && wStep.times && wStep.times.length % 2 === 1) {
      wStep.times.push(new Date());
    }
    const vOrder = viewingOrder(state);
    const wOrder = workingOrder(state);
    const newWOrder = wOrder === vOrder ? null : wOrder;
    return {
      ...state,
      workingOrder: newWOrder
    };
  },
  [ORDER.STEP.NEXT]: (state) => {
    const newIndex = limitIndex(viewingOrder(state), viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex
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
    const wOrder = workingOrder(state);
    const newIndex = workingIndex(wOrder) + 1;
    const vIndex = workingStep(wOrder) === viewingStep(state) ? newIndex
      : viewingIndex(state);
    if (wOrder) {
      wOrder.workingIndex = newIndex;
    }
    return {
      ...state,
      viewingIndex: vIndex
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: (state) => {
    const wOrder = workingOrder(state);
    const newIndex = limitIndex(wOrder, workingIndex(wOrder) - 1);
    if (wOrder) {
      wOrder.workingIndex = newIndex;
    }
    const revokedState = {
      ...state,
      viewingIndex:
        workingStep(wOrder) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
    return setStepStatus(revokedState, STEP_STATUS.READY);
  },
  [ORDER.STEP.DATA]: (state, action) => {
    const { reducer } = action;
    return reduceStepData(reducer, state, action);
  },
  [ORDER.STEP.TIME]: (state, action) => {
    const { idx, time } = action;
    const newState = { ...state };
    const newStep = getStep(workingOrder(newState), idx);
    if (newStep) {
      if (newStep.times) {
        newStep.times.push(time);
      } else {
        newStep.times = [time];
      }
    }
    return newState;
  }
};

export default genReducers(orderReducer, initState);
