import { LOGO } from '../actions/actionTypes';
import defaultWorkingImg from '../../resources/imgs/defaultWorking.jpg';

const initOperationViewer = defaultWorkingImg;
export default function logo(
  state: object = initOperationViewer,
  action: actionType
) {
  switch (action.type) {
    case LOGO.FETCH_OK: {
      return action.logo || defaultWorkingImg;
    }
    default:
      return defaultWorkingImg;
  }
}
