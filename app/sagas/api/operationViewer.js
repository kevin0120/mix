import { defaultClient } from '../../common/utils';

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
      throw e.toString();
    });
}

export function operationDetailApi(baseURL, operationID) {
  const fullUrl = `${baseURL}/mrp.routing.workcenter/${operationID}`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
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
      throw e.toString();
    });
}
