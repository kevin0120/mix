/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { defaultClient } from '../common/utils';
import { lError } from '../logger';

// eslint-disable-next-line import/prefer-default-export
export const apiToolEnable = (fullUrl, data) =>
  defaultClient
    .put(fullUrl, data)
    .then(resp => resp)
    .catch(e => {
      lError(e.toString(), {
        at: 'apiToolEnable',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e.toString();
    });
