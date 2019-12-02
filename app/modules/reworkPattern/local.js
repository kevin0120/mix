import { makeLocalBundle, lng } from '../../i18n';
import { redoPointConstants } from './constants';

import {reworkNS} from './constants';

export const translation = {
  ...redoPointConstants
};


const trans = [
  makeLocalBundle(lng.zh_CN, reworkNS, {
    confirm: '确认返工',
    cancel: '取消返工',
    redoSpecScrewPointTitle: '返工当前拧紧点',
    redoSpecScrewPointTitleNoPoint: '返工当前工步',
    redoSpecScrewPointContent: '返工当前拧紧点:',
    redoSpecScrewPointContentNoPoint: '返工当前工步'
  }),
  makeLocalBundle(lng.en, reworkNS, {
    confirm: 'Confirm Rework',
    cancel: 'Cancel',
    redoSpecScrewPointTitle: 'Rework Current Tightening Point',
    redoSpecScrewPointTitleNoPoint: 'Rework Current Step',
    redoSpecScrewPointContent: 'Make Sure. Do You Want To Rework This Tightening Joint? Current Joint:',
    redoSpecScrewPointContentNoPoint: 'Make Sure. Do You Want To Rework This Step?'
  })
];

export default trans;
