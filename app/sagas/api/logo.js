import { defaultClient } from '../../common/utils';

import isURL from 'validator/lib/isURL';

export function fetchLogo(baseURL) {
  if (!isURL(baseURL, { require_protocol: true })) {
    throw new Error('fetchLogo baseURL is empty');
  }
  const fullUrl = `${baseURL}/logo`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}
