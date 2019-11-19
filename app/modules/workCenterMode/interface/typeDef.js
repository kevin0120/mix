// @flow

import { translation as trans } from '../../../components/NavBar/local';
import type { tCommonActionType } from '../../../common/type';

export type tWorkCenterMode = trans.normWorkCenterMode | trans.reworkWorkCenterMode;

export type tAction = tCommonActionType & { mode: tWorkCenterMode };


