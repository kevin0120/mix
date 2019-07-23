export const demoOrder = {
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: 'TODO',
  steps: [
    {
      name: 'skippable',
      info: 'this step is skippable',
      type: 'input',
      skippable: true,
      data: {},
      payload: {
        label: 'name'
      }
    },
    {
      name: 'undoable',
      info: 'this step is undoable',
      type: 'input',
      undoable: true,

      payload: {
        label: 'address'
      }
    },
    {
      name: '地址',
      info: 'this step collects a input',
      type: 'input',
      revocable: true,

      payload: {
        label: 'address'
      }
    }
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

export const demoOrderExcp = {
  name: 'fill info',
  info:
    'this is a demo order this is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo orderthis is a demo order',
  type: 'step',
  status: 'FAIL',
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
