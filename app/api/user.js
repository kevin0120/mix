// @flow
import { CommonLog, defaultClient } from '../common/utils';
import type { tUuid } from '../modules/user/interface/typeDef';

// 获取用户信息
// eslint-disable-next-line import/prefer-default-export
export function getUserInfo(url: string, uuid: tUuid) {
  const fullUrl = `${url}/res.users/${uuid}`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      CommonLog.lError(e.toString(), {
        at: 'getUserInfo',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
