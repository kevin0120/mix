// @flow
import type {tScrewStepPayload} from './typeDef';
import type { IWorkStep } from '../../interface/IWorkStep';

export interface IScrewStep extends IWorkStep{
  _payload: ?tScrewStepPayload
}