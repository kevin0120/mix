// @flow

import { toast } from 'react-toastify';
import axios from 'axios';
import axiosRetry, { exponentialDelay } from 'axios-retry';
import { isNil, cloneDeep } from 'lodash-es';
import { Info, lError, Warn, Debug, Maintenance } from '../logger';
import moment, { DurationInputArg1 } from 'moment';

const VINMap = {
  A: 1,
  B: 2,
  C: 3,
  D: 4,
  E: 5,
  F: 6,
  G: 7,
  H: 8,
  J: 1,
  K: 2,
  L: 3,
  M: 4,
  N: 5,
  P: 7,
  R: 9,
  S: 2,
  T: 3,
  U: 4,
  V: 5,
  W: 6,
  X: 7,
  Y: 8,
  Z: 9,
  1: 1,
  2: 2,
  3: 3,
  4: 4,
  5: 5,
  6: 6,
  7: 7,
  8: 8,
  9: 9,
  0: 0
};

// eslint-disable-next-line flowtype/no-weak-types
export function sortObj(obj: any, orderKey: string): any {
  const orderedKey = Object.keys(obj).sort(
    (a, b) => obj[a][orderKey] - obj[b][orderKey]
  );
  return orderedKey.map(key => ({
    key,
    value: cloneDeep(obj[key])
  }));
}

export function normalSortObj(obj) {
  const orderedKey = Object.keys(obj);
  return orderedKey.map(key => ({
    key,
    value: cloneDeep(obj[key])
  }));
}

// type: default、error、warning、success、info
export function showNoty(type, text) {
  toast(text, {
    type,
    position: 'top-right',
    autoClose: 5000,
    hideProgressBar: false,
    closeOnClick: true,
    pauseOnHover: true,
    draggable: true,
    draggablePercent: 60
  });
}

export function isCarID(data: string): boolean {
  return !!(isVin(data) || isKnr(data) || isLongPin(data));
}

export function isVin(data: string): boolean {
  if (data.length !== 17) return false;
  const cd: 'X' | number = getCheckDigit(data);
  if (cd === 'X') {
    return cd === data[8];
  }
  return cd === VINMap[data[8]];
}

function isKnr(data: string): boolean {
  // if (data.length === 8 && checkNaturalNumber(data)) {
  //   return true;
  // }

  return data.length === 8 && checkNaturalNumber(data);
}

function isLongPin(data: string): boolean {
  if (data.length === 14) {
    const year = data.substring(2, 6);
    const knr = data.substring(6);

    if (checkNaturalNumber(year) && isKnr(knr)) {
      return true;
    }
  }

  return false;
}

function checkNaturalNumber(data: string): boolean {
  return /^[0-9]+$/.test(data) && data > 0;
}

// source: https://en.wikipedia.org/wiki/Vehicle_identification_number#Example_Code
// ---------------------------------------------------------------------------------
// function transliterate(c) {
//   return '0123456789.ABCDEFGH..JKLMN.P.R..STUVWXYZ'.indexOf(c) % 10;
// }

function getCheckDigit(vin: string): 'X' | number {
  const weights = [8, 7, 6, 5, 4, 3, 2, 10, 0, 9, 8, 7, 6, 5, 4, 3, 2];
  let sum = 0;
  for (let i = 0; i < 17; i += 1) sum += VINMap[vin[i]] * weights[i];
  let checkDigit = sum % 11;
  if (checkDigit === 10) checkDigit = 'X';
  return checkDigit;
}

export class HttpClient {
  // eslint-disable-next-line flowtype/no-weak-types
  instance: any = null;

  constructor() {
    this.instance = null;
  }

  getInstance() {
    if (isNil(this.instance)) {
      this.instance = axios.create({
        timeout: 3000,
        headers: { 'Content-Type': 'application/json' }
      });
      axiosRetry(this.instance, {
        retries: 2,
        retryDelay: exponentialDelay,
        retryCondition: err => err.message.indexOf('200') === -1
      });
      return this.instance;
    }
    return this.instance;
  }
}

export const defaultClient = new HttpClient().getInstance();

export type CommonLogLvl = 'Warn' | 'Info' | 'Error' | 'Debug' | 'Maintenance';

// eslint-disable-next-line flowtype/no-weak-types
type tCommonLogMeta = Object;

/* eslint-disable prefer-rest-params */
/* eslint-disable no-unused-vars */
export const CommonLog = {
  Info(msg: string, meta: ?tCommonLogMeta) {
    _logger('Info', ...arguments);
  },

  Warn(msg: string, meta: ?tCommonLogMeta) {
    _logger('Warn', ...arguments);
  },

  Debug(msg: string, meta: ?tCommonLogMeta) {
    _logger('Debug', ...arguments);
  },

  Maintenance(msg: string, meta: ?tCommonLogMeta) {
    _logger('Maintenance', ...arguments);
  },

  lError(msg: mixed, meta: ?tCommonLogMeta) {
    _logger('Error', ...arguments);
  }
};
/* eslint-enable prefer-rest-params */
/* eslint-enable no-unused-vars */


const fn = {
  Error: console.error,
  Info: console.info,
  Warn: console.warn,
  Debug: console.log,
  Maintenance: console.log
};

const fnAlways = {
  Error: (...rest) => {
    if (rest) {
      if (typeof rest[0] === 'string') {
        lError(...rest);
      }
      if (rest[0] instanceof Error) {
        lError(rest[0].message, rest[1]);
      }
    }
  },
  Info,
  Warn,
  Debug,
  Maintenance
};

// eslint-disable-next-line no-underscore-dangle
function _logger(lvl: CommonLogLvl, ...rest) {
  if (process.env.NODE_ENV !== 'production' || process.env.DEBUG_PROD === 'true') {
    // eslint-disable-next-line no-unused-expressions
    fn[lvl] && fn[lvl](...rest);
  }
  // eslint-disable-next-line no-unused-expressions
  fnAlways[lvl] && fnAlways[lvl](...rest);
}

export const timeCost = (times: Array<Date>) =>
  ((times || []).length % 2 === 0 ? times || [] : [...times, new Date()])
    .reduce((total, currentTime, idx) =>
      idx % 2 === 0 ? total - currentTime : total - (0 - currentTime), 0);


export const formatTime = (t: number) => `${t}`.length <= 2 ? `00${t}`.slice(-2) : t;

export function durationString(duration: DurationInputArg1): string {
  const h = moment.duration(duration).hours();
  const m = moment.duration(duration).minutes();
  const s = moment.duration(duration).seconds();
  return `${formatTime(h)}:${formatTime(m)}:${formatTime(s)}`;
}
