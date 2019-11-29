// @flow
import Moment from 'moment';
import { rushSendApi } from './rush';
import { ORDER_WS_TYPES } from '../modules/order/constants';

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
export function orderUpdateApi(id: number, orderStatus: string): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.UPDATE, {
    id,
    status: orderStatus
  });
}

// 更新工步状态
export function orderStepUpdateApi(id: number, status: string): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.STEP_UPDATE, {
    id,
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

export function stepDataApi(id: number, data: Object): Promise<any> {
  let idNum = id;
  if (typeof id === 'string') {
    idNum = parseInt(id, 10);
  }
  const json = JSON.stringify(data);
  console.warn(json);
  return rushSendApi(ORDER_WS_TYPES.STEP_DATA, { id: idNum, data: json });
}

export function orderDataApi(id: number, data: Object): Promise<any> {
  let idNum = id;
  if (typeof id === 'string') {
    idNum = parseInt(id, 10);
  }
  const json = JSON.stringify(data);
  return rushSendApi(ORDER_WS_TYPES.ORDER_DATA, { id: idNum, data: json });
}

export function orderDetailByCodeApi(
  code: string,
  workcenter?: string
): Promise<any> {
  return rushSendApi(ORDER_WS_TYPES.ORDER_DETAIL_BY_CODE, { code, workcenter });
}
