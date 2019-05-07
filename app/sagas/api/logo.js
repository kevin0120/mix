import isURL from 'validator/lib/isURL';
import { defaultClient } from '../../common/utils';
import { Error } from '../../logger';


export default function fetchLogo(baseURL) {
  if (!isURL(baseURL, { require_protocol: true })) {
    throw new Error('fetchLogo baseURL is empty');
  }
  const fullUrl = `${baseURL}/logo`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(),{
        at:fetchLogo.name
      });
      console.log(e);
      throw e;
    });
}
