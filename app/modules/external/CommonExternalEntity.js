// @flow
import { isEqual } from 'lodash-es';
import type { Saga } from 'redux-saga';
import { all, put, call } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import type { ICommonExternalEntity } from './ICommonExternalEntity';
import { makeListener } from '../util';
import type { tAction, tListener } from '../typeDef';

export default class CommonExternalEntity implements ICommonExternalEntity {
  _name: string;

  _isHealthz: boolean = false;

  _enable: boolean = false;

  _children: Set<ICommonExternalEntity> = new Set();

  _healthzListener = makeListener();

  constructor(name: string) {
    this._name = name;

    /* eslint-disable flowtype/no-weak-types */
    (this: any).Enable = this.Enable.bind(this);
    (this: any).Disable = this.Disable.bind(this);
    (this: any).ToggleEnable = this.ToggleEnable.bind(this);
    /* eslint-enable flowtype/no-weak-types */
  }

  appendChildren(
    children: Array<ICommonExternalEntity> | ICommonExternalEntity
  ) {
    if (children instanceof Array) {
      children.forEach(c => {
        this._children.add(c);
      });
    } else {
      this._children.add(children);
    }
  }

  deleteChildren(children: Array<ICommonExternalEntity>) {
    if (children instanceof Array) {
      children.forEach(c => {
        this._children.delete(c);
      });
    }
  }

  getChildren(
    patten: ICommonExternalEntity => boolean
  ): Array<ICommonExternalEntity> {
    if (!patten) {
      return [...this._children];
    }
    return [...this._children].filter(patten);
  }

  // eslint-disable-next-line flowtype/no-weak-types
  bindOnHealthzAction(
    predicate: boolean => boolean,
    action: boolean => tAction<any, any>
  ) {
    return this._healthzListener.add(predicate, action);
  }

  removeOnHealthzAction(listener: tListener<boolean>) {
    return this._healthzListener.remove(listener);
  }

  *setHealthz(isHealthz: boolean): Saga<void> {
    try {
      if (isEqual(this._isHealthz, isHealthz)) {
        return;
      }
      this._isHealthz = isHealthz;
      const actions = this._healthzListener.check(isHealthz);

      if (!isHealthz) {
        this.Disable();
      }

      if (this._children.size > 0 && !isHealthz) {
        yield all([...this._children].map(c => call(c.setHealthz, false)));
      }

      const msg = `${
        this._name
      } Healthz Status Change: ${isHealthz.toString()}`;
      CommonLog.Info(msg);

      yield all(actions.map(a => put(a)));
    } catch (e) {
      CommonLog.lError(e, { at: 'setHealthz', name: this._name });
    }
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

  Enable(): void | Saga<void> {
    this._enable = true;
    CommonLog.Info(`${this.source} Is Enabled!`);
  }

  Disable(): void | Saga<void> {
    this._enable = false;
    CommonLog.Info(`${this.source} Is Disabled!`);
  }

  // eslint-disable-next-line require-yield
  *ToggleEnable(): Saga<void> {
    this._enable = !this._enable;
  }
}
