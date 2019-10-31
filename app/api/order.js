// @flow
import Moment from 'moment';
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
export function orderStepUpdateApi(id: number, status: string): ?Promise<any> {
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
export function orderReportStartApi(
  code: string,
  trackCode: string,
  workCenter: string,
  productCode: string,
  dateStart: Date,
  resources: {
    user: [{
      name: string,
      code: string
    }],
    equipments: [
      {
        name: string,
        code: string
      }
    ]

  }
): ?Promise<{
  error_code: number,
  msg: string,
  extra: string
}> {
  try {
    const dateStartString=Moment(dateStart).format();
    return rushSendApi(ORDER_WS_TYPES.START_REQUEST, {
      code,
      track_code: trackCode,
      workcenter: workCenter,
      product_code: productCode,
      date_start: dateStartString,
      resources
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportStartApi'
    });
  }
}

export function orderReportFinishApi(
  code: string,
  trackCode: string,
  workCenter: string,
  productCode: string,
  operation: {}
) {
  try {
    return rushSendApi(ORDER_WS_TYPES.FINISH_REQUEST, {
      code,
      track_code: trackCode,
      workcenter: workCenter,
      product_code: productCode,
      operation
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportFinishApi'
    });
  }
}

export function stepDataApi(data: Object): ?Promise<any> {
  try {
    return rushSendApi(ORDER_WS_TYPES.STEP_DATA, { data });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'orderReportFinishApi'
    });
  }
}
