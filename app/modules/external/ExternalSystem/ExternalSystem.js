// @flow
import { isURL } from 'validator/lib/isURL';
import { call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../../common/utils';
import CommonExternalEntity from '../CommonExternalEntity';

export default class ExternalSystem extends CommonExternalEntity {
  _endpoint: ?string = null;

  _name: string = '';

  constructor(name: string, endpoint: string) {
    super(name);
    this.Endpoint = endpoint;
  }

  get Endpoint(): ?string {
    return this._endpoint;
  }

  set Endpoint(endpoint: string) {
    if (!isURL(endpoint)) {
      CommonLog.lError(`Endpoint: ${endpoint} Is Not Valid!`);
    } else {
      this._endpoint = endpoint;
    }
  }

  *Connect(): Saga<boolean> {
    try {
      CommonLog.Warn('Please Override Connect Method!!!');
      yield call(this.setHealthz, true);
      return true;
    } catch (e) {
      CommonLog.lError(e, { at: 'Connect', name: this._name });
      return false;
    }
  }

  *Close(): Saga<boolean> {
    try {
      CommonLog.Warn('Please Override Close Method!!!');
      yield call(this.setHealthz, false);
      return true;
    } catch (e) {
      CommonLog.lError(e, { at: 'Connect', name: this._name });
      return false;
    }
  }
}
