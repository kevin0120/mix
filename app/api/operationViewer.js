import { defaultClient } from '../common/utils';
import { Error } from '../logger';

export function operationListApi(baseURL, hmiSn) {
  const fullUrl = `${baseURL}/mrp.routing.workcenter`;

  return defaultClient
    .get(fullUrl, {
      params: {
        hmi_sn: hmiSn
      }
    })
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(),{
        at:operationListApi.name,
        response:e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e.toString();
    });
}

export function operationDetailApi(baseURL, operationID) {
  const fullUrl = `${baseURL}/mrp.routing.workcenter/${operationID}`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(),{
        at:operationDetailApi.name,
        response:e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e.toString();
    });
}

// 图片编辑
export function imageEditApi(baseURL, operationID, points, img) {
  const fullUrl = `${baseURL}/mrp.routing.workcenter/${operationID}/edit`;

  return defaultClient
    .put(fullUrl,{
      points,
      img
    })
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(),{
        at:imageEditApi.name,
        response:e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e.toString();
    });
}
