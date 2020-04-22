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

// 开工
export function orderReportStartApi(
  code: 工单号,
  workCenterCode: 工位号,
  startTime: Date
): anyPromise {
  const dateStartString = Moment(startTime).format();
  return rushSendApi(ORDER_WS_TYPES.START_REQUEST, {
    code,
    start_time: dateStartString,
    workcenter_code: workCenterCode
  });
}

// 完工
export function orderReportFinishApi(
  code: 工单号,
  workCenterCode: 工位号,
  endTime: Date
): anyPromise {
  const dateCompleteString = Moment(endTime).format();
  return rushSendApi(ORDER_WS_TYPES.FINISH_REQUEST, {
    except_type: 'finish',
    except_code: 'finish',
    end_time: dateCompleteString,
    order_name: code,
    workcenter_code: workCenterCode
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

export function orderPendingApi(
  exceptType,
  exceptCode,
  endTime: Date,
  orderCode: 工单号,
  workCenterCode: 工位号
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
