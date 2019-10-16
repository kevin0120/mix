// @flow
import { isURL } from 'validator/lib/isURL';
import { CommonLog } from '../../../common/utils';
import CommonExternalEntity from '../CommonExternalEntity';

export default class ExternalSystem extends CommonExternalEntity {
  _endpoint: ?string = null;

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

  Connect(): boolean {
    CommonLog.Warn('Please Override Connect Method!!!');
    this.Healthz = false;
    return false;
  }

  Close(): boolean {
    CommonLog.Warn('Please Override Close Method!!!');
    this.Healthz = false;
    return true;
  }
}
