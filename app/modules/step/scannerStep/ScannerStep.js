// @flow
import type { IWorkStep } from '../interface/IWorkStep';
import { onLeave, scannerStepStatusMixin } from './stepStatusTasks';

const ScannerStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsScannerStep extends ClsBaseStep {
    _scanners = [];

    _listeners = [];

    _onLeave = onLeave;

    constructor(...args) {
      super(...args);
      this._statusTasks = scannerStepStatusMixin(this._statusTasks);
    }
  };
export default ScannerStepMixin;
