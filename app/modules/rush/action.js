// @flow

import type {tCommonActionType} from "../../common/type"

export const RUSH = {
  INIT: 'RUSH_INIT',
  STARTED: 'RUSH_STARTED',
  STOPPED: 'RUSH_STOPPED',
  NEW_RESULTS: 'RUSH_NEW_RESULTS'
};


// eslint-disable-next-line flowtype/no-weak-types
export function NewResults(results: Array<any>): tCommonActionType {
  return {
    type: RUSH.NEW_RESULTS,
    data: results
  };
}
