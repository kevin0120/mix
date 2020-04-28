const { events } = require('../constants');

const exsysTypes = {
  EXSYS_STATUS: 'WS_EXSYS_STATUS'
};

const exsysHandlers = {
  [exsysTypes.EXSYS_STATUS]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          name: 'Mes',
          status: 'online'
        }
      },
      events.exsys
    );
  }
};

module.exports = {
  exsysHandlers,
  exsysTypes
};
