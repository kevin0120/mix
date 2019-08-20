// @flow
import dataImg from '../../../resources/imgs/workingImg.png';
import dataImg2 from '../../../resources/imgs/working-page-test-pic2.png';
import screwImg from '../../../resources/imgs/screw.jpg';
import { ORDER_STATUS } from './model';

// console.log(screwImg);
export const demoOrder = {
  id:'demoOrder1',
  name: 'demo Order 1',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.TODO,
  steps: [
    {
      id:'s1',

      name: '获取物料',
      info: {
        time: '00:02:00'
      },
      description: '获取下列物料',
      type: 'material',
      skippable: true,
      data: {},
      payload: {
        items: [
          {
            name: 'test1',
            index: 0
          }, {
            name: 'test2',
            index: 2
          }
        ],
        ioSN: '1',
        confirmIO: 4
      }
    },
    {
      id:'s2',

      name: '拧紧作业指导 pset',
      info: {
        time: '00:02:00'
      },
      type: 'screw',
      skippable: false,
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
            toolSN: 'demo'
          }, {
            sequence: 2,
            group_sequence: 2,
            x: 20,
            y: 20,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: 'demo'

          }, {
            sequence: 3,
            group_sequence: 3,

            x: 30,
            y: 30,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: 'demo'

          }, {
            sequence: 4,
            group_sequence: 3,
            x: 40,
            y: 40,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: 'demo'

          }, {
            sequence: 5,
            group_sequence: 4,
            x: 50,
            y: 50,
            maxRetryTimes: 3,
            pset: 1,
            toolSN: 'demo'

          }
        ],
        image: dataImg
      }
    },
    {
      id:'s3',
      name: '拧紧作业指导 job',
      info: {
        time: '00:02:00'
      },
      type: 'screw',
      skippable: false,
      undoable: true,
      payload: {
        controllerMode: 'job',
        jobID: 1,
        points: [
          {
            toolSN: 'demo2',
            sequence: 1,
            group_sequence: 1,
            x: 10,
            y: 10,
            maxRetryTimes: 3
          }, {
            toolSN: 'demo2',
            sequence: 2,
            group_sequence: 2,
            x: 20,
            y: 20,
            maxRetryTimes: 3
          }, {
            toolSN: 'demo2',
            sequence: 3,
            group_sequence: 3,
            x: 30,
            y: 30,
            maxRetryTimes: 3
          }, {
            toolSN: 'demo2',
            sequence: 4,
            group_sequence: 4,
            x: 40,
            y: 40,
            maxRetryTimes: 3
          }, {
            toolSN: 'demo2',
            sequence: 5,
            group_sequence: 5,
            x: 50,
            y: 50,
            maxRetryTimes: 3
          }
        ],
        image: dataImg2
      }
    },
    {
      id:'s4',
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      id:'s5',
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
  ]
};

export const demoOrderLong = {
  id:'demoOrderLong',

  name: 'Long Demo Order',
  info: 'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.TODO,
  steps: [
    {
      name: '拧紧作业指导',
      info: {
        time: '00:02:00'
      },
      type: 'screw',
      skippable: false,
      undoable: true,
      payload: {
        maxRetryTimes: 3,
        points: [
          {
            id: 1,
            x: 10,
            y: 10,
            status: 'waiting'
          }, {
            id: 2,
            x: 20,
            y: 20,
            status: 'waiting'
          }, {
            id: 3,
            x: 30,
            y: 30,
            status: 'waiting'
          }, {
            id: 4,
            x: 40,
            y: 40,
            status: 'waiting'
          }, {
            id: 5,
            x: 50,
            y: 50,
            status: 'waiting'
          }
        ],
        image: dataImg
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '扫码',
      info: {
        time: '00:02:00'
      },
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
  ]
};

export const demoOrderCancel = {
  id:'demoOrderCancel',
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.CANCEL,
  steps: [
    {
      name: '姓名',
      info: {
        time: '00:02:00'
      },
      type: 'input',
      skippable: true,
      revocable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
  ]
};

export const demoOrderPending = {
  id:'demoOrderPending',
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.PENDING,
  steps: [
    {
      name: '姓名',
      info: {
        time: '00:02:00'
      },
      type: 'input',

      skippable: true,
      revocable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
  ]
};

export const demoOrderDone = {
  id:'demoOrderDone',
  name: 'fill info',
  info: 'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.DONE,
  steps: [
    {
      name: '姓名',
      info: {
        time: '00:02:00'
      },
      type: 'input',
      skippable: true,
      revocable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
  ]
};

