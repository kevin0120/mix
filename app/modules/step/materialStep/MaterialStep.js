// @flow
import type { IWorkStep } from '../interface/IWorkStep';
import type { IMaterialStep } from './interface/IMaterialStep';
import { materialStepStatusMixin, onLeave } from './stepStatusTasks';

const MaterialStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsMaterialStep extends ClsBaseStep implements IMaterialStep {
    _ports = new Set([]);

    _io = new Set([]);

    _items = new Set([]);

    _confirm = null;

    _onLeave = onLeave;

    constructor(...args) {
      super(...args);
      this._statusTasks = materialStepStatusMixin(this._statusTasks);
    }
  };

export default MaterialStepMixin;
