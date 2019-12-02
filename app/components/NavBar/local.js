import { makeLocalBundle, lng } from '../../i18n';
import { workModes  } from '../../modules/workCenterMode/constants';
import { reworkDialogConstants  } from '../../modules/reworkPattern/constants';

export const navBarNs = 'navBar';

export const translation = {
  ...workModes,
  ...reworkDialogConstants
};


const trans = [
  makeLocalBundle(lng.zh_CN, navBarNs, {
    [translation.normWorkCenterMode]: '工作模式',
    [translation.reworkWorkCenterMode]: '返工模式',
    [translation.confirm]: '确认切换工作模式',
    [translation.cancel]: '取消',
    [translation.switchWorkCenterModeTitle]: '切换工位工作模式',
    [translation.switchWorkCenterModeContent]: '确认是否切换工作模式？当前模式为:'
  }),
  makeLocalBundle(lng.en, navBarNs, {
    [translation.normWorkCenterMode]: 'Normal Mode',
    [translation.reworkWorkCenterMode]: 'Rework Mode',
    [translation.confirm]: 'Confirm',
    [translation.cancel]: 'Cancel',
    [translation.switchWorkCenterModeTitle]: 'Switch WorkCenter Work Mode',
    [translation.switchWorkCenterModeContent]: 'Make Sure. Do You Want To Switch WorkCenter Mode? Current Work Center Mode:'
  })
];

export default trans;
