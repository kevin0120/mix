const { events } = require('../constants');

const aiisTypes = {
  AIIS_STATUS: 'WS_AIIS_STATUS'
};

const aiisHandlers = {
  [aiisTypes.AIIS_STATUS]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          name: 'AIIS',
          status: 'online'
        }
      },
      events.aiis
    );
  }
};

module.exports = {
  aiisHandlers,
  aiisTypes
};
