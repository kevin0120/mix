const { events } = require('../constants');
const { demoOrder, demoOrderJson } = require('../demo/demoOrder');
const orderDetailData = require('../demo/orderDetailData.json');

const orderTypes = Object.freeze({
  LIST: 'WS_ORDER_LIST',
  DETAIL: 'WS_ORDER_DETAIL',
  UPDATE: 'WS_ORDER_UPDATE',
  STEP_UPDATE: 'WS_ORDER_STEP_UPDATE',
  NEW: 'WS_NEW_ORDER',
  START_REQUEST: 'WS_ORDER_START_REQUEST',
  FINISH_REQUEST: 'WS_ORDER_FINISH_REQUEST',
  STEP_DATA: 'STEP_DATA',
  ORDER_DETAIL_BY_CODE: 'WS_ORDER_DETAIL_BY_CODE',
  ORDER_STEP_DATA_UPDATE: 'WS_ORDER_STEP_DATA_UPDATE'
});

const orderHandlers = {
  [orderTypes.LIST]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.LIST,
        data: [
          // demoOrder
          {
            id: 1,
            code: 'MO0002',
            track_code: 'track_code for order list',
            product_code: 'demo_product_code',
            workcenter: 'workcenter',
            date_planned_start: '2019-10-16T03:20:30Z',
            date_planned_complete: 'date_planned_complete',
            product_type_image: 'this is an image',
            status: 'todo'
          },
          {
            id: 2,
            code: 'MO0003',
            track_code: 'track_code for order list',
            product_code: 'product_code',
            workcenter: 'workcenter',
            date_planned_start: '2019-10-16T03:25:30Z',
            date_planned_complete: 'date_planned_complete',
            status: 'todo'
          }
        ]
      },
      events.order
    );
  },
  [orderTypes.START_REQUEST]: (data, reply) => {
    reply({
      data: {
        result: 0
      },
      sn: data.sn,
      type: orderTypes.START_REQUEST
    });
  },
  [orderTypes.FINISH_REQUEST]: (data, reply) => {
    reply({
      data: {
        result: 0
      },
      sn: data.sn,
      type: orderTypes.FINISH_REQUEST
    });
  },
  [orderTypes.DETAIL]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.DETAIL,
        data: orderDetailData
      },
      events.order
    );
  },
  [orderTypes.ORDER_DETAIL_BY_CODE]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.DETAIL,
        data: orderDetailData
      },
      events.reply
    );
  },
  [orderTypes.STEP_UPDATE]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.STEP_UPDATE,
        data: {
          result: 0
        }
      },
      events.reply
    );
  },
  [orderTypes.ORDER_STEP_DATA_UPDATE]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.STEP_UPDATE,
        data: {
          result: 0
        }
      },
      events.reply
    );
  },
  [orderTypes.UPDATE]: (data, reply) => {
    demoOrder.status = data.status;
    reply(
      {
        sn: data.sn,
        type: orderTypes.UPDATE,
        data: {
          result: 0
        }
      },
      events.reply
    );
  },
  [orderTypes.NEW]: (data, reply) => {
    reply(
      {
        sn: 0,
        type: orderTypes.NEW,
        data: [demoOrderJson]
      },
      events.order
    );
  }
};

module.exports = {
  orderTypes,
  orderHandlers
};
