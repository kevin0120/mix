import axios from 'axios';
import axiosRetry from 'axios-retry';

const defaultClient = axios.create({
  timeout: 3000,
  headers: { 'Content-Type': 'application/json' }
});

axiosRetry(defaultClient, {
  retries: 2,
  retryDelay: axiosRetry.exponentialDelay,
  retryCondition: err => {
    if (
      err.message.indexOf('200') !== -1 ||
      err.message.indexOf('409') !== -1
    ) {
      return false;
    }

    return true;
  }
});

// 获取作业
export function fetchRoutingWorkcenter(url, workCenterCode, carType, job) {
  const fullUrl = `${url}/rush/v1/operation/${workCenterCode}`;

  let paramObj = {
    carType
  };

  if (job) {
    paramObj = {
      job
    };
  }

  return defaultClient
    .get(fullUrl, { params: paramObj })
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}

// 获取工单
export function fetchWorkorder(url, hmiSN, code) {
  const fullUrl = `${url}/rush/v1/workorder`;

  return defaultClient
    .get(fullUrl, {
      params: {
        hmi_sn: hmiSN,
        code
      }
    })
    .then(resp => resp)
    .catch(e => {
      throw e;
    });
}

// 拧紧枪使能控制
export function toolEnable(url, controllerSN, enable) {
  const fullUrl = `${url}/rush/v1/tool-enable`;
  return defaultClient
    .put(fullUrl, {
      controller_sn: controllerSN,
      enable
    })
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}

// 拧紧模式设置
export function controllerMode(url, mode, controllerSN) {
  const fullUrl = `${url}/rush/v1/controller-mode`;
  const bodyData = {
    mode,
    controller_sn: controllerSN
  };
  return defaultClient
    .put(fullUrl, bodyData)
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}

// pset
export function pset(url, controllerSN, gunSN, resultID, count, userID, pset) {
  const fullUrl = `${url}/rush/v1/psets`;
  return defaultClient
    .put(fullUrl, {
      pset,
      controller_sn: controllerSN,
      gun_sn: gunSN,
      result_id: resultID,
      count,
      user_id: userID
    })
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}

// 手动job
export function jobManual(
  url,
  controllerSN,
  carType,
  vin,
  userID,
  jobID,
  points,
  hmiSn,
  productID,
  workcenterId,
  skip,
  manualJob,
  mode
) {
  const fullUrl = `${url}/rush/v1/jobs-manual`;
  return defaultClient
    .put(fullUrl, {
      job_id: jobID,
      controller_sn: controllerSN,
      user_id: userID,
      car_type: carType,
      vin,
      mode,
      points,
      hmi_sn: hmiSn,
      workcenter_id: workcenterId,
      product_id: productID,
      skip,
      has_set: manualJob
    })
    .then(resp => resp)
    .catch(e => {
      throw e.toString();
    });
}
