// @flow

import { ORDER } from './action';
import { genReducers } from '../util';
import { ORDER_STATUS } from './model';
import type { tOrder, tOrderState, tOrderStatus, tOrderStepIdx, tStep } from './model';
import STEP_STATUS from '../step/model';
import {
  demoOrder,
  demoOrderLong,
  demoOrderCancel,
  demoOrderDone,
  demoOrderPending
} from './demoData';
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
  list: [
    demoOrder,
    demoOrderLong,
    demoOrderCancel,
    demoOrderPending,
    demoOrderDone
  ]
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

function clearWorkingOrderIfMatch(state, order) {
  const wOrder = workingOrder(state);
  if (wOrder === order) {
    return {
      ...state,
      workingOrder: null
    };
  }
  return state;
}

function orderStatus(state: tOrderState, order: tOrder, status: tOrderStatus) {
  const newOrder = order;
  if (newOrder) {
    newOrder.status = status;
  }
  return {
    ...state
  };
}

function stepTime(step: ?tStep, isStart: boolean) {
  const newStep = step;
  if (!(newStep && newStep.times)) {
    return;
  }
  const isStarted = newStep.times.length % 2 === 1;
  if ((isStart && !isStarted) || (!isStart && isStarted)) {
    newStep.times.push(new Date());
  }
}


const orderReducer: {
  // eslint-disable-next-line flowtype/no-weak-types
  [key: string]: (tOrderState, { type: string, [key: any]: any }) => tOrderState
} = {
  [ORDER.LIST.SUCCESS]: (state, { list }) => {
    let newList = state.list.filter(
      (o) => !!list.find(newO => o.id === newO.id)
    );
    newList.forEach((o) => {
      Object.assign(o, list.find(newO => o.id === newO.id));
    });
    newList = newList.concat(list.filter(newO =>
      !newList.find(o => o.id === newO.id)
    ));
    return {
      ...state,
      list: newList
    };
  },
  [ORDER.DETAIL.SUCCESS]: (state, { order }) => {
    const newList = [...state.list];
    const newOrder = newList.find(o => o.id === order.id);
    if (newOrder) {
      Object.assign(newOrder, order);
    }
    return {
      ...state,
      list: newList
    };
  },
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
  [ORDER.WORK_ON]: (state, {order}) => {
    const wOrder = workingOrder(state);
    if (wOrder) {
      return state;
    }
    const startIndex = workingIndex(order);
    const newState={
      ...state,
      workingOrder: order || null,
      viewingIndex: startIndex
    };
    return orderStatus(newState,order,ORDER_STATUS.WIP);
  },
  [ORDER.FINISH]: (state, { order }) => clearWorkingOrderIfMatch(
    orderStatus(state, order, ORDER_STATUS.DONE),
    order
  ),
  [ORDER.CANCEL]: (state, { order }) => {
    stepTime(workingStep(order), false);
    return clearWorkingOrderIfMatch(
      orderStatus(state, order, ORDER_STATUS.CANCEL),
      order
    );
  },
  [ORDER.PENDING]: (state, {order}) => {
    stepTime(workingStep(order), false);
    return clearWorkingOrderIfMatch(
      orderStatus(state, order, ORDER_STATUS.PENDING),
      order
    );
  },
  [ORDER.STEP.VIEW_NEXT]: state => {
    const newIndex = limitIndex(viewingOrder(state), viewingIndex(state) + 1);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },
  [ORDER.STEP.VIEW_PREVIOUS]: state => {
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
  [ORDER.STEP.DO_NEXT]: state => {
    const wOrder: ?tOrder = workingOrder(state);
    const newIndex = workingIndex(wOrder) + 1;
    const vIndex =
      workingStep(wOrder) === viewingStep(state)
        ? newIndex
        : viewingIndex(state);
    if (wOrder) {
      wOrder.workingIndex = newIndex;
    }
    return {
      ...state,
      viewingIndex: vIndex
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: state => {
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
