// @flow
import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';

type tIOSn = string;

export function ioSetApi(sn: tIOSn, index: number, status: string) {
  try {
    return rushSendApi('WS_IO_SET', {
      sn,
      index,
      status
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioSetApi'
    });
  }
}

export function ioContactApi(sn: tIOSn) {
  try {
    return rushSendApi('WS_IO_CONTACT', {
      sn
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioContact'
    });
  }
}

export function ioStatusApi(sn: tIOSn) {
  try {
    return rushSendApi('WS_IO_STATUS', {
      sn
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: 'ioStatusApi'
    });
  }
}
