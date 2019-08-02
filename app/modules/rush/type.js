// @flow

export type tWebSocketEvent =
  'controller'
  | 'result'
  | 'register'
  | 'selector'
  | 'job'
  | 'scanner'
  | 'io'
  | 'odoo'
  | 'maintenance'
  | 'tool'
  | 'reader';

/* eslint-disable flowtype/no-weak-types */
export type tRushWebSocketData = {
  +type: string,
  +data: { [key: string]: any }
};

export type tBarcode = {
  +id: string,
  +barcode: string
};


/* eslint-enable flowtype/no-weak-types */
