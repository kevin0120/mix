import { put, call } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import { orderActions } from './action';
import { ORDER_STATUS } from './model';

const { getListSuccess } = orderActions;

export const ORDER_WS_TYPES = {
  LIST: 'WS_ORDER_LIST',
  DETAIL: 'WS_ORDER_DETAIL',
  UPDATE: 'WS_ORDER_UPDATE',
  STEP_UPDATE: 'WS_ORDER_STEP_UPDATE'
};


const dataHandlers = {
  * [ORDER_WS_TYPES.LIST](data) {
    try {
      const list = data.map(d => ({
        id: d.id,
        info: d.desc,
        name: d.name,
        image: '',
        status: d.status
      }));
      yield put(getListSuccess(list));
    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.LIST' });
    }
  },
  * [ORDER_WS_TYPES.DETAIL](data) {
    try {
      if (data.result === 0) {
        yield put(orderActions.getDetailSuccess(data));
      } else {
        yield put(orderActions.getDetailSuccess(data));
      }

    } catch (e) {
      CommonLog.lError(e, { at: 'ORDER_WS_TYPES.DETAIL' });
    }
  }
};

export default function* orderData(rushData) {
  try {
    const { type, data } = rushData;
    yield call(dataHandlers[type], data);
  } catch (e) {
    CommonLog.lError(e);
  }
}
