// @flow
import { ClsOrderOperationPoints } from './classes/ClsOrderOperationPoints';
import type { IWorkStep } from '../interface/IWorkStep';
import type { IScrewStep } from './interface/IScrewStep';
import { onLeave, screwStepStatusTasks } from './stepStatusTasks';


const ScrewStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsScrewStep extends ClsBaseStep implements IScrewStep {
    _tools = [];

    isValid: boolean = false;

    _orderOperationPoints: ClsOrderOperationPoints;

    _pointsToActive = [];

    _listeners = [];

    _statusTasks = screwStepStatusTasks;

    _onLeave = onLeave;

    // eslint-disable-next-line flowtype/no-weak-types
    constructor(...args: Array<any>) {
      super(...args);
      this.isValid = true; // 设置此工步是合法的
      this._onLeave = this._onLeave.bind(this);
    }

    get points() {
      if (this._pointsManager) {
        return this._pointsManager.points;
      }
      return [];
    }
  };
export default ScrewStepMixin;
