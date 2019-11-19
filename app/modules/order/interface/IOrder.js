// @flow
import type { Saga } from 'redux-saga';
import type { tOrder, tOrderStatus } from './typeDef';
import type { IWorkable } from '../../workable/IWorkable';
import type { IWorkStep } from '../../step/interface/IWorkStep';

export interface IOrder extends IWorkable {
  _datePlannedStart: ?Date;
  +datePlannedStart: ?Date;

  _datePlannedComplete: ?Date;
  +datePlannedComplete: ?Date;

  _productTypeImage: string;
  +productTypeImage: ?string;

  _workcenter: string;
  +workcenter: ?string;

  _workingIndex: number;
  +workingIndex: number;

  _status: tOrderStatus;
  +status: tOrderStatus;

  _trackCode: string;
  +trackCode: string;

  _productCode: string;
  +productCode: ?string;

  +workingStep: IWorkStep;
  constructor(tOrder, ...Array<any>): void;
  update(dataObj: ?$Shape<tOrder>): void;
  updateStatus({ status: tOrderStatus }): Saga<void>;
}
