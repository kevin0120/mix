import STEP_STATUS from '../../step/constents';
import { ORDER_STATUS } from '../constents';
import type { tClsOrder } from '../Order';

export type tStockMove = {
  lot: string, // 批次号或者序列号
  product: string // 所选的产品名称
};

export type tOrder = {
  id: number,
  canRework: boolean, // 是否能够返工
  incomingProducts: Array<tStockMove>,
  finishedProducts: Array<tStockMove>, //
  steps: tStepArray, // 工步
  status: tOrderStatus, // 工单状态
  plannedDateTime: string, // 计划时间
  name: string, // 工单号
  desc: string, // 工单信息
  workingIndex: ?number // 正在执行的工步索引
};

export type tOrderStepIdx = number;

export type tOrderState = {
  viewingOrder: tClsOrder | null,
  workingOrder: tClsOrder | null,
  viewingIndex: tOrderStepIdx,
  list: Array<tClsOrder>
};

interface tStepPayload {
  // eslint-disable-next-line flowtype/no-weak-types
  [key: string]: any
}


export type tStep = {
  +name: string,
  // eslint-disable-next-line flowtype/no-weak-types
  info: Object,
  status: tStepStatus,
  +type: tStepType, // check,collect,instruct,enable,...
  payload: tStepPayload, // 工步的数据
  data: tStepPayload, // 工步执行过程中生成的数据
  steps: tStepArray,
  // startTime: Date,
  // endTime: Date,
  times: Array<Date>,
  skippable: boolean, // 此工步是否可跳过
  undoable: boolean, // 是否可重做
  description: string
};

export type tStepType = 'check' | 'collect' | 'instruct' | 'enable';

export type tStepStatus =
  | STEP_STATUS.DOING
  | STEP_STATUS.ENTERING
  | STEP_STATUS.FAIL
  | STEP_STATUS.FINISHED
  | STEP_STATUS.LEAVING
  | STEP_STATUS.READY;

export type tOrderStatus = $Values<typeof ORDER_STATUS>;

export type tStepArray = Array<tStep>;
