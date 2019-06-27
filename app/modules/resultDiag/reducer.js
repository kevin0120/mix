/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

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
