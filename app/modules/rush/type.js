// @flow

import { WEBSOCKET_EVENTS } from './constants';

export type tWebSocketEvent = $Keys<typeof WEBSOCKET_EVENTS>;

/* eslint-disable flowtype/no-weak-types */
export type tRushWebSocketData = {
  +type: string,
  +data: { [key: string]: any } | Array<any>
};

export type tBarcode = {
  +id: string,
  +barcode: string
};

/* eslint-enable flowtype/no-weak-types */
