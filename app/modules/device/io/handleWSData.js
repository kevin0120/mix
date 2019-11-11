// @flow
import { call, all } from 'redux-saga/effects';
import { isNil } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../../common/utils';
import type { tIOWSDataContact, tIOWSMsgType } from './type';
import { IO_WS_TYPES, ioDirection } from './constants';
import { getDevice } from '../../deviceManager/devices';
import type { rushHandlerMap } from '../../rush/type';

// eslint-disable-next-line flowtype/no-weak-types
const handlers: rushHandlerMap<tIOWSMsgType, any> = {
  *[IO_WS_TYPES.CONTACT](data: { data: tIOWSDataContact }): Saga<void> {
    try {
      // TODO: test it
      const { sn, inputs, outputs } = data.data;
      if (isNil(sn)) {
        throw new Error(`io contact with invalid sn: ${sn}`);
      }
      const contacts = [];
      [inputs, outputs].forEach(contact => {
        if (isNil(contact)) {
          return;
        }
        let direction = ioDirection.input;
        if (contact === outputs) {
          direction = ioDirection.output;
        }
        contacts.push({
          sn,
          direction,
          contact
        });
      });

      const io = getDevice(sn);
      const dispatches = contacts.map(c => call(io.doDispatch, c));
      yield all(dispatches);
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  [IO_WS_TYPES.STATUS](data) {
    CommonLog.Info('WS_IO_STATUS', data);
  }
};

export default handlers;
