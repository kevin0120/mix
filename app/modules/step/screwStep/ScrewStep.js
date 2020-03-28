// @flow
import { ClsOrderOperationPoints } from './classes/ClsOrderOperationPoints';
import type { IWorkStep } from '../interface/IWorkStep';
import type { IScrewStep } from './interface/IScrewStep';
import { onLeave, screwStepStatusTasksMixin } from './stepStatusTasks';
import type { IWorkable } from '../../workable/IWorkable';


const ScrewStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsScrewStep extends ClsBaseStep implements IScrewStep {
    _tools = [];

    isValid: boolean = false;

    _pointsToActive = [];

    _listeners = [];

    _onLeave = onLeave;

    // eslint-disable-next-line flowtype/no-weak-types
    constructor(...args: Array<any>) {
      super(...args);
      this.isValid = true; // 设置此工步是合法的
      this._onLeave = this._onLeave.bind(this);
      (this: IWorkable)._statusTasks = screwStepStatusTasksMixin((this: IWorkable)._statusTasks);
    }

    get points() {
      if (this._pointsManager) {
        return this._pointsManager.points;
      }
      return [];
    }
  };
export default ScrewStepMixin;
