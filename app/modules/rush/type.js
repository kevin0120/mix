// @flow

import type { Saga } from 'redux-saga';
import { WEBSOCKET_EVENTS } from './constants';

export type tWebSocketEvent = $Values<typeof WEBSOCKET_EVENTS>;

/* eslint-disable flowtype/no-weak-types */
export type tRushWebSocketData = {
  +type: string,
  +data: { [key: string]: any } | Array<any>
};

export type tBarcode = {
  +id: string,
  +barcode: string
};

export type tRushHandler<TData> = (TData)=>void | Saga<void>;

export type rushHandlerMap<TTypes, TData> = {
  [key: TTypes]: tRushHandler<TData>
};

export type tRushData<TType, TData> = {|
  type: TType,
  data: TData
|};

/* eslint-enable flowtype/no-weak-types */
