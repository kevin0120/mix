// @flow
import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';
import { ORDER_WS_TYPES } from '../modules/order/constents';

export function orderListApi() {
  try {
    return rushSendApi(ORDER_WS_TYPES.LIST);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi'
    });
  }
}

export function orderDetailApi(id: number) {
  try {
    return rushSendApi(ORDER_WS_TYPES.DETAIL, {
      id
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi'
    });
  }
}


// 更新工单状态
export function orderUpdateApi(id: number, orderStatus: string) {
  try {
    return rushSendApi(ORDER_WS_TYPES.UPDATE, {
      id,
      status: orderStatus
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi'
    });
  }
}

// 更新工步状态
export function orderStepUpdateApi(id: number, status: string) {
  try {
    return rushSendApi(ORDER_WS_TYPES.STEP_UPDATE, {
      id, status
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'jobApi'
    });
  }
}


