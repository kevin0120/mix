import { cloneDeep } from 'lodash';
import { toast } from 'react-toastify';

import axios from 'axios';
import axiosRetry from 'axios-retry';


const VINMAP = {A : 1, B : 2, C : 3, D : 4, E : 5, F : 6, G : 7, H : 8,
  J : 1, K : 2, L : 3, M : 4, N : 5,        P : 7,        R : 9,
  S : 2, T : 3, U : 4, V : 5, W : 6, X : 7, Y : 8, Z : 9,
  1 : 1, 2 : 2, 3 : 3, 4 : 4, 5 : 5, 6 : 6, 7 : 7, 8 : 8, 9 : 9, 0 : 0
};

export function sortObj(obj, orderKey) {
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

export function isCarID(data) {
  if (isVin(data) || isKnr(data) || isLongpin(data)) {
    return true;
  }

  return false;
}

export function isVin(data) {
  if (data.length !== 17) return false;
  const cd = getCheckDigit(data);
  if (cd === 'x') {
    return cd === data[8];
  }
  return cd === VINMAP[data[8]];
}

function isKnr(data) {
  // if (data.length === 8 && checkNaturalnumber(data)) {
  //   return true;
  // }

  return data.length === 8 && checkNaturalnumber(data);
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

function checkNaturalnumber(data) {
  if (/^[0-9]+$/.test(data) && data > 0) {
    return true;
  }

  return false;
}

// source: https://en.wikipedia.org/wiki/Vehicle_identification_number#Example_Code
// ---------------------------------------------------------------------------------
// function transliterate(c) {
//   return '0123456789.ABCDEFGH..JKLMN.P.R..STUVWXYZ'.indexOf(c) % 10;
// }

function getCheckDigit(vin) {
  const weights = [8, 7, 6, 5, 4, 3, 2, 10, 0, 9, 8, 7, 6, 5, 4, 3, 2];
  let sum = 0;
  for (let i = 0; i < 17; i += 1)
    sum += VINMAP[vin[i]] * weights[i];
  let checkDigit = sum % 11;
  if(checkDigit === 10) checkDigit = 'X';
  return checkDigit;
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
