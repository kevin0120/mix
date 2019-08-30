// @flow

import { ORDER } from './action';
import { genReducers } from '../util';
import type { tOrder, tOrderState, tOrderStepIdx } from './model';
import { demoOrder, demoOrderCancel, demoOrderDone, demoOrderLong, demoOrderPending } from './demoData';
import {
  orderLength,
  viewingIndex,
  viewingOrder,
  viewingStep,
  workingIndex,
  workingOrder,
  workingStep
} from './selector';
import stepTypes from '../step/stepTypes';
import Order from './Order';

const initState = {
  workingOrder: null,
  viewingOrder: null,
  viewingIndex: 0,
  list: [
    demoOrder,
    // demoOrderLong,
    // demoOrderCancel,
    // demoOrderPending,
    // demoOrderDone
  ].map(o => new Order(o, stepTypes))
};


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

const orderReducer: {
  // eslint-disable-next-line flowtype/no-weak-types
  [key: string]: (tOrderState, { type: string, [key: any]: any }) => tOrderState
} = {
  [ORDER.LIST.SUCCESS]: (state, { list }) => {
    // get exist orders, orders not in the new list will be removed!!
    let newList = state.list.filter(
      (o) => !!list.find(newO => o.id === newO.id)
    );

    // update order data
    newList.forEach((o) => {
      const orderData = list.find(newO => o.id === newO.id);
      o.update(orderData, stepTypes);
    });

    // make new orders
    newList = newList.concat(
      list.filter(newO =>
        !newList.find(o => o.id === newO.id)
      ).map(
        oD => new Order(oD, stepTypes)
      )
    );
    console.log(newList);

    // const newList = list.map(oD => new Order(oD, stepTypes));

    return {
      ...state,
      list: newList
    };
  },
  [ORDER.NEW]: (state, { list }) => {
    // get exist orders
    let newList = state.list;

    // update order data
    newList.forEach((o) => {
      const orderData = list.find(newO => o.id === newO.id);
      if(orderData){
        o.update(orderData, stepTypes);
      }
    });

    // make new orders
    newList = newList.concat(
      list.filter(newO =>
        !newList.find(o => o.id === newO.id)
      ).map(
        oD => new Order(oD, stepTypes)
      )
    );

    return {
      ...state,
      list: newList
    };
  },
  [ORDER.DETAIL.SUCCESS]: (state, { order }) => {
    const newList = [...state.list];
    const newOrder = newList.find(o => o.id === order.id);
    if (newOrder) {
      newOrder.update(order, stepTypes);
    }
    return {
      ...state,
      list: newList
    };
  },
  [ORDER.UPDATE_STATE]: (state, action) => ({
    ...state
  }),
  [ORDER.VIEW]: (state, action) => {
    const { order } = action;
    const firstIndex = 0;
    return {
      ...state,
      viewingOrder: order || null,
      viewingIndex: firstIndex
    };
  },
  [ORDER.WORK_ON]: (state, { order }) => {
    const wOrder = workingOrder(state);
    if (wOrder) {
      return state;
    }
    const startIndex = workingIndex(order) >= order.steps.length ? 0 : workingIndex(order);
    return {
      ...state,
      workingOrder: order || null,
      viewingIndex: startIndex
    };
  },
  [ORDER.FINISH]: (state, { order }) => clearWorkingOrderIfMatch(state, order),
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
  [ORDER.STEP.DO_NEXT]: state => {
    const wOrder: ?tOrder = workingOrder(state);
    const newIndex = workingIndex(wOrder) + 1;
    const vIndex =
      workingStep(wOrder) === viewingStep(state)
        ? newIndex
        : viewingIndex(state);
    return {
      ...state,
      viewingIndex: vIndex
    };
  },
  [ORDER.STEP.FINISH]: state => {
    const wOrder: ?tOrder = workingOrder(state);
    const newIndex = workingIndex(wOrder) + 1;
    const vIndex =
      workingStep(wOrder) === viewingStep(state)
        ? newIndex
        : viewingIndex(state);
    return {
      ...state,
      viewingIndex: vIndex
    };
  },
  [ORDER.STEP.DO_PREVIOUS]: state => {
    const wOrder = workingOrder(state);
    const newIndex = limitIndex(wOrder, workingIndex(wOrder) - 1);
    return {
      ...state,
      viewingIndex:
        workingStep(wOrder) === viewingStep(state)
          ? newIndex
          : viewingIndex(state)
    };
  }
};

export default genReducers(orderReducer, initState);
