// @flow
import type { Saga } from "redux-saga";

export interface ICommonExternalEntity {
  Healthz: boolean,
  +isEnable: boolean,
  +Name: string,
  +source: string,
  Enable(): void | Saga<void>,
  Disable(): void | Saga<void>
}