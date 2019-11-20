import { makeLocalBundle, lng } from '../../../i18n';

export const screwStepNS = 'screwStepNS';

export const translation = {
  redoSpecScrewPointTitle: 'switchWorkCenterModeTitle',
  redoSpecScrewPointContent: 'switchWorkCenterModeContent',
  confirm: 'confirm',
  cancel: 'cancel'
};


const trans = [
  makeLocalBundle(lng.zh_CN, 'screwStepNS', {
    confirm: '确认切换工作模式',
    cancel: '取消',
    redoSpecScrewPointTitle: '切换工位工作模式',
    redoSpecScrewPointContent: '确认是否切换工作模式？当前模式为:'
  }),
  makeLocalBundle(lng.en, 'screwStepNS', {
    confirm: 'Confirm',
    cancel: 'Cancel',
    redoSpecScrewPointTitle: 'Switch WorkCenter Work Mode',
    redoSpecScrewPointContent: 'Make Sure. Do You Want To Switch WorkCenter Mode? Current Work Center Mode:'
  })
];

export default trans;
