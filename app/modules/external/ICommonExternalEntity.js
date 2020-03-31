// @flow
import type { Saga } from 'redux-saga';

export interface ICommonExternalEntity {
  _name: string,
  +Healthz: boolean,
  +isEnable: boolean,
  +Name: string,
  +source: string,

  Enable(): void | Saga<void>,

  setHealthz(isHealthz: boolean): Saga<void>,

  Disable(): void | Saga<void>,

  appendChildren(children: Array<ICommonExternalEntity> | ICommonExternalEntity): void,

  setParent(parent: ICommonExternalEntity): void,

  +parent: ICommonExternalEntity
}
