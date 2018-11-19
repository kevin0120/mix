import { defaultClient } from '../../common/utils';

// 获取用户信息
export function getUserInfo(url, uuid) {
  const fullUrl = `${url}/res.users/${uuid}`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}
