// @flow
import { call } from 'redux-saga/effects';
import { CommonLog } from '../../../../common/utils';
import type { tIOWSDataContact, tIOWSMsgType } from './type';
import { IO_WS_TYPES } from './constants';
import { getDevice } from '../index';
import type { rushHandlerMap } from '../../../rush/type';


// eslint-disable-next-line flowtype/no-weak-types
const handlers: rushHandlerMap<tIOWSMsgType, any> = {
  * [IO_WS_TYPES.CONTACT](data: { data: tIOWSDataContact }) {
    try {
      const io = getDevice(data.data.sn);
      yield call(io.doDispatch, {
        sn: data.data.sn,
        direction: data.data.type,
        contact: data.data.contact
      });
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  [IO_WS_TYPES.STATUS](data) {
    CommonLog.Info('WS_IO_STATUS', data);
  },
};

export default handlers;