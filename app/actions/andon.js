import {ANDON} from './actionTypes';

export function andonNewData(data) {
  return{
    type:ANDON.NEW_DATA,
    data
  };
}

export function andonScanner(vin) {
  return{
    type:ANDON.SCANNER,
    vin
  };
}
