import { cloneDeep } from 'lodash';
import { toast } from 'react-toastify';

export function sortObj(obj, orderKey) {
  const orderedKey = Object.keys(obj).sort((a, b) => (
    obj[a][orderKey] - obj[b][orderKey]));
  return orderedKey.map(key => ({
    key,
    value: cloneDeep(obj[key]),
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
    draggablePercent: 60,
  });
}

export function isVin(vin) {
  if (vin.length !== 17) return false;
  return getCheckDigit(vin) === vin[8];
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

