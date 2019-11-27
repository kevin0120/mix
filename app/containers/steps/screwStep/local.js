import { makeLocalBundle, lng } from '../../../i18n';

export const screwStepNS = 'screwStepNS';

export const translation = {
  redoSpecScrewPointTitle: 'redoSpecScrewPointTitle',
  redoSpecScrewPointContent: 'redoSpecScrewPointContent',
  confirm: 'confirm',
  cancel: 'cancel'
};


const trans = [
  makeLocalBundle(lng.zh_CN, 'screwStepNS', {
    confirm: '确认返工拧紧点',
    cancel: '取消返工此拧紧点',
    redoSpecScrewPointTitle: '返工当前拧紧点',
    redoSpecScrewPointContent: '返工当前拧紧点:'
  }),
  makeLocalBundle(lng.en, 'screwStepNS', {
    confirm: 'Confirm Rework',
    cancel: 'Cancel',
    redoSpecScrewPointTitle: 'Rework Current Tightening Point',
    redoSpecScrewPointContent: 'Make Sure. Do You Want To Rework This Tightening Joint? Current Joint:'
  })
];

export default trans;
