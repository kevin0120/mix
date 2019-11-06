// @flow
import type { IWorkStep } from '../../step/interface/IWorkStep';

export interface IOrder extends IWorkStep{
  +plannedDateTime: ?number,
  +productTypeImage: ?string,
  +datePlannedStart: Date,
  +workcenter: ?string,
  +productCode: ?string,
  +workingIndex: number
}
