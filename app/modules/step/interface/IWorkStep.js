// @flow
import type { Saga } from 'redux-saga';
import type { tStep, tStepInfo, tStepStatus } from './typeDef';
import type { IWorkable } from '../../workable/IWorkable';

export interface IWorkStep extends IWorkable {
  _failureMsg: string;
  +failureMsg: string;

  _info: ?tStepInfo;
  +info: ?tStepInfo;

  _type: string;
  +type: string;

  _skippable: boolean;
  +skippable: boolean;

  _undoable: boolean;
  +undoable: boolean;

  _image: string;
  +image: string;

  update(stepData: ?$Shape<tStep>): void;

  updateStatus({ status: tStepStatus }): Saga<void>;
}
