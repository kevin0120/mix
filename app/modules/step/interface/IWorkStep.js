// @flow
import type {
  tStep,
  tStepInfo
} from './typeDef';
import type { IWorkable } from '../../workable/IWorkable';

export interface IWorkStep extends IWorkable {

  _info: ?tStepInfo,
  +info: ?tStepInfo,

  _type: string,
  +type: string,

  _skippable: boolean,
  +skippable: boolean,

  _undoable: boolean,
  +undoable: boolean,

  _image: string,
  +image: string,

  update(stepData: ?$Shape<tStep>): void
}
