import { ORDER } from './action';
import { genReducers } from '../indexReducer';
import { ORDER_STEP_STATUS } from './model';
import { demoOrder, demoOrder2 } from './demoData';
import { viewingStep, processingStep, processingIndex, currentOrderLength, viewingIndex, currentOrder } from './selector';

const initState = {
  currentOrder: null,
  processingIndex: 0,
  viewingIndex: 0,
  status: '',
  list: [demoOrder, demoOrder2]
};

function setStepStatus(state, status) {
  const newState = {
    ...state
  };
  const newStep = processingStep(newState);
  newStep.status = status;
  return newState;
}

function limitIndex(index, state) {
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
      processingIndex: firstIndex
    };
  },

  [ORDER.STEP.NEXT]: (state) => {
    const newIndex = limitIndex(viewingIndex(state) + 1, state);
    return {
      ...state,
      viewingIndex: newIndex
    };
  },

  [ORDER.STEP.PREVIOUS]: (state) => {
    const newIndex = limitIndex(viewingIndex(state) - 1, state);
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
  [ORDER.STEP.PUSH]: (state) => {
    const newIndex = processingIndex(state) + 1;
    return {
      ...state,
      processingIndex: newIndex,
      viewingIndex: processingStep(state) === viewingStep(state) ? newIndex : viewingIndex(state)
    };
  },
  [ORDER.STEP.REVOKE]:(state)=>{
    const newIndex = processingIndex(state) - 1;
    return {
      ...state,
      processingIndex: newIndex,
      viewingIndex: processingStep(state) === viewingStep(state) ? newIndex : viewingIndex(state)
    };
  },
  [ORDER.STEP.UPDATE]: (state, action) => {
    const { newStep } = action;
    const newOrder = currentOrder(state).slice();
    newOrder.splice(processingIndex(state), 1, newStep);
    return {
      ...state,
      currentOrder: newOrder
    };
  }
};

export default genReducers(orderReducer, initState);
