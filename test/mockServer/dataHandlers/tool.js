const { events } = require('../constants');

const toolTypes = {
  ENABLE: 'WS_TOOL_ENABLE',
  PSET: 'WS_TOOL_PSET'
};

const toolHandlers = {
  [toolTypes.ENABLE]: (data, reply) => {
    reply(
      {
        sn: data.sn,
        data: {
          result: 0
        },
        type: toolTypes.ENABLE
      },
      events.reply
    );
  },
  [toolTypes.PSET]: (data, reply) => {
    // eslint-disable-next-line camelcase
    const { sequence, total, tool_sn } = data.data;
    reply(
      {
        sn: data.sn,
        data: {
          result: 0
        },
        type: toolTypes.PSET
      },
      events.reply
    );
    setTimeout(() => {
      reply(
        {
          sn: 0,
          data: [
            {
              tool_sn,
              sequence,
              ti: 1,
              mi: 1,
              wi: 1,
              result: 'ok', // result
              batch: `${sequence}/${total}`
            }
          ]
        },
        events.result
      );
    }, 1000);
  }
};

module.exports = {
  toolHandlers,
  toolTypes
};
