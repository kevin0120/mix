// @flow

import type { IWorkStep } from '../../interface/IWorkStep';
import type { IIOModule } from '../../../external/device/io/interface/IIOModule';
import type { tIOPort } from '../../../external/device/io/type';

export type tMaterialStepPayload = {
  items: Array<any>
};

export interface IMaterialStep extends IWorkStep {
  _payload: ?tMaterialStepPayload,
  _io: Set<IIOModule>,
  _ports: Set<tIOPort>,
  _confirm: ?{
    io: IIOModule,
    port: tIOPort
  },
  _items: Set<{
    in: {
      io: IIOModule,
      port: tIOPort
    },
    out: {
      io: IIOModule,
      port: tIOPort
    }
  }>
}
