// @flow

import type { Action, Reducer } from 'redux';
import { ORDER } from './constants';
import { genReducers } from '../util';

import type { tAction, tReducer } from '../typeDef';
import type { tOrder, tOrderState, tOrderStepIdx, tOrderActionTypes, tOrderListData } from './interface/typeDef';
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
import OrderMixin from './Order';
import Step from '../step/Step';
import type { IOrder } from './interface/IOrder';

const initState: tOrderState = {
  workingOrder: null,
  viewingOrder: null,
  viewingIndex: 0,
  list: [
    demoOrder,
    demoOrderLong,
    demoOrderCancel,
    demoOrderPending,
    demoOrderDone
  ].map(o => new (OrderMixin(Step))(o))
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

function clearWorkingOrderIfMatch(state: tOrderState, order: IOrder) {
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
  [key: tOrderActionTypes]: tReducer<tOrderState, Action<tOrderActionTypes>>
} = {
  [ORDER.LIST.SUCCESS](state, { list }: { list: Array<tOrderListData> }) {
    // get exist orders, orders not in the new list will be removed!!
    let newList = state && state.list.filter(o => !!list.find(newO => o.id === newO.id));

    // update order data
    newList.forEach(o => {
      const orderData = list.find(newO => o.id === newO.id);
      o.update(orderData);
    });

    // make new orders
    newList = newList.concat(
      list
        .filter(newO => !newList.find(o => o.id === newO.id))
        .map(oD => new (OrderMixin(Step))(oD))
    );
    // const newList = list.map(oD => new Order(oD, stepTypes));
    return {
      ...state,
      list: newList
    };
  },
  [ORDER.NEW](state, { list }: { list: Array<tOrderListData> }) {
    // get exist orders
    let newList = state.list;

    // update order data
    newList.forEach(o => {
      const orderData = list.find(newO => o.id === newO.id);
      if (orderData) {
        o.update(orderData);
      }
    });

    // make new orders
    newList = newList.concat(
      list
        .filter(newO => !newList.find(o => o.id === newO.id))
        .map(oD => new (OrderMixin(Step))(oD))
    );

    return {
      ...state,
      list: newList
    };
  },
  [ORDER.DETAIL.SUCCESS](state, { order }: { order: IOrder }) {
    const newList = [...state.list];
    const newOrder = newList.find(o => o.id === order.id);
    if (newOrder) {
      newOrder.update(order);
    }
    return {
      ...state,
      list: newList
    };
  },
  [ORDER.UPDATE_STATE](state) {
    return { ...state };
  },
  [ORDER.VIEW](state, { order }: { order: IOrder }) {
    const firstIndex = 0;
    return {
      ...state,
      viewingOrder: order || null,
      viewingIndex: firstIndex
    };
  },
  [ORDER.WORK_ON](state, { order }: { order: IOrder }) {
    const wOrder = workingOrder(state);
    if (wOrder) {
      return state;
    }
    return {
      ...state,
      workingOrder: order || null,
      viewingIndex: workingIndex(order) >= order.steps.length ? 0 : workingIndex(order)
    };
  },
  [ORDER.FINISH]: (state, { order }: { order: IOrder }) => clearWorkingOrderIfMatch(state, order),
  [ORDER.STEP.VIEW_NEXT]: state => ({
    ...state,
    viewingIndex: limitIndex(viewingOrder(state), viewingIndex(state) + 1)
  }),
  [ORDER.STEP.VIEW_PREVIOUS]: state => ({
    ...state,
    viewingIndex: limitIndex(viewingOrder(state), viewingIndex(state) - 1)
  }),
  [ORDER.STEP.JUMP_TO]: (state, { stepId }) => ({
    ...state,
    viewingIndex: stepId
  }),
  [ORDER.STEP.DO_NEXT](state) {
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
  [ORDER.STEP.FINISH](state) {
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
  [ORDER.STEP.DO_PREVIOUS](state) {
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

export default genReducers<tOrderState, tOrderActionTypes>(orderReducer, initState);
