// @flow

import type { Saga } from "redux-saga";
import type { IDevice } from '../../IDevice';

export interface IScrewTool extends IDevice{
  Enable(): void | Saga<void>
}