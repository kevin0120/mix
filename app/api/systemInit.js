// @flow

import { defaultClient } from '../common/utils';
import { lError } from '../logger';

// eslint-disable-next-line import/prefer-default-export
export const fetchConnectionInfo = (fullUrl: string) =>
  defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      lError(e.toString(), {
        at: 'fetchConnectionInfo',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e.toString();
    });
