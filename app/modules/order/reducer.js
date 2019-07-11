import { ORDER } from './action';
import { genReducers } from '../indexReducer';
import { ORDER_STEP_STATUS, hasSubStep } from './model';
import { demoOrder, demoOrder2 } from './demoData';

const initState = {
  currentOrder: null,
  currentProcessingIndex: 0,
  currentProcessingStep: {},
  currentViewingIndex: 0,
  currentViewingStep: {},
  status: '',
  list: [demoOrder, demoOrder2]
};

function setStepStatus(state, status) {
  const { currentProcessingIndex, currentOrder, currentViewingIndex } = state;
  const newState = {
    ...state,
    currentOrder: {
      ...currentOrder
    }
  };
  const newStep = getStepByIndex(currentProcessingIndex, newState.currentOrder);
  newStep.status = status;
  return {
    ...newState,
    currentProcessingStep: newStep,
    currentViewingStep: getStepByIndex(currentViewingIndex, newState.currentOrder)
  };
}

function getStepByIndex(index, order) {
  return order.steps && order.steps[index];
}

function nextIndex(index) {
  return index + 1;
}

function limitIndex(index, order) {
  if (index < 0) {
    return 0;
  }
  if (index >= order.steps.length) {
    return order.steps.length - 1;
  }
  return index;
}

function prevIndex(index) {
  return index - 1;
}

const orderReducer = {
  [ORDER.TRIGGER]: (state, action) => {
    const firstIndex = 0;
    return {
      ...state,
      currentOrder: action.order || null,
      currentViewingIndex: firstIndex,
      currentViewingStep: getStepByIndex(firstIndex, action.order),
      currentProcessingIndex: firstIndex,
      currentProcessingStep: getStepByIndex(firstIndex, action.order)

    };
  },

  [ORDER.STEP.NEXT]: (state) => {
    const { currentViewingIndex, currentOrder } = state;
    const newIndex = limitIndex(nextIndex(currentViewingIndex), currentOrder);
    return {
      ...state,
      currentViewingIndex: newIndex,
      currentViewingStep: getStepByIndex(newIndex, currentOrder)
    };
  },

  [ORDER.STEP.PREVIOUS]: (state) => {
    const { currentViewingIndex, currentOrder } = state;
    const newIndex = limitIndex(prevIndex(currentViewingIndex), currentOrder);
    return {
      ...state,
      currentViewingIndex: newIndex,
      currentViewingStep: getStepByIndex(newIndex, currentOrder)
    };
  },

  [ORDER.STEP.JUMP_TO]: (state, action) => {
    const { stepId } = action;
    return {
      ...state,
      currentViewingIndex: stepId,
      currentViewingStep: getStepByIndex(stepId, state.currentOrder)
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
    const { currentProcessingIndex, currentOrder, currentViewingStep, currentProcessingStep } = state;
    const newIndex = nextIndex(currentProcessingIndex);
    const newProcessingStep = getStepByIndex(newIndex, currentOrder);
    const newState = {
      ...state,
      currentProcessingIndex: newIndex,
      currentProcessingStep: newProcessingStep
    };
    if (currentProcessingStep === currentViewingStep) {
      newState.currentViewingStep = newProcessingStep;
      newState.currentViewingIndex = newIndex;
    }
    return newState;
  }
  // [ORDER.FINISH]:
};

export default genReducers(orderReducer, initState);
