// @flow
import { isEqual } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';

interface IHealthChecker {
  Healthz: boolean
}

export default class CommonExternalEntity implements IHealthChecker {
  _name: string;

  _isHealthz: boolean = false;

  _enable: boolean = false;

  _children: Set<CommonExternalEntity> = new Set();

  constructor(name: string) {
    this._name = name;

    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  appendChildren(children: Array<CommonExternalEntity> | CommonExternalEntity) {
    if (children instanceof Array) {
      children.forEach(c => {
        this._children.add(c);
      });
    } else {
      this._children.add(children);
    }
  }

  deleteChildren(children: Array<CommonExternalEntity>) {
    if (children instanceof Array) {
      children.forEach(c => {
        this._children.delete(c);
      });
    }
  }

  getChildren(patten: CommonExternalEntity => boolean) {
    if (!patten) {
      return [...this._children];
    }
    return [...this._children].filter(patten);
  }

  set Healthz(isHealthz: boolean) {
    if (isEqual(this._isHealthz, isHealthz)) {
      return;
    }
    this._isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
    }
    if (this._children.size > 0 && !isHealthz) {
      this._children.forEach(c => {
        c.Healthz = false;
      });
    }

    if (this._children.size > 0 && isHealthz) {
      this._children.forEach(c => {
      });
    }

    const msg = `${this._name} Healthz Status Change: ${isHealthz.toString()}`;
    CommonLog.Info(msg);
  }

  get Healthz(): boolean {
    return this._isHealthz;
  }

  get Name(): string {
    return this._name;
  }

  get source(): string {
    return this._name;
  }

  get isEnable(): boolean {
    return this._enable;
  }

  Enable(): any {
    this._enable = true;
    CommonLog.Info(`${this.source} Is Enabled!`);
  }

  Disable(): any {
    this._enable = false;
    CommonLog.Info(`${this.source} Is Disabled!`);
  }

  // eslint-disable-next-line require-yield
  * ToggleEnable(): Saga<void> {
    this._enable = !this._enable;
  }
}
