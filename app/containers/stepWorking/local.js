import { makeLocalBundle, lng } from '../../i18n';

export const stepWorkingNS='stepWorking';

export const translation = {
  skip: 'skip',
  undo: 'undo',
  submit: 'submit',
  viewModel: 'viewModel',
  cancel: 'cancel',
  pending: 'pending',
  continueDoing: 'continueDoing',
  view: 'view',
  name:'name',
  desc: 'desc',
  action: 'action',
};


const trans = [
    makeLocalBundle(lng.zh_CN, 'stepWorking', {
      skip: '跳过',
      undo: '取消上次动作',
      submit: '提交',
      finish: '完成',
      viewModel: '查看模型',
      cancel: '取消订单',
      pending: '阻塞订单',
      continueDoing: '继续订单',
      view: '查看',
      name:'名称',
      desc: '描述',
      action: '操作',
      notSelected:'未选中工单'
    }),
    makeLocalBundle(lng.en, 'stepWorking', {
      skip: 'skip',
      undo: 'undo',
      submit: 'submit',
      finish: 'finish',
      viewModel: 'view model',
      cancel: 'cancel',
      pending: 'pending',
      continueDoing: 'continue',
      view: 'view',
      name:'name',
      desc: 'description',
      action: 'action',
      notSelected:'Order Not Selected'
    })
  ];

export default trans;
