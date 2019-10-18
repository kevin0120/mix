// @flow
import type { IWorkStep } from '../../step/interface/IWorkStep';

export interface IOrder extends IWorkStep{
  +plannedDateTime: ?number,
  +workingIndex: number
}
