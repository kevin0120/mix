import { LOGO } from './action';
import defaultWorkingImg from '../../../resources/imgs/defaultWorking.jpg';

export default function logo(
  state: object = defaultWorkingImg,
  action: actionType
) {
  switch (action.type) {
    case LOGO.FETCH_OK: {
      return action.logo || defaultWorkingImg;
    }
    default:
      return state;
  }
}
