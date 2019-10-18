// @flow
import LOGO from './action';
import defaultWorkingImg from '../../../resources/imgs/defaultWorking.jpg';
import type { tActionLogo } from './type';

export default function logo(
  state: Object = defaultWorkingImg,
  action: tActionLogo
) {
  switch (action.type) {
    case LOGO.FETCH_OK: {
      return action.logo || defaultWorkingImg;
    }
    default:
      return state;
  }
}
