const { events } = require('../constants');

const scannerTypes = {
  STATUS: 'WS_SCANNER'
};

const scannerHandlers = {
  [scannerTypes.STATUS]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          sn: '3',
          barcode: 1
        }
      },
      events.scanner
    );
  }
};

module.exports = {
  scannerHandlers,
  scannerTypes
};
