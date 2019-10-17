// @flow

import type { IWorkStep } from '../../interface/IWorkStep';

export type tMaterialStepPayload = {
  items: Array<any>
};

export interface IMaterialStep extends IWorkStep{
  _payload: ?tMaterialStepPayload
}
