import { RESULT_DIAG } from './action';

const defaultResultDiag = {
  show: false
};

type actionType = {
  +type: string,
  +show: boolean
};

export default function resultDiag(
  state: object = defaultResultDiag,
  action: actionType
) {
  switch (action.type) {
    case RESULT_DIAG.SHOW: {
      const { show } = action;
      return { show };
    }
    default:
      return state;
  }
}
