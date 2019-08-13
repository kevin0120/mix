// @flow

// import type { tCommonActionType } from '../../common/type';

export const RUSH = {
  INIT: 'RUSH_INIT',
  STARTED: 'RUSH_STARTED',
  STOPPED: 'RUSH_STOPPED',
  NEW_RESULTS: 'RUSH_NEW_RESULTS',
  DATA: 'RUSH_DATA'
};

export default {
  // eslint-disable-next-line flowtype/no-weak-types
  data: (payload: { type: string, sn: number, data: Object }) => ({
    type: RUSH.DATA,
    payload
  })
};

// eslint-disable-next-line flowtype/no-weak-types
// export function NewResults(results: Array<any>): tCommonActionType {
//   return {
//     type: RUSH.NEW_RESULTS,
//     data: results
//   };
// }
