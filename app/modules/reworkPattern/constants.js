// @flow

// action types
export const REWORK_PATTERN = {
  TRY_REWORK: 'REWORK_PATTERN_TRY_REWORK',
  DO_REWORK: 'REWORK_PATTERN_DO_REWORK',
  CANCEL_REWORK: 'REWORK_PATTERN_CANCEL_REWORK'
};

export const reworkDialogConstants = {
  confirm: 'confirm',
  cancel: 'cancel',
  switchWorkCenterModeTitle: 'switchWorkCenterModeTitle',
  switchWorkCenterModeContent: 'switchWorkCenterModeContent'
};

export const redoPointConstants = {
  confirm: 'confirm',
  cancel: 'cancel',
  redoSpecScrewPointTitle: 'redoSpecScrewPointTitle',
  redoSpecScrewPointContent: 'redoSpecScrewPointContent'
};

export const reworkNS = 'reworkPattern';
