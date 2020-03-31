// @flow

import isURL from 'validator/lib/isURL';
import { CommonLog, defaultClient } from '../../common/utils';

export default function fetchLogo(baseURL: string) {
  if (!isURL(baseURL, { require_protocol: true })) {
    throw new Error('fetchLogo baseURL Is Empty');
  }
  const fullUrl = `${baseURL}/logo`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      CommonLog.lError(e.toString(), {
        at: 'fetchLogo',
        response: e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
