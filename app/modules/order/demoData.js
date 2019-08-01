// @flow
import dataImg from '../../../resources/imgs/working-page-test-pic1.png';
import { ORDER_STATUS }from './model';

export const demoOrder = {
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: ORDER_STATUS.TODO,
  steps: [

    {
      name: 'screw',
      info: 'this is an screw step',
      type: 'screw',
      skippable: false,
      undoable: true,
      payload: {
        maxRetryTimes:3,
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
      name: 'skippable',
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
      name: 'skippable',
      info: 'this step is skippable',
      description: '扫描二维码或在输入框中输入',
      type: 'scanner',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    }
    // {
    //   name: 'undoable',
    //   info: 'this step is undoable',
    //   type: 'input',
    //   undoable: true,
    //   skippable: true,
    //
    //   payload: {
    //     label: 'address'
    //   }
    // },
    // {
    //   name: 'instruction',
    //   info: 'this is an instruction step',
    //   type: 'instruction',
    //   skippable: true,
    //
    //   payload: {
    //     instruction: '根据这段文字进行作业'
    //   }
    // },

    // {
    //   name: '扫码',
    //   info: 'this step collects a input',
    //   type: 'scanner',
    //   payload: {
    //     label: 'address'
    //   }
    // }
    // {
    //   name: '联系方式',
    //   info: 'this step collects 2 inputs',
    //   type:'parallel_input',
    //   steps: [
    //     [{
    //       name: '座机',
    //       info: 'this step collects a input',
    //       type: 'input', // check,collect,instruct,enable,...
    //       payload: {
    //         label: 'phone'
    //       }
    //     }, {
    //       name: '手机',
    //       info: 'this step collects another input',
    //       type: 'input', // check,collect,instruct,enable,...
    //       payload: {
    //         label: 'mobile'
    //       }
    //     }],
    //     {
    //       name: '电子邮箱',
    //       info: 'this step collects another input',
    //       type: 'input', // check,collect,instruct,enable,...
    //       payload: {
    //         label: 'email'
    //       }
    //     }
    //   ]
    // },
    // {
    //   name: '评价',
    //   info: 'this step collects 2 inputs',
    //   steps: [
    //     {
    //       name: '打分',
    //       info: 'this step collects a input',
    //       type: 'input', // check,collect,instruct,enable,...
    //       payload: {
    //         label: 'score'
    //       }
    //     }, {
    //       name: '评论',
    //       info: 'this step collects another input',
    //       type: 'input', // check,collect,instruct,enable,...
    //       payload: {
    //         label: 'comment'
    //       }
    //     }
    //   ]
    // }
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
      type:'input',

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
      type:'input',
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
