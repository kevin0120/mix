import { lng, makeLocalBundle } from '../../i18n';
import { stepTypeKeys } from '../../modules/step/constants';

export const stepWorkingNS = 'stepWorking';

export const translation = {
  skip: 'skip',
  undo: 'undo',
  submit: 'submit',
  viewModel: 'viewModel',
  viewFile: 'viewFile',
  cancel: 'cancel',
  pending: 'pending',
  continueDoing: 'continueDoing',
  view: 'view',
  name: 'name',
  desc: 'desc',
  action: 'action',
  reportFinish:'reportFinish',
  ...stepTypeKeys
};


const trans = [
  makeLocalBundle(lng.zh_CN, 'stepWorking', {
    skip: '跳过',
    undo: '取消上次动作',
    submit: '提交',
    finish: '完成',
    viewModel: '查看模型',
    viewFile: '查看工艺',
    cancel: '取消订单',
    pending: '阻塞订单',
    continueDoing: '继续订单',
    view: '查看',
    name: '名称',
    desc: '描述',
    action: '操作',
    notSelected: '未选中工单',
    [translation.material]: '物料',
    [translation.scanner]: '扫码',
    [translation.screw]: '拧紧',
    [translation.input]: '输入',
    [translation.instruction]: '指示',
    [translation.text]: '文字',
    [translation.passFail]: '检测',
    [translation.measure]: '测量',
    [translation.video]: '视频',
    [translation.reportFinish]: '报工'
  }),
  makeLocalBundle(lng.en, 'stepWorking', {
    skip: 'skip',
    undo: 'undo',
    submit: 'submit',
    finish: 'finish',
    viewModel: 'view model',
    viewFile: 'view file',
    cancel: 'cancel',
    pending: 'pending',
    continueDoing: 'continue',
    view: 'view',
    name: 'name',
    desc: 'description',
    action: 'action',
    notSelected: 'Order Not Selected',
    [translation.material]: 'material',
    [translation.scanner]: 'scanner',
    [translation.screw]: 'tightening',
    [translation.input]: 'input',
    [translation.instruction]: 'instruction',
    [translation.text]: 'text',
    [translation.passFail]: 'pass/fail',
    [translation.measure]: 'measure',
    [translation.video]: 'video',
    [translation.reportFinish]: 'Report finish'
  })
];

export default trans;
