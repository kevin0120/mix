// @flow

type tCommonActionType = {
  +type: string
};

type tDeviceNewData = {
  +data: string,
  +source: string
};

interface IInputDevice {
  static validate(data: string | number): boolean
}

interface IOutputDevice {
  validate(data: any): boolean
}

class Device {
  source: string;

  constructor(source: string) {
    this.source = source;
  }
}

export type {tCommonActionType,tDeviceNewData, IInputDevice, IOutputDevice, Device};
