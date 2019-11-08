const demoTool = {
  sn: 'xx1233',
  type: 'tool',
  name: 'tool1',
  status: 'online'
};
const demoTool2 = {
  sn: '2',
  type: 'tool',
  name: 'tool1',
  status: 'online'
};
const demoController = {
  sn: 'c1',
  type: 'controller',
  name: 'controller1',
  status: 'online',
  children: ['2', 'xx1233'],
  config: {
    input_num: 8,
    output_num: 8
  }
};

const demoScanner = {
  sn: '3',
  type: 'scanner',
  status: 'online'
};

const demoReader = {
  sn: '4',
  type: 'reader',
  status: 'online'
};

module.exports = {
  demoDevices: [demoTool, demoTool2, demoScanner, demoReader, demoController]
};
