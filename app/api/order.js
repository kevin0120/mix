// @flow
import Moment from 'moment';
import { rushSendApi } from './rush';
import { ORDER_WS_TYPES } from '../modules/order/constants';
import type {
  工位号,
  人员列表,
  工单号
} from '../modules/order/interface/typeDef';

type anyPromise = Promise<any>;

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

export function orderDetailApi(id: number): anyPromise {
  return rushSendApi(ORDER_WS_TYPES.DETAIL, {
    id
  });
}

// 更新工单状态
export function orderUpdateApi(code: 工单号, orderStatus: string): anyPromise {
  return rushSendApi(ORDER_WS_TYPES.UPDATE, {
    workorder_code: code,
    status: orderStatus
  });
}

// 更新工步状态
export function orderStepUpdateApi(code: string, status: string): anyPromise {
  return rushSendApi(ORDER_WS_TYPES.STEP_UPDATE, {
    workstep_code: code,
    status
  });
}

type 开工人员 = Array<{ name: string, code: string }>;
type 开工设备 = Array<{ name: string, code: string }>;

// 开工
export function orderReportStartApi(
  code: 工单号,
  trackCode: string,
  workCenter: string,
  productCode: string,
  dateStart: Date,
  resources: { user: 开工人员, equipments: 开工设备 }
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
  code: 工单号,
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

// eslint-disable-next-line flowtype/no-weak-types
export function stepDataApi(code: string, data: Object): anyPromise {
  // let idNum = code;
  // if (typeof id === 'string') {
  //   idNum = parseInt(code, 10);
  // }
  const json = JSON.stringify(data);
  return rushSendApi(ORDER_WS_TYPES.STEP_DATA, {
    workstep_code: code,
    data: json
  });
}

// eslint-disable-next-line flowtype/no-weak-types
export function orderDataApi(code: 工单号, data: Object): anyPromise {
  const json = JSON.stringify(data);
  return rushSendApi(ORDER_WS_TYPES.ORDER_DATA, {
    workorder_code: code,
    data: json
  });
}

export function orderDetailByCodeApi(
  code: 工单号,
  workcenter?: 工位号
): anyPromise {
  return rushSendApi(ORDER_WS_TYPES.ORDER_DETAIL_BY_CODE, { code, workcenter });
}

// 产前模拟
export function apiOrderStartSimulate(
  code: 工单号,
  employee: 人员列表,
  workCenter: 工位号
): anyPromise {
  return rushSendApi(ORDER_WS_TYPES.WS_ORDER_SIMULATE, {
    employee,
    order_name: code,
    workcenter_code: workCenter
  })
    .then(resp => {
      const { Message, Result, Status } = resp;
      if (Message) {
        throw new Error(Message);
      }
      if (Status === 'E') {
        throw new Error('产前模拟失败');
      }
      // eslint-disable-next-line promise/always-return
      if (Result === 'N') {
        throw new Error('产前模拟失败');
      }
      return { success: '产前模拟成功' };
    })
    .catch(err => ({ errorMessage: err.message }));
}
