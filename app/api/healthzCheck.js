// @flow

import isURL from 'validator/lib/isURL';
import { lError } from '../logger';

const axios = require('axios');

const defaultClient = axios.create({
  timeout: 1500,
  headers: { 'Content-Type': 'application/json' }
});

export function masterPCHealthCheck(conn) {
  if (!isURL(conn, { require_protocol: true })) {
    throw new Error('conn is empty');
  }
  const fullUrl = `${conn}/rush/v1/healthz`;

  return defaultClient
    .get(fullUrl)
    .then(resp => resp)
    .catch(e => {
      lError(e.toString(), {
        at: 'masterPCHealthCheck',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}

export function controllerHealthCheck(conn, controllers) {
  if (!isURL(conn, { require_protocol: true }) || controllers.length === 0) {
    throw new Error('conn or controller is empty');
  }
  const url = `${conn}/rush/v1/controller-status`;
  const queryparam = controllers.join(',');
  return defaultClient
    .get(url, {
      params: {
        controller_sn: queryparam
      }
    })
    .then(response => response)
    .catch(e => {
      lError(e.toString(), {
        at: 'controllerHealthCheck',
        response:
          e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
