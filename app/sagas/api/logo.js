import { defaultClient } from '../../common/utils';


export function fetchLogo(baseURL) {
  const fullUrl = `${baseURL}/logo`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}
