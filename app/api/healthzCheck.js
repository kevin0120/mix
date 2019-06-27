/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import isURL from 'validator/lib/isURL';
import { Error } from '../logger';

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
      Error(e.toString(),{
        at:masterPCHealthCheck.name,
        response:e.response && e.response.data && JSON.stringify(e.response.data)
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
      Error(e.toString(),{
        at:controllerHealthCheck.name,
        response:e.response && e.response.data && JSON.stringify(e.response.data)
      });
      throw e;
    });
}
