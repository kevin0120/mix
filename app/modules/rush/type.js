// @flow

export const WEBSOCKET_EVENTS = {
  controller: 'controller',
  result: 'result',
  register: 'register',
  selector: 'selector',
  job: 'job',
  scanner: 'scanner',
  io: 'io',
  odoo: 'odoo',
  maintenance: 'maintenance',
  tool: 'tool',
  reader: 'reader',
  tightening_device: 'tightening_device',
  reply: 'reply'
};

export type tWebSocketEvent = $Keys<typeof WEBSOCKET_EVENTS>;

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
