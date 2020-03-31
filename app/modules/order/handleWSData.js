// @flow
import { put, select } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { orderActions } from './action';
import { ORDER_WS_TYPES } from './constants';
import type {
  tOrderWSTypes,
  tOrderRushData,
  tOrderListData,
  tOrder
} from './interface/typeDef';
import type { rushHandlerMap } from '../rush/type';
import OrderMixin from './Order';
import Workable from '../workable/Workable';
import NotifierActions from '../Notifier/action';

// rush data handlers
const dataHandlers: rushHandlerMap<tOrderWSTypes,
  $PropertyType<tOrderRushData, 'data'>> = {
  // 工单列表
  * [ORDER_WS_TYPES.LIST](data: Array<tOrderListData>) {
    try {
      const list = (data || []).map(d => ({
        ...d
      }));
      const orderState = yield select(s => s.order);
      // get exist orders, orders not in the new list will be removed!!
      let newList =
        orderState &&
        orderState.list.filter(o => !!list.find(newO => o.code === newO.code));
      // update order data
      newList.forEach(o => {
        const orderData = list.find(newO => o.code === newO.code);
        o.update(orderData);
      });
      // make new orders
      newList = newList.concat(
        list
          .filter(newO => !newList.find(o => o.code === newO.code))
          .map(oD => new (OrderMixin(Workable))(oD))
      );
      yield put(orderActions.newList(newList));
      yield put(orderActions.getListSuccess());
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.LIST' });
    }
  },
  // 工单详情
  * [ORDER_WS_TYPES.DETAIL](data: tOrder) {
    try {
      const orderState = yield select(s => s.order);
      const newList = [...orderState.list];
      if (!data || !data.code) {
        yield put(orderActions.getDetailSuccess());
        return;
      }
      const newOrder = newList.find(o => o.code === data.code);
      if (newOrder) {
        newOrder.update(data);
      } else {
        newList.push(new (OrderMixin(Workable))(data));
      }
      yield put(orderActions.newList(newList));
      yield put(orderActions.getDetailSuccess());
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.DETAIL' });
    }
  },
  * [ORDER_WS_TYPES.ORDER_DETAIL_BY_CODE](data: tOrder) {
    try {
      const orderState = yield select(s => s.order);
      const newList = [...orderState.list];
      const newOrder = newList.find(o => o.code === data.code);
      if (newOrder) {
        newOrder.update(data);
      } else {
        newList.push(new (OrderMixin(Workable))(data));
      }
      yield put(orderActions.newList(newList));
      yield put(orderActions.getDetailSuccess());
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.DETAIL' });
    }
  },
  // 新工单
  * [ORDER_WS_TYPES.NEW](data: Array<tOrder>) {
    try {
      const orderState = yield select(s => s.order);
      // get exist orders
      let newList = orderState.list;

      // update order data
      newList.forEach(o => {
        const orderData = data.find(newO => o.code === newO.code);
        if (orderData) {
          o.update(orderData);
        }
      });

      // make new orders
      newList = newList.concat(data
        .filter(newO => !newList.find(o => o.code === newO.code))
        .map(oD => new (OrderMixin(Workable))(oD))
      );

      yield put(orderActions.newList(newList));
      yield put(NotifierActions.enqueueSnackbar('Info', '收到新工单'));
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.NEW' });
    }
  }
};

export default dataHandlers;
