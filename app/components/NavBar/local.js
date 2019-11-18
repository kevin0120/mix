import { makeLocalBundle, lng } from '../../i18n';

export const navBarNs = 'navBar';

export const translation = {
  normWorkCenterMode: 'normWorkCenterMode',
  reworkWorkCenterMode: 'reworkWorkCenterMode',
  switchWorkCenterModeTitle: 'switchWorkCenterModeTitle',
  switchWorkCenterModeContent: 'switchWorkCenterModeContent',
  confirm: 'confirm',
  cancel: 'cancel'
};


const trans = [
  makeLocalBundle(lng.zh_CN, 'navBar', {
    normWorkCenterMode: '工作模式',
    reworkWorkCenterMode: '返工模式',
    confirm: '确认切换工作模式',
    cancel: '取消',
    switchWorkCenterModeTitle: '切换工位工作模式',
    switchWorkCenterModeContent: '确认是否切换工作模式？'
  }),
  makeLocalBundle(lng.en, 'navBar', {
    normWorkCenterMode: 'Normal Mode',
    reworkWorkCenterMode: 'Rework Mode',
    confirm: 'Confirm',
    cancel: 'Cancel',
    switchWorkCenterModeTitle: 'Switch WorkCenter Work Mode',
    switchWorkCenterModeContent: 'Make Sure. Do You Want To Switch WorkCenter Mode'
  })
];

export default trans;
