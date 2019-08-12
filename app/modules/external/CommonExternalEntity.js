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

  constructor(name: string) {
    this.#name = name;

    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  set Healthz(isHealthz: boolean) {
    if (isEqual(this.#isHealthz, isHealthz)) {
      return;
    }
    this.#isHealthz = isHealthz;
    if (!isHealthz) {
      this.Disable();
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

  // eslint-disable-next-line require-yield
  * Enable(): Saga<void> {
    this.#enable = true;
    CommonLog.Debug(`${this.source} Is Enabled!`);
  }

  // eslint-disable-next-line require-yield
  * Disable(): Saga<void> {
    this.#enable = false;
    CommonLog.Debug(`${this.source} Is Disabled!`);
  }

  // eslint-disable-next-line require-yield
  * ToggleEnable(): Saga<void> {
    this.#enable = !this.#enable;
  }

}