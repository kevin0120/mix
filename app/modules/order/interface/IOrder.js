// @flow
import type { tOrder, tOrderStatus } from './typeDef';
import type { IWorkable } from '../../workable/IWorkable';
import type { IWorkStep } from '../../step/interface/IWorkStep';

export interface IOrder extends IWorkable {
  _datePlannedStart: Date,
  +datePlannedStart: Date,
  _datePlannedComplete: Date,
  +datePlannedComplete: Date,
  _productTypeImage: string,
  +productTypeImage: ?string,
  +workcenter: ?string,
  +productCode: ?string,
  +workingIndex: number,
  _workcenter: string,
  _workingIndex: number,
  _stateToRun: tOrderStatus,
  _status: tOrderStatus,
  _trackCode: string,
  _productCode: string,
  +workingStep: IWorkStep,
  constructor(tOrder, ...Array<any>): void
}
