// @flow
import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';
import { ORDER_WS_TYPES } from '../modules/order/constants';

export function orderListApi() {
  try {
    return rushSendApi(ORDER_WS_TYPES.LIST);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderListApi'
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
      at: 'orderDetailApi'
    });
  }
}

// 更新工单状态
export function orderUpdateApi(id: number, orderStatus: string): ?Promise<any> {
  try {
    return rushSendApi(ORDER_WS_TYPES.UPDATE, {
      id,
      status: orderStatus
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderUpdateApi'
    });
  }
}

// 更新工步状态
export function orderStepUpdateApi(id: number, status: string): any {
  try {
    return rushSendApi(ORDER_WS_TYPES.STEP_UPDATE, {
      id,
      status
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderStepUpdateApi'
    });
  }
}

// TODO
export function orderReportStartApi() {
  try {
    return rushSendApi(ORDER_WS_TYPES.REPORT_START, {});
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportStartApi'
    });
  }
}

export function orderReportFinishApi() {
  try {
    return rushSendApi(ORDER_WS_TYPES.REPORT_FINISH, {});
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportFinishApi'
    });
  }
}

export function stepDataApi(data){
  try {
    return rushSendApi(ORDER_WS_TYPES.STEP_DATA, {data});
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportFinishApi'
    });
  }
}
