import { ORDER } from './action';
import { genReducers } from '../indexReducer';
import { ORDER_STEP_STATUS, hasSubStep } from './model';
import { demoOrder, demoOrder2 } from './demoData';

const initState = {
  currentOrder: null,
  currentProcessingIndex: [],
  currentProcessingStep: {},
  currentViewingIndex: [],
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

function stepNext(index, order) {
  if (!(order && order.steps && order.steps.length && order.steps.length > 0)) {
    return [];
  }
  let depthHasNextBrother = -1;
  let newIndex = index.slice();
  const depth = newIndex.length;
  const currentIndex = newIndex.reduce((s, idx, i) => {
    depthHasNextBrother = s.steps && s.steps[idx + 1] ? i : depthHasNextBrother;
    if (!hasSubStep(s)) {
      return s;
    }
    return s.steps[idx];
  }, order);
  if (hasSubStep(currentIndex)) {
    newIndex.push(0);
  } else if (depthHasNextBrother >= 0) {
    const backDepth = depth - depthHasNextBrother;
    newIndex.push(newIndex.splice(-backDepth, backDepth)[0] + 1);
    // newIndex = firstChildIndex(newIndex, order);
  }
  return newIndex;
}

function stepPrevious(index, order) {
  if (!(order && order.steps && order.steps.length && order.steps.length > 0)) {
    return [];
  }
  let depthHasPreviousBrother = -1;
  let newIndex = index.slice();
  const depth = newIndex.length;
  newIndex.reduce((s, idx, i) => {
    depthHasPreviousBrother = s.steps && s.steps[idx - 1] ? i : depthHasPreviousBrother;
    if (!hasSubStep(s)) {
      return s;
    }
    return s.steps[idx];
  }, order);
  if (depthHasPreviousBrother >= 0) {
    const backDepth = depth - depthHasPreviousBrother;
    newIndex.push(newIndex.splice(-backDepth, backDepth)[0] - 1);
    // newIndex = lastChildIndex(newIndex, order);
  }
  return newIndex;
}

function getStepByIndex(index, order) {
  if (!(order && order.steps && order.steps.length && order.steps.length > 0)) {
    return {};
  }
  return index.reduce((s, idx) => s && s.steps && s.steps[idx], order);
}


function lastChildIndex(index, order) {
  const newIndex = index.slice();
  // eslint-disable-next-line no-constant-condition
  while (true) {
    const currStep = getStepByIndex(newIndex, order);
    if (!(currStep && currStep.steps && currStep.steps.length >= 0)) {
      break;
    }
    newIndex.push(currStep.steps.length - 1);
  }
  return newIndex;
}

function firstChildIndex(index, order) {
  const newIndex = index.slice();
  // eslint-disable-next-line no-constant-condition
  while (true) {
    const currStep = getStepByIndex(newIndex, order);
    if (!(currStep && currStep.steps && currStep.steps.length >= 0)) {
      break;
    }
    newIndex.push(0);
  }
  return newIndex;
}

const orderReducer = {
  [ORDER.TRIGGER]: (state, action) => {
    const firstIndex=firstChildIndex([],action.order);
    return {
    ...state,
    currentOrder: action.order || null,
    currentViewingIndex: firstIndex,
    currentViewingStep: getStepByIndex(firstIndex, action.order),
    currentProcessingIndex: firstIndex,
    currentProcessingStep: getStepByIndex(firstIndex, action.order)

  }},

  [ORDER.STEP.NEXT]: (state) => {
    const { currentViewingIndex, currentOrder } = state;
    const newIndex = stepNext(currentViewingIndex, currentOrder);
    return {
      ...state,
      currentViewingIndex: newIndex,
      currentViewingStep: getStepByIndex(newIndex, currentOrder)
    };
  },

  [ORDER.STEP.PREVIOUS]: (state) => {
    const { currentViewingIndex, currentOrder } = state;
    const newIndex = stepPrevious(currentViewingIndex, currentOrder);
    return {
      ...state,
      currentViewingIndex: newIndex,
      currentViewingStep: getStepByIndex(newIndex, currentOrder)
    };
  },

  [ORDER.STEP.JUMP_TO]: (state, action) => {
    const firstChild = firstChildIndex(action.stepId, state.currentOrder);
    return {
      ...state,
      currentViewingIndex: firstChild,
      currentViewingStep: getStepByIndex(firstChild, state.currentOrder)
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
  [ORDER.COLLECT]: (state, action) => {
    const { payload } = action;
    const { currentProcessingIndex, currentOrder } = state;
    const newOrder = currentOrder.slice();
    newOrder.splice(currentProcessingIndex, 1, {
      ...newOrder[currentProcessingIndex],
      stash: {
        ...newOrder[currentProcessingIndex],
        ...payload
      }
    });
    return {
      ...state,
      newOrder
    };
  },
  //
  [ORDER.STEP.PUSH]: (state) => {
    const { currentProcessingIndex, currentOrder, currentViewingStep, currentProcessingStep } = state;
    const newIndex = stepNext(currentProcessingIndex, currentOrder);
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
};

export default genReducers(orderReducer, initState);
