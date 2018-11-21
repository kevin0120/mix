import { cloneDeep } from 'lodash';
import { toast } from 'react-toastify';

import axios from 'axios';
import axiosRetry from 'axios-retry';

export function sortObj(obj, orderKey) {
  const orderedKey = Object.keys(obj).sort(
    (a, b) => obj[a][orderKey] - obj[b][orderKey]
  );
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

export function isCarID(data) {
  if (isVin(data) || isKnr(data) || isLongpin(data)) {
    return true;
  }

  return false;
}

export function isVin(vin) {
  if (vin.length !== 17) return false;
  return getCheckDigit(vin) === vin[8];
}

function isKnr(data) {
  if (data.length === 8 && checkNaturalnumber(data)){
    return true;
  }

  return false
}

function isLongpin(data) {
  if (data.length === 14) {
    const year = data.substring(2, 6);
    const knr = data.substring(6);

    if (checkNaturalnumber(year) && isKnr(knr)) {
      return true;
    }
  }

  return false;
}

function checkNaturalnumber(data){
  if(/^[0-9]+$/.test(data) && (data>0)){
    return true;
  }

  return false;
}

// source: https://en.wikipedia.org/wiki/Vehicle_identification_number#Example_Code
// ---------------------------------------------------------------------------------
function transliterate(c) {
  return '0123456789.ABCDEFGH..JKLMN.P.R..STUVWXYZ'.indexOf(c) % 10;
}

function getCheckDigit(vin) {
  const map = '0123456789X';
  const weights = '8765432X098765432';
  let sum = 0;
  for (let i = 0; i < 17; i += 1)
    sum += transliterate(vin[i]) * map.indexOf(weights[i]);
  return map[sum % 11];
}

export class HttpClient {
  constructor() {
    this.instance = null;
  }

  static getInstance() {
    if (!this.instance) {
      this.instance = axios.create({
        timeout: 3000,
        headers: { 'Content-Type': 'application/json' }
      });
      axiosRetry(this.instance, {
        retries: 2,
        retryDelay: axiosRetry.exponentialDelay,
        retryCondition: err => err.message.indexOf('200') === -1
      });
      return this.instance;
    }
    return this.instance;
  }
}

export const defaultClient = HttpClient.getInstance();
