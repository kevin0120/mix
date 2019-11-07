// @flow
import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import Device from '../Device';
import ClsIOModule from '../io/ClsIOModule';
import type { tIOContact, tIOConfig } from '../io/type';
import { CommonLog } from '../../../common/utils';

class ClsController extends Device {

  _io = null;

  // eslint-disable-next-line no-unused-vars,flowtype/no-weak-types
  constructor(name: string, serialNumber: string, config: tIOConfig, ...rest: Array<any>) {
    super(name, serialNumber);
    this._serialNumber = serialNumber;
    this._io = new ClsIOModule(name, serialNumber, config);
  }

  * doDispatch(data: tIOContact): Saga<void> {
    try {
      if (!this._io) {
        return;
      }
      yield call(this._io.doDispatch, data);
    } catch (e) {
      CommonLog.lError(e, { at: 'ClsController.doDispatch' });
    }
  }
}

export default ClsController;
