// @flow
import { ORDER_STATUS, ORDER_WS_TYPES } from '../constants';
import type { IOrder } from './IOrder';
import type { tStep } from '../../step/interface/typeDef';
import type { tRushData } from '../../rush/type';

export type tStockMove = {
  lot: string, // 批次号或者序列号
  product: string // 所选的产品名称
};

export type tOrderListData = {|
  id: number,
  desc: string, // 工单描述
  name: string, // 工单名称（工单号）
  image: string,
  status: tOrderStatus
|};

export type tOrder = {|
  ...tOrderListData,
  canRework: boolean, // 是否能够返工
  incomingProducts: Array<tStockMove>,
  finishedProducts: Array<tStockMove>, //
  steps: Array<tStep>, // 工步
  status: tOrderStatus, // 工单状态
  plannedDateTime: string, // 计划时间
  workingIndex: ?number // 正在执行的工步索引
|};

export type tOrderStepIdx = number;

export type tOrderState = {
  viewingOrder: IOrder | null,
  workingOrder: IOrder | null,
  viewingIndex: tOrderStepIdx,
  list: Array<IOrder>
};

export type tOrderActionTypes = string;

export type tOrderStatus = $Values<typeof ORDER_STATUS>;
export type tOrderWSTypes = $Values<typeof ORDER_WS_TYPES>;

export type tOrderRushData = tRushData<tOrderWSTypes, any>;
