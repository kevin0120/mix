// @flow
import dataImg from '../../../resources/imgs/working-page-test-pic1.png';
import dataImg2 from '../../../resources/imgs/working-page-test-pic2.png';
import { ORDER_STATUS } from './model';

export const demoOrder = {
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.TODO,
  steps: [
    {
      name: '拧紧作业指导 pset',
      info: 'this is an screw step',
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
    }, {
      name: '拧紧作业指导 job',
      info: 'this is an screw step',
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
      name: '扫码',
      info: 'this step is skippable',
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: '获取物料',
      info: 'fetch the following material',
      description: '获取下列物料',
      type: 'material',
      skippable: true,
      data: {},
      payload: {}
    },
    {
      name: '扫码',
      info: 'this step is skippable',
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
  name: 'Long Demo Order',
  info: 'this demo order is long!',
  type: 'step',
  status: ORDER_STATUS.TODO,
  steps: [
    {
      name: '拧紧作业指导',
      info: 'this is an screw step',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
      info: 'this step is skippable',
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
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.CANCEL,
  steps: [
    {
      name: '姓名',
      info: 'this step does one checking',
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
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.PENDING,
  steps: [
    {
      name: '姓名',
      info: 'this step does one checking',
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
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.DONE,
  steps: [
    {
      name: '姓名',
      info: 'this step does one checking',
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

export const demoOrder2 = {
  name: '扫码触发 拧紧作业',
  type: 'screw',
  steps: [
    {
      name: '扫码',
      steps: [
        [
          {
            name: 'carID',
            type: 'scanner',
            payload: {}
          },
          {
            name: 'CarType',
            type: 'scanner',
            payload: {}
          }
        ]
      ]
    },
    {
      name: '获取作业',
      type: 'getOperation',
      payload: {
        keys: ['carID', 'carType']
      }
    },
    {
      name: '进行作业',
      type: 'doScrew',
      payload: {}
    },
    {
      name: '显示结果',
      type: 'showResult',
      payload: {
        keys: []
      }
    }
  ]
};
