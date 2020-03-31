const { events } = require('../constants');

const readerTypes = {
  DATA: 'WS_READER'
};

const readerHandlers = {
  [readerTypes.DATA]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          id: '4',
          uid: 1
        }
      },
      events.reader
    );
  }
};

module.exports = {
  readerHandlers,
  readerTypes
};
