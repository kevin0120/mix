// @flow
import {stepTypeKeys as s} from './constants';
import InputStepMixin from './inputStep/InputStep';
import ScannerStepMixin from './scannerStep/ScannerStep';
import InstructionStepMixin from './instructionStep/InstructionStep';
import ScrewStepMixin from './screwStep/ScrewStep';
import MaterialStepMixin from './materialStep/MaterialStep';
import CheckStepMixin from './checkStep/CheckStep';
import videoStepMixin from './videoStep/videoStep';
import StepMixin from './Step';
import type { IWorkable } from '../workable/IWorkable';

export default {
  [s.input]: (ClsWorkable: Class<IWorkable>)=>InputStepMixin(StepMixin(ClsWorkable)),
  [s.scanner]: (ClsWorkable: Class<IWorkable>)=>ScannerStepMixin(StepMixin(ClsWorkable)),
  [s.instruction]: (ClsWorkable: Class<IWorkable>)=>InstructionStepMixin(StepMixin(ClsWorkable)),
  [s.screw]: (ClsWorkable: Class<IWorkable>)=>ScrewStepMixin(StepMixin(ClsWorkable)),
  [s.material]: (ClsWorkable: Class<IWorkable>)=>MaterialStepMixin(StepMixin(ClsWorkable)),
  [s.check]: (ClsWorkable: Class<IWorkable>)=>CheckStepMixin(StepMixin(ClsWorkable)),
  [s.video]: (ClsWorkable: Class<IWorkable>)=>videoStepMixin(StepMixin(ClsWorkable))
};
