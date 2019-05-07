import { defaultClient } from '../../common/utils';
import { Error } from '../../logger';

// 获取用户信息
export function getUserInfo(url, uuid) {
  const fullUrl = `${url}/res.users/${uuid}`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(),{
        at:getUserInfo.name
      });
      throw e.toString();
    });
}
