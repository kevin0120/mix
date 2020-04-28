import { rushSendApi } from './rush';

export function deviceStatusApi(): Promise<any> {
  return rushSendApi('WS_DEVICE_STATUS');
}
