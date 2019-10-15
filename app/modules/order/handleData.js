// @flow
import { put, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';
import { orderActions } from './action';
import { ORDER_WS_TYPES } from './constants';
import type { tOrderWSTypes, tOrderRushData, tOrderListData, tOrder } from './interface/typeDef';
import type { rushHandlerMap } from '../rush/type';

const { getListSuccess } = orderActions;


// rush data handlers
const dataHandlers: rushHandlerMap<tOrderWSTypes, $PropertyType<tOrderRushData, 'data'>> = {
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
  * [ORDER_WS_TYPES.DETAIL](data: tOrder) {
    try {
      yield put(orderActions.getDetailSuccess(data));
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.DETAIL' });
    }
  },
  * [ORDER_WS_TYPES.NEW](data: Array<tOrder>) {
    try {
      yield put(orderActions.newOrder(data));
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};

export default function* orderData(rushData: tOrderRushData): Saga<void> {
  try {
    const { type, data } = rushData;
    yield call(dataHandlers[type], data);
  } catch (e) {
    CommonLog.lError(e);
  }
}
