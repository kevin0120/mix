// @flow
import { put } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { orderActions } from './action';
import { ORDER_WS_TYPES } from './constants';
import type { tOrderWSTypes, tOrderRushData, tOrderListData, tOrder } from './interface/typeDef';
import type { rushHandlerMap } from '../rush/type';

const { getListSuccess } = orderActions;


// rush data handlers
const dataHandlers: rushHandlerMap<tOrderWSTypes, $PropertyType<tOrderRushData, 'data'>> = {

  // 工单列表
  * [ORDER_WS_TYPES.LIST](data: Array<tOrderListData>) {
    try {
      const list = data.map(d => ({
        id: d.id,
        desc: d.desc,
        name: d.name,
        image: d.image || '',
        status: d.status
      }));
      yield put(getListSuccess(list));
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.LIST' });
    }
  },
  // 工单详情
  * [ORDER_WS_TYPES.DETAIL](data: tOrder) {
    try {
      yield put(orderActions.getDetailSuccess(data));
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.DETAIL' });
    }
  },
  // 新工单
  * [ORDER_WS_TYPES.NEW](data: Array<tOrder>) {
    try {
      yield put(orderActions.newOrder(data));
    } catch (e) {
      CommonLog.lError(e,{ at: 'ORDER_WS_TYPES.NEW' });
    }
  }
};

export default dataHandlers;
