const { events } = require('../constants');

const registerTypes = {
  RUSH_DATA: 'WS_RUSH_DATA'
};

const registerHandlers = {
  [registerTypes.RUSH_DATA]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          workcenter: '2'
        },
        type: registerTypes.RUSH_DATA
      },
      events.register
    );
  }
};

module.exports = {
  registerHandlers,
  registerTypes
};
