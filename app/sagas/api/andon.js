import axios from "axios";

const defaultClient = axios.create({
  timeout: 3000,
  headers: { 'Content-Type': 'application/json' }
});

export function andonVehicleApi(baseUrl, vin, workcenterCode) {
  console.log(vin,workcenterCode);
  const fullUrl = `${baseUrl}/aiis/v1/vehicle`;
  return defaultClient
    .get(fullUrl, {
      params: {
        vin,
        workcenter_code:workcenterCode
      }
    })
    .then(resp => resp)
    .catch(e => {
      throw e;
    });
}
