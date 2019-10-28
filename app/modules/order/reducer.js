// @flow

import { ORDER } from './constants';
import { genReducers } from '../util';

import type { tAction, tReducer } from '../typeDef';
import type {
  tOrderState,
  tOrderStepIdx,
  tOrderActionTypes
} from './interface/typeDef';
import {
  demoOrder,
  demoOrderCancel,
  demoOrderDone,
  demoOrderLong,
  demoOrderPending
} from './demoData';
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

function limitIndex(order: ?IOrder, index: tOrderStepIdx): tOrderStepIdx {
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
  [key: tOrderActionTypes]: tReducer<
    tOrderState,
    tAction<tOrderActionTypes, any>
  >
} = {
  [ORDER.NEW_LIST]: (state, { list }) => ({
    ...state,
    list
  }),
  [ORDER.UPDATE_STATE]: state => ({
    ...state
  }),
  [ORDER.VIEW]: (state, { order }: { order: IOrder }) => ({
    ...state,
    viewingOrder: order || null,
    viewingIndex: 0 // first index
  }),
  [ORDER.WORK_ON]: (state, { order }: { order: IOrder }) => ({
    ...state,
    workingOrder: order || null,
    viewingIndex:
      workingIndex(order) >= order.steps.length ? 0 : workingIndex(order)
  }),
  [ORDER.FINISH]: (state, { order }: { order: IOrder }) =>
    clearWorkingOrderIfMatch(state, order),
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
  [ORDER.STEP.FINISH](state, { step }) {
    const wOrder: ?IOrder = workingOrder(state);
    const newIndex = workingIndex(wOrder) + 1;
    const vIndex = step === viewingStep(state) ? newIndex : viewingIndex(state);
    return {
      ...state,
      viewingIndex: vIndex
    };
  },
  [ORDER.STEP.DO_PREVIOUS](state) {
    const wOrder = workingOrder(state);
    const newIndex = limitIndex(wOrder, workingIndex(wOrder) - 1);
    const vIndex =
      workingStep(wOrder) === viewingStep(state)
        ? newIndex
        : viewingIndex(state);
    return {
      ...state,
      viewingIndex: vIndex
    };
  }
};

export default genReducers<tOrderState, tOrderActionTypes>(
  orderReducer,
  initState
);
