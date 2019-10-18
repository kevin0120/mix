// @flow

import type { Saga } from 'redux-saga';
import { WEBSOCKET_EVENTS } from './constants';
import type { tAction } from '../typeDef';

export type tWebSocketEvent = string;// $Values<typeof WEBSOCKET_EVENTS>;

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
  data: TData,
  sn: number
|};

export type tRushConnectionCallback = (...Array<any>)=>tAction<any, any>;
