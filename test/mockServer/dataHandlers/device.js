const { events } = require('../constants');
const { demoDevices } = require('../demo/demoDevice');

const deviceType = Object.freeze({
  controller: 'controller',
  scanner: 'scanner',
  tool: 'tool',
  reader: 'reader',
  io: 'io'
});

const deviceTypes = {
  STATUS: 'WS_DEVICE_STATUS'
};

const deviceHandlers = {
  [deviceTypes.STATUS]: (data, reply) => {
    reply(
      {
        type: deviceTypes.STATUS,
        sn: data.sn,
        data: demoDevices
      },
      events.device
    );
  }
};

module.exports = {
  deviceType,
  deviceHandlers
};
