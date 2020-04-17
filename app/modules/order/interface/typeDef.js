// @flow
import { ORDER_STATUS, ORDER_WS_TYPES } from '../constants';
import type { IOrder } from './IOrder';
import type { tRushData } from '../../rush/type';
import type { tWorkableData } from '../../workable/typeDef';

export type 工单号 = string;
export type 人员列表 = Array<string>;
export type 工位号 = string;

export type tStockMove = {
  lot: string, // 批次号或者序列号
  product: string // 所选的产品名称
};

export type tOrderListData = {|
  id?: number,
  code?: 工单号,
  track_code: string,
  product_code: string,
  workcenter: 工位号,
  date_planned_start: string,
  date_planned_complete: string,
  status: tOrderStatus,
  product_type_image: string // 产成品图片
|};

export type tOrder = {
  ...tWorkableData,
  ...tOrderListData,
  payload: tOrderPayload
};

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

export type tProduct = {
  url: string,
  code: string
};

export type tWorksheet = {
  url: string,
  name: string,
  revision: string
};

export type tEnvironments = {
  text: string,
  test_type: string,
  code: string,
  sequence: number,
  desc: string
};

export type tComponent = {
  is_key: boolean,
  code: string
};

export type tResources = {
  equipments: Array<string>,
  users: Array<string>
};

export type tOperation = {
  code: string,
  resources: tResources,
  desc: string
};

export type tOrderPayload = {
  product: tProduct,
  worksheet: tWorksheet,
  environments: tEnvironments,
  components: Array<tComponent>,
  operation: tOperation
};

export type tWorkOnOrderConfig = Object;
