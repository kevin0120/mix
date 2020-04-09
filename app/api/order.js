// @flow
import Moment from 'moment';
import { rushSendApi } from './rush';
import { ORDER_WS_TYPES } from '../modules/order/constants';
import type { tUuid } from '../modules/user/interface/typeDef';
import { CommonLog, defaultClient } from '../common/utils';

export function orderListApi({
  timeFrom,
  timeTo,
  status,
  pageSize,
  pageNo
}: {
  timeFrom?: string,
  timeTo?: string,
  status?: string,
  pageSize?: number,
  pageNo?: number
}): Promise<{}> {
  return rushSendApi(ORDER_WS_TYPES.LIST, {
    time_from: timeFrom, // UTC
    time_to: timeTo,
    status,
    page_size: pageSize,
    page_no: pageNo
  });
}

export function orderDetailApi(id: number): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.DETAIL, {
    id
  });
}

// 更新工单状态
export function orderUpdateApi(code: string, orderStatus: string): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.UPDATE, {
    workorder_code: code,
    status: orderStatus
  });
}

// 更新工步状态
export function orderStepUpdateApi(code: string, status: string): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.STEP_UPDATE, {
    workstep_code: code,
    status
  });
}

// 开工
export function orderReportStartApi(
  code: string,
  trackCode: string,
  workCenter: string,
  productCode: string,
  dateStart: Date,
  resources: {
    user: Array<{
      name: string,
      code: string
    }>,
    equipments: Array<{
      name: string,
      code: string
    }>
  }
): Promise<any> {
  const dateStartString = Moment(dateStart).format();
  return rushSendApi(ORDER_WS_TYPES.START_REQUEST, {
    code,
    track_code: trackCode,
    workcenter: workCenter,
    product_code: productCode,
    date_start: dateStartString,
    resources
  });
}

// 完工
export function orderReportFinishApi(
  code: string,
  trackCode: string,
  productCode: string,
  workCenter: string,
  dateComplete: Date,
  operation: {}
) {
  const dateCompleteString = Moment(dateComplete).format();
  return rushSendApi(ORDER_WS_TYPES.FINISH_REQUEST, {
    code,
    track_code: trackCode,
    workcenter: workCenter,
    product_code: productCode,
    date_complete: dateCompleteString,
    operation
  });
}

export function stepDataApi(code: string, data: Object): Promise<any> {
  // let idNum = code;
  // if (typeof id === 'string') {
  //   idNum = parseInt(code, 10);
  // }
  const json = JSON.stringify(data);
  return rushSendApi(ORDER_WS_TYPES.STEP_DATA, { workstep_code: code, data: json });
}

export function orderDataApi(id: code, data: Object): Promise<any> {
  // let idNum = id;
  // if (typeof id === 'string') {
  //   idNum = parseInt(id, 10);
  // }
  const json = JSON.stringify(data);
  return rushSendApi(ORDER_WS_TYPES.ORDER_DATA, { workorder_code: code, data: json });
}

export function orderDetailByCodeApi(
  code: string,
  workcenter?: string
): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.ORDER_DETAIL_BY_CODE, { code, workcenter });
}

export function orderPendingApi(
  exceptType,
  exceptCode,
  endTime,
  orderCode,
  workCenterCode
) {
  return rushSendApi(ORDER_WS_TYPES.ORDER_PENDING, {
    except_type: exceptType,
    except_code: exceptCode,
    end_time: endTime,
    order_name: orderCode,
    workcenter_code: workCenterCode
  });
}

export function orderResumeApi(
  startTime,
  orderCode,
  workCenterCode
) {
  return rushSendApi(ORDER_WS_TYPES.ORDER_RESUME, {
    start_time: startTime,
    order_name: orderCode,
    workcenter_code: workCenterCode
  });
}

export function getBlockReasonsApi(url: string) {
  const fullUrl = `${url}/mrp.workcenter.productivity.loss`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      CommonLog.lError(e, {
        at: 'getBlockReasonsApi',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
