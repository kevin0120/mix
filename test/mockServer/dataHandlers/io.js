const { events } = require('../constants');
const demoDevices = require('../demo/demoDevice');

const ioTypes = {
  IO_CONTACT: 'WS_IO_CONTACT',
  IO_CONTACT_IN: 'WS_IO_CONTACT_IN', // mock server only
  IO_CONTACT_OUT: 'WS_IO_CONTACT_OUT', // mock server only
  IO_SET: 'WS_IO_SET'
};

const ioHandlers = {
  [ioTypes.IO_CONTACT]: (data, reply) => {
    reply(
      {
        sn: data.sn || 0,
        type: ioTypes.IO_CONTACT,
        data: {
          src: 'controller',
          sn: demoDevices.demoController.sn,
          inputs: '00000000'
        }
      },
      events.io
    );
  },
  [ioTypes.IO_CONTACT_IN]: (data, reply) => {
    reply(
      {
        sn: data.sn || 0,
        type: ioTypes.IO_CONTACT,
        data: {
          src: 'controller',
          sn: demoDevices.demoIO.sn,
          inputs: '00000000'
        }
      },
      events.io
    );
  },
  [ioTypes.IO_CONTACT_OUT]: (data, reply) => {
    reply(
      {
        sn: data.sn || 0,
        type: ioTypes.IO_CONTACT,
        data: {
          src: 'controller',
          sn: demoDevices.demoIO.sn,
          outputs: '00000000'
        }
      },
      events.io
    );
  },
  [ioTypes.IO_SET]: (data, reply) => {
    reply(
      {
        sn: data.sn || 0,
        type: ioTypes.IO_SET,
        data: {
          result: 0
        }
      },
      events.io
    );
  }
};

module.exports = {
  ioHandlers,
  ioTypes
};
