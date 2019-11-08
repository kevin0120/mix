const dataImg = ''; // require('../res/working-page-test-pic1.png');
const dataImg2 = ''; // require('../res/working-page-test-pic2.png');
// 工单状态
const ORDER_STATUS = Object.freeze({
  TODO: 'todo',
  WIP: 'wip',
  DONE: 'done',
  CANCEL: 'cancel',
  PENDING: 'pending'
});

const demoOrderJson = require('./demoOrder.json');

const demoOrder = {
  id: 1,
  name: 'demo Order 1',
  desc: 'a:b\t\tb:c\t\tc:d',
  canRework: true,
  image: '',
  status: ORDER_STATUS.TODO,
  steps: [
    {
      id: 1,

      name: '获取物料',
      info: {
        time: '00:02:00'
      },
      description: '获取下列物料',
      type: 'material',
      skippable: true,
      undoable: true,
      payload: {
        items: [
          {
            name: 'test material 1',
            image: '',
            in: {
              sn: 'io1',
              index: 0
            },
            out: {
              sn: 'io2',
              index: 0
            }
          }
        ],
        confirm: {
          sn: 'io1',
          index: 0
        }
      }
    },
    {
      id: 2,

      name: '拧紧作业指导 pset',
      info: {
        time: '00:02:00'
      },
      type: 'screw',
      skippable: true,
      undoable: true,
      payload: {
        controllerMode: 'pset',
        points: [
          {
            sequence: 1,
            group_sequence: 1,
            x: 10,
            y: 10,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: '1',
            is_key: true
          },
          {
            sequence: 2,
            group_sequence: 2,
            x: 20,
            y: 20,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: '1',
            is_key: true
          },
          {
            sequence: 3,
            group_sequence: 3,

            x: 30,
            y: 30,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: '2',
            is_key: false
          },
          {
            sequence: 4,
            group_sequence: 3,
            x: 40,
            y: 40,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: '1',
            is_key: true
          },
          {
            sequence: 5,
            group_sequence: 4,
            x: 50,
            y: 50,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: '1',
            is_key: true
          }
        ],
        image: dataImg
      }
    },
    {
      id: 3,
      name: '拧紧作业指导 job',
      info: {
        time: '00:02:00'
      },
      type: 'screw',
      skippable: true,
      undoable: true,
      payload: {
        controllerMode: 'job',
        jobID: 1,
        points: [
          {
            toolSN: '2',
            sequence: 1,
            group_sequence: 1,
            minPassCount: 0,
            x: 10,
            y: 10,
            maxRetryTimes: 3
          },
          {
            toolSN: '2',
            sequence: 2,
            group_sequence: 2,
            minPassCount: 0,
            x: 20,
            y: 20,
            maxRetryTimes: 3
          },
          {
            toolSN: '2',
            sequence: 3,
            group_sequence: 3,
            minPassCount: 0,
            x: 30,
            y: 30,
            maxRetryTimes: 3
          },
          {
            toolSN: '2',
            sequence: 4,
            group_sequence: 4,
            minPassCount: 0,
            x: 40,
            y: 40,
            maxRetryTimes: 3
          },
          {
            toolSN: '2',
            sequence: 5,
            group_sequence: 5,
            minPassCount: 0,
            x: 50,
            y: 50,
            maxRetryTimes: 3
          }
        ],
        image: dataImg2
      }
    },
    {
      id: 4,
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      undoable: true,
      payload: {
        label: 'name'
      }
    },
    {
      id: 6,
      name: '指导',
      info: {
        time: '00:02:00'
      },
      description: 'pdf',
      type: 'instruction',
      skippable: true,
      undoable: true,
      payload: {
        url: '\\\\192.168.74.31\\share\\D0000177249.pdf',
        page: 0
      }
    },
    {
      id: 5,
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      undoable: true,
      payload: {
        label: 'name'
      }
    }
  ]
};

module.exports = {
  demoOrder,
  demoOrderJson
};
