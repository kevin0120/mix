
import axios from 'axios';
import axiosRetry from 'axios-retry';
import isURL from 'validator/lib/isURL';

const defaultClient = axios.create({
  timeout: 3000,
  headers: { 'Content-Type': 'application/json' }
});

axiosRetry(defaultClient, {
  retries: 2,
  retryDelay: axiosRetry.exponentialDelay,
  retryCondition: err => {
    return !(err.message.indexOf('200') !== -1 ||
      err.message.indexOf('409') !== -1 ||
      err.message.indexOf('404') !== -1);

  }
});

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

