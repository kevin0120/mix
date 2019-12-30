// @flow
import { stepTypeKeys as s } from './constants';
import InputStepMixin from './inputStep/InputStep';
import ScannerStepMixin from './scannerStep/ScannerStep';
import ScrewStepMixin from './screwStep/ScrewStep';
import MaterialStepMixin from './materialStep/MaterialStep';
import PassFailStepMixin from './PassFailStep/PassFailStep';
import StepMixin from './Step';
import type { IWorkable } from '../workable/IWorkable';

export default {
  [s.input]: (ClsWorkable: Class<IWorkable>) => InputStepMixin(StepMixin(ClsWorkable)),
  [s.scanner]: (ClsWorkable: Class<IWorkable>) => ScannerStepMixin(StepMixin(ClsWorkable)),
  [s.instruction]: (ClsWorkable: Class<IWorkable>) => StepMixin(ClsWorkable),
  [s.text]: (ClsWorkable: Class<IWorkable>) => StepMixin(ClsWorkable),
  [s.screw]: (ClsWorkable: Class<IWorkable>) => ScrewStepMixin(StepMixin(ClsWorkable)),
  [s.material]: (ClsWorkable: Class<IWorkable>) => MaterialStepMixin(StepMixin(ClsWorkable)),
  [s.passFail]: (ClsWorkable: Class<IWorkable>) => PassFailStepMixin(StepMixin(ClsWorkable)),
  [s.measure]: (ClsWorkable: Class<IWorkable>) => StepMixin(ClsWorkable),
  [s.video]: (ClsWorkable: Class<IWorkable>) => StepMixin(ClsWorkable)
};
