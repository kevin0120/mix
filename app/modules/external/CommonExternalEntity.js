// @flow
import { isEqual } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { CommonLog } from '../../common/utils';

interface IHealthChecker {
  Healthz: boolean
}

export default class CommonExternalEntity implements IHealthChecker {
  #name: string;

  #isHealthz: boolean = false;

  #enable: boolean = false;

  #_children: Set<CommonExternalEntity> = new Set([]);

  constructor(name: string) {
    this.#name = name;

    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  appendChildren(children: Array<CommonExternalEntity> | CommonExternalEntity) {
    if (children instanceof Array) {
      children.forEach((c) => {
        this.#_children.add(c);
      });
    } else {
      this.#_children.add(children);
    }
  }

  deleteChildren(children: Array<CommonExternalEntity>) {
    if (children instanceof Array) {
      children.forEach(c => {
        this.#_children.delete(c);
      });
    }
  };

  getChildren(patten: (CommonExternalEntity)=>boolean) {
    if (!patten) {
      return [...this.#_children];
    }
    return [...this.#_children].filter(patten);
  }

  set Healthz(isHealthz: boolean) {
    if (isEqual(this.#isHealthz, isHealthz)) {
      return;
    }
    this.#isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
    }
    if (this.#_children.length > 0 && !isHealthz) {
      this.#_children.forEach(c => {
        c.Healthz = false;
      });
    }

    if (this.#_children.length > 0 && isHealthz) {
      this.#_children.forEach(c => {

      });
    }

    const msg = `${this.#name} Healthz Status Change: ${isHealthz.toString()}`;
    CommonLog.Info(msg);
  }

  get Healthz(): boolean {
    return this.#isHealthz;
  }

  get Name(): string {
    return this.#name;
  }

  get source(): string {
    return this.#name;
  }

  get isEnable(): boolean {
    return this.#enable;
  }

  Enable(): Saga<void> {
    this.#enable = true;
    CommonLog.Info(`${this.source} Is Enabled!`);
  }

  Disable() {
    this.#enable = false;
    CommonLog.Info(`${this.source} Is Disabled!`);
  }

  // eslint-disable-next-line require-yield
  * ToggleEnable(): Saga<void> {
    this.#enable = !this.#enable;
  }

}
