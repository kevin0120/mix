// @flow
export type tCommonActionType = {
  +type: string
};

export type tTheme = {
  spacing: number => number,
  transitions: {
    // eslint-disable-next-line flowtype/no-weak-types
    [keys: string]: any
  },
  breakpoints: {
    // eslint-disable-next-line flowtype/no-weak-types
    [keys: string]: any
  },
  // eslint-disable-next-line flowtype/no-weak-types
  spacing: any,
  palette: {
    // eslint-disable-next-line flowtype/no-weak-types
    [keys: string]: any
  }
};
