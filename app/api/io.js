// @flow
import { rushSendApi } from './rush';

type tIOSn = string;

export function ioSetApi(sn: tIOSn, index: number, status: 0 | 1): Promise<any> {
  return rushSendApi('WS_IO_SET', { sn, index, status });
}

export function ioContactApi(sn: tIOSn): Promise<any> {
  return rushSendApi('WS_IO_CONTACT', { sn });
}

export function ioStatusApi(sn: tIOSn): Promise<any> {
  return rushSendApi('WS_IO_STATUS', { sn });
}
