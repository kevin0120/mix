const { events } = require('../constants');
const demoDevices = require('../demo/demoDevice');

const resultTypes = {
  RESULT: 'WS_TIGHTENING_RESULT'
};

const resultHandlers = {
  [resultTypes.RESULT]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: [
          {
            tool_sn: demoDevices.demoTool.sn,
            sequence: 1,
            ti: 1,
            mi: 1,
            wi: 1,
            result: 'ok', // result
            batch: '1/1'
          },
          {
            tool_sn: demoDevices.demoTool.sn,
            sequence: 2,
            ti: 1,
            mi: 1,
            wi: 1,
            result: 'ok', // result
            batch: '1/1'
          },
          {
            tool_sn: demoDevices.demoTool.sn,
            sequence: 3,
            ti: 1,
            mi: 1,
            wi: 1,
            result: 'ok', // result
            batch: '1/1'
          }
        ]
      },
      events.result
    );
  }
};

module.exports = {
  resultHandlers,
  resultTypes
};
