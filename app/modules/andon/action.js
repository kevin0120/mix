export const ANDON = {
  NEW_DATA: 'ANDON_NEW_DATA',
  SCANNER: 'ANDON_SCANNER'
};

export function andonNewData(data) {
  return {
    type: ANDON.NEW_DATA,
    data
  };
}

export function andonScanner(vin) {
  return {
    type: ANDON.SCANNER,
    vin
  };
}
