// @flow

import { OPERATION_SOURCE } from '../operation/model';

export const SCANNER = {
  DEVICE_NOT_FOUND: 'SCANNER_DEVICE_NOT_FOUND',
  START_ERROR: 'SCANNER_START_ERROR',
  READ_ERROR: 'SCANNER_READ_ERROR',
  READ_NEW_DATA: 'SCANNER_NEW_DATA'
};

let scanner = null;

const usbDetect = require('usb-detection');
const HID = require('node-hid');

const lodash = require('lodash');

function getBarcode(mode, chunk) {
  let buf;
  switch (mode) {
    case 'HID': {
      const idx = lodash.findIndex(chunk, o => o === 0);
      buf = Buffer.from(chunk.slice(5, idx), 'hex');
      return buf.toString();
    }

    case 'BT_HID': {
      const idx1 = lodash.indexOf(chunk, 0);
      const idx2 = lodash.indexOf(chunk, 0, idx1 + 1);
      buf = Buffer.from(chunk.slice(idx1 + 1, idx2), 'hex');
      return buf.toString();
    }
    default:
      console.log('error mode');
  }
}

export function NewCar(
  aBarcode: string,
  source: string = OPERATION_SOURCE.SCANNER
) {
  return {
    type: SCANNER.READ_NEW_DATA,
    data: aBarcode,
    source
  };
}

export const listenToNewCar = (dispatch, getState) => {
  const config = getState().setting;
  const { vendorId, mode } = config.system.workcenter.hardware.scanner;
  if (scanner == null) {
    try {
      const devices = HID.devices();
      const deviceInfo = devices.find(d => d.vendorId === vendorId);
      if (deviceInfo) {
        scanner = new HID.HID(deviceInfo.path);
        ScannerProcess(scanner, dispatch, mode);
      }
    } catch (error) {
      scanner = null;
      dispatch({
        type: SCANNER.READ_ERROR,
        data: '请插拔扫码枪'.concat(error.toString())
      });
    }
  }
  usbDetect.startMonitoring();

  usbDetect.on('add:'.concat(vendorId), device => {
    try {
      scanner = new HID.HID(device.vendorId, device.productId);
      ScannerProcess(scanner, dispatch, mode);
    } catch (e) {
      if (scanner != null) {
        scanner.close();
        scanner = null;
      }
    }
  });
  usbDetect.on('remove:'.concat(vendorId), () => {
    if (scanner != null) {
      scanner.close();
      scanner = null;
    }
  });
};

const ScannerProcess = (s, dispatch, mode) => {
  s.on('data', chunk => {
    const aBarcode = getBarcode(mode, chunk);
    dispatch(NewCar(aBarcode));
  });
  scanner.on('error', error => {
    scanner.close();
    scanner = null;
    dispatch({
      type: SCANNER.READ_ERROR,
      data: '请插拔扫码枪'.concat(error.toString())
    });
  });
};

export const stopListernScanner = () => {
  usbDetect.stopMonitoring();
  if (scanner != null) {
    scanner.close();
    scanner = null;
  }
};
