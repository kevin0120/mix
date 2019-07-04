// @flow

export const RESULT_DIAG = {
  SHOW: 'RESULT_DIAG.SHOW'
};
export function setResultDiagShow(show: boolean) {
  return {
    type: RESULT_DIAG.SHOW,
    show
  };
}
