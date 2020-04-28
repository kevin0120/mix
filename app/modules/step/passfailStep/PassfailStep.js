// @flow
import type { IWorkStep } from '../interface/IWorkStep';
import { passfailStepStatusMixin } from './stepStatusTasks';

const PassFailStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsPassFailStep extends ClsBaseStep {
    constructor(...args) {
      super(...args);
      this._statusTasks = passfailStepStatusMixin(this._statusTasks);
    }
  };
export default PassFailStepMixin;
