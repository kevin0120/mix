// @flow

// 工单状态
export const ORDER_STATUS = {
  TODO: 'todo',
  WIP: 'wip',
  DONE: 'done',
  CANCEL: 'cancel',
  PENDING: 'pending'
};
// rush type
export const ORDER_WS_TYPES = {
  LIST: 'WS_ORDER_LIST',
  DETAIL: 'WS_ORDER_DETAIL',
  UPDATE: 'WS_ORDER_UPDATE',
  STEP_UPDATE: 'WS_ORDER_STEP_UPDATE',
  NEW: 'WS_NEW_ORDER'
};

// order action types
export const ORDER = {
  WORK_ON: 'ORDER_WORK_ON',
  VIEW: 'ORDER_VIEW',
  FINISH: 'ORDER_FINISH',
  // update the store
  UPDATE_STATE: 'ORDER_UPDATE_STATE',
  NEW: 'ORDER_NEW',
  LIST: {
    GET: 'ORDER_LIST_GET',
    SUCCESS: 'ORDER_LIST_SUCCESS',
    FAIL: 'ORDER_LIST_FAIL'
  },
  DETAIL: {
    GET: 'ORDER_DETAIL_GET',
    SUCCESS: 'ORDER_DETAIL_SUCCESS',
    FAIL: 'ORDER_DETAIL_FAIL'
  },
  STEP: {
    // 仅移动指针，不修改step状态
    NEXT: 'ORDER_STEP_NEXT',
    PREVIOUS: 'ORDER_STEP_PREVIOUS',
    VIEW_PREVIOUS: 'ORDER_STEP_VIEW_PREVIOUS', // 防抖后事件
    VIEW_NEXT: 'ORDER_STEP_VIEW_NEXT', // 防抖后事件
    JUMP_TO: 'ORDER_STEP_JUMP_TO',
    // 修改step的状态
    STATUS: 'ORDER_STEP_STATUS',
    // 步进、步退
    DO_NEXT: 'ORDER_STEP_DO_NEXT',
    DO_PREVIOUS: 'ORDER_STEP_DO_PREVIOUS',
    FINISH: 'ORDER_STEP_FINISH'
  }
};


