// @flow

import type { tCommonActionType } from '../../../common/type';


export type tReworkDetailSpecialScrewType = 'REWORK_SCREW_S_POINT'; // 返工一个具体拧紧点

export type tReworkDetailScrewType = 'REWORK_SCREW_STEP'; // 返工整个拧紧工步


export type tReworkDetailType = tReworkDetailSpecialScrewType | tReworkDetailScrewType;

export type tAction = tCommonActionType & { dType: tReworkDetailType, extra: any }; // extra为额外信息， 如具体拧紧点为一个具体的拧紧点信息


