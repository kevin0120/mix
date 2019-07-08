export const demoOrder = {
  name: 'fill info',
  info: 'this is a demo order',
  steps: [
    {
      name: '姓名',
      info: 'this step does one checking',
      type: 'input', // check,collect,instruct,enable,...
      payload: {
        label: 'name'
      }
    },
    {
      name: '地址',
      info: 'this step collects a input',
      type: 'input', // check,collect,instruct,enable,...
      payload: {
        label: 'address'
      }
    },
    {
      name: '联系方式',
      info: 'this step collects 2 inputs',
      steps: [
        [{
          name: '座机',
          info: 'this step collects a input',
          type: 'input', // check,collect,instruct,enable,...
          payload: {
            label: 'phone'
          }
        }, {
          name: '手机',
          info: 'this step collects another input',
          type: 'input', // check,collect,instruct,enable,...
          payload: {
            label: 'mobile'
          }
        }],
        {
          name: '电子邮箱',
          info: 'this step collects another input',
          type: 'input', // check,collect,instruct,enable,...
          payload: {
            label: 'email'
          }
        }
      ]
    },
    {
      name: '评价',
      info: 'this step collects 2 inputs',
      steps: [
        {
          name: '打分',
          info: 'this step collects a input',
          type: 'input', // check,collect,instruct,enable,...
          payload: {
            label: 'score'
          }
        }, {
          name: '评论',
          info: 'this step collects another input',
          type: 'input', // check,collect,instruct,enable,...
          payload: {
            label: 'comment'
          }
        }
      ]
    }

  ]
};
