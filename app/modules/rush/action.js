export const RUSH = {
  INIT: 'RUSH_INIT',
  STARTED: 'RUSH_STARTED',
  STOPPED: 'RUSH_STOPPED',
  NEW_RESULTS: 'RUSH_NEW_RESULTS'
};

export function NewResults(results) {
  return {
    type: RUSH.NEW_RESULTS,
    data: results
  };
}
