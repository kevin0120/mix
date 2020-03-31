// @flow

import type {tWorkCenterMode} from  './interface/typeDef'

export const sGetWorkCenterMode = (state: Object): tWorkCenterMode =>
  state?.workCenterMode;
