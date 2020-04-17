// @flow

import { ORDER } from './constants';
import { genReducers } from '../util';

import type { tAction, tReducer } from '../typeDef';
import type { tOrderActionTypes, tOrderState, tOrderStepIdx } from './interface/typeDef';
import { orderLength, viewingIndex, viewingOrder, workingIndex, workingOrder } from './selector';
import type { IOrder } from './interface/IOrder';

const initState: tOrderState = {
  workingOrder: null,
  viewingOrder: null,
  viewingIndex: 0,
  list: [],
  blockReasons: [{
    lossType: 'availability',
    name: 'Equipment Failure'
  }]
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

const orderReducer: {
  [key: tOrderActionTypes]: tReducer<tOrderState,
    // eslint-disable-next-line flowtype/no-weak-types
    tAction<tOrderActionTypes, any>>
} = {
  [ORDER.NEW_LIST]: (state, { list }: { list: Array<IOrder> }) => ({
    ...state,
    list: list.sort((a, b) => {
      const prev = a.datePlannedStart || 0;
      const latter = b.datePlannedStart || 0;
      return prev - latter;
    })
  }),
  [ORDER.UPDATE_STATE]: state => ({
    ...state
  }),
  [ORDER.VIEW]: (state, { order }: { order: IOrder }) => ({
    ...state,
    viewingOrder: order || null,
    viewingIndex: 0 // first index
  }),
  [ORDER.WORK_ON]: (state, { order, config }: { order: IOrder, config: Object }) => {
    const { step } = config || {};
    let vIdx = workingIndex(order) >= order.steps.length ? 0 : workingIndex(order);
    if (step) {
      const stepIndex = order.steps.findIndex(s => s.code === step.code);
      if (stepIndex >= 0) {
        vIdx = stepIndex;
      }
    }
    vIdx = limitIndex(viewingOrder(state), vIdx);
    return {
      ...state,
      workingOrder: order || null,
      viewingOrder: order || null,
      viewingIndex: vIdx
    };
  },
  [ORDER.DID_FINISH]: state => ({
    ...state,
    workingOrder: null
  }),
  [ORDER.STEP.VIEW_NEXT]: state => ({
    ...state,
    viewingIndex: limitIndex(viewingOrder(state), viewingIndex(state) + 1)
  }),
  [ORDER.STEP.VIEW_PREVIOUS]: state => ({
    ...state,
    viewingIndex: limitIndex(viewingOrder(state), viewingIndex(state) - 1)
  }),
  [ORDER.STEP.JUMP_TO]: (state, { stepIdx }) => ({
    ...state,
    viewingIndex: stepIdx
  }),
  [ORDER.STEP.FINISH](state, { step }) {
    // const wOrder: ?IOrder = workingOrder(state);
    // const newIndex = workingIndex(wOrder) + 1;
    // const vIndex = limitIndex(viewingOrder(state),
    //   step === viewingStep(state) ? newIndex : viewingIndex(state)
    // );
    return {
      ...state
      // viewingIndex: vIndex
    };
  },
  [ORDER.STEP.DO_PREVIOUS](state) {
    // const wOrder = workingOrder(state);
    // const newIndex = limitIndex(wOrder, workingIndex(wOrder) - 1);
    // const vIndex =
    //   workingStep(wOrder) === viewingStep(state)
    //     ? newIndex
    //     : viewingIndex(state);
    return {
      ...state
      // viewingIndex: vIndex
    };
  },
  [ORDER.CLEAR_DATA](state, { order }) {
    const vOrder = viewingOrder(state);
    const wOrder = workingOrder(state);
    if (wOrder !== order) {
      order.clearData();
    }
    if (vOrder === order) {
      return {
        ...state,
        viewingIndex: 0,
        viewingOrder: null
      };
    }
    return state;
  },
  [ORDER.SET_BLOCK_REASON_LIST](state, { blockReasons }) {
    return {
      ...state,
      blockReasons
    };
  }
};

export default genReducers<tOrderState, tOrderActionTypes>(
  orderReducer,
  initState
);
