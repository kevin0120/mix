import { defaultClient } from '../../common/utils'

// 获取作业
export function fetchRoutingWorkcenter(url, workCenterCode, carType, job) {
  const fullUrl = `${url}/rush/v1/operation/${workCenterCode}`;

  let paramObj = {
    carType: carType
  };

  if (job) {
    paramObj = {
      job: job
    };
  }

  return defaultClient.get(fullUrl, {params: paramObj}).then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });
}

// 获取工单
export function fetchWorkorder(url, hmiSN, code) {
  const fullUrl = `${url}/rush/v1/wokorders`;

  return defaultClient.get(fullUrl, {params: {
      hmi_sn: hmiSN,
      code: code
    }}).then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });
}

// 拧紧枪使能控制
export function toolEnable(url, controllerSN, enable) {
  const fullUrl = `${url}/rush/v1/tool-enable`;
  return defaultClient.put(fullUrl, {
    controller_sn: controllerSN,
    enable: enable,
  }).then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });
}

// 拧紧模式设置
export function controllerMode(url, mode, controllerSN) {
  const fullUrl = `${url}/rush/v1/controller-mode`;
  const bodyData = {
    mode: mode,
    controller_sn: controllerSN,
  };
  return defaultClient.put(fullUrl, bodyData)
    .then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });
}

// pset
export function pset(url, controllerSN, gunSN, resultID, count, userID) {
  const fullUrl = `${url}/rush/v1/psets`;
  return defaultClient.put(fullUrl, {
    pset: pset,
    controller_sn: controllerSN,
    gun_sn: gunSN,
    result_id: resultID,
    count: count,
    user_id: userID,
  })
    .then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });

}

// 手动job
export function jobManual(url, controllerSN, carType, vin, userID, jobID, points, hmiSn, productID, workcenterId, skip, manualJob, mode) {

  const fullUrl = `${url}/rush/v1/jobs-manual`;
  return defaultClient.put(fullUrl, {
    job_id: jobID,
    controller_sn: controllerSN,
    user_id: userID,
    car_type: carType,
    vin: vin,
    mode: mode,
    points: points,
    hmi_sn: hmiSn,
    workcenter_id: workcenterId,
    product_id: productID,
    skip: skip,
    has_set: manualJob
  })
    .then(resp => (resp))
    .catch(e => {
      throw e.toString();
    });
}
