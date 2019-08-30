import { CommonLog } from '../common/utils';
import { rushSendApi } from './rush';

export function deviceStatusApi() {
  try {
    return rushSendApi('WS_DEVICE_STATUS');
  } catch (e) {
    CommonLog.lError(e, {
      at: 'deviceStatusApi'
    });
  }
}
