const { events } = require('../constants');
const { read, write, list } = require('../utils/storage');

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

const demoOrdersObj = {};

const getListInfo = (order) => {
  const { id, code, track_code, product_code, workcenter, date_planned_start, date_planned_complete, status } = order;
  return {
    id,
    code,
    track_code,
    product_code,
    workcenter,
    date_planned_start,
    date_planned_complete,
    status
  };
};

const orderNames = list('order');
Promise.all(
  orderNames.map(n => read('order', n))
).then((resp) => {
  resp.forEach(r => {
    const { name, data: order } = r;
    demoOrdersObj[name] = order;
  });
  const orders = Object.values(demoOrdersObj);
  orders.forEach((o, idx) => {
    orders[idx].id = idx;
  });
  let stepId = 0;
  orders.forEach((o) => {
    (o.steps || []).forEach((s, idx) => {
      // eslint-disable-next-line no-param-reassign
      o.steps[idx].id = stepId;
      stepId += 1;
    });
  });
}).catch(console.error);

const orderHandlers = {
  [orderTypes.LIST]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        type: orderTypes.LIST,
        data: Object.values(demoOrdersObj).map(getListInfo)
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
  // [orderTypes.DETAIL]: (data, reply) => {
  //   const order = demoOrderList.find(d => d.id === data.data.id);
  //   reply(
  //     {
  //       sn: data.sn,
  //       type: orderTypes.DETAIL,
  //       data: order
  //     },
  //     events.order
  //   );
  // },
  [orderTypes.ORDER_DETAIL_BY_CODE]: (data, reply) => {
    const order = Object.values(demoOrdersObj).find(d => d.code === data.data.code);
    reply(
      {
        sn: data.sn,
        type: data.type,
        data: {
          result: 0
        }
      },
      events.reply
    );
    reply(
      {
        sn: data.sn,
        type: orderTypes.DETAIL,
        data: order
      },
      events.order
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
    const { data: { id, status } } = data;
    console.log(Object.values(demoOrdersObj));
    const order = Object.values(demoOrdersObj).find(o => o.id === id);
    if (order) {
      order.status = status;
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
      return;
    }
    reply({
        sn: data.sn,
        type: orderTypes.UPDATE,
        data: {
          result: -1,
          msg: `cannot find order with id: ${id}`
        }
      },
      events.reply);
  },
  [orderTypes.ORDER_STEP_DATA_UPDATE]: (data, reply) => {
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
  }
  // [orderTypes.NEW]: (data, reply) => {
  //   reply(
  //     {
  //       sn: 0,
  //       type: orderTypes.NEW,
  //       data: [demoOrderJson]
  //     },
  //     events.order
  //   );
  // }
};

module.exports = {
  orderTypes,
  orderHandlers
};
