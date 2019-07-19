import axios from 'axios';
import { Error } from '../logger';

const defaultClient = axios.create({
  timeout: 3000,
  headers: { 'Content-Type': 'application/json' }
});

// eslint-disable-next-line import/prefer-default-export
export function andonVehicleApi(baseUrl, vin, workcenterCode) {
  console.log(vin, workcenterCode);
  const fullUrl = `${baseUrl}/aiis/v1/vehicle`;
  return defaultClient
    .get(fullUrl, {
      params: {
        vin,
        workcenter_code: workcenterCode
      }
    })
    .then(resp => resp)
    .catch(e => {
      Error(e.toString(), {
        at: 'andonVehicleApi',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
