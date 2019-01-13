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

