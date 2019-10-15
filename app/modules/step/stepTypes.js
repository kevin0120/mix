// @flow
import {stepTypeKeys as s} from './constants';
import InputStepMixin from './inputStep/InputStep';
import ScannerStepMixin from './scannerStep/ScannerStep';
import InstructionStepMixin from './instructionStep/InstructionStep';
import ScrewStepMixin from './screwStep/ScrewStep';
import MaterialStepMixin from './materialStep/MaterialStep';
import CheckStepMixin from './checkStep/CheckStep';
import videoStepMixin from './videoStep/videoStep';

export default {
  [s.input]: InputStepMixin,
  [s.scanner]: ScannerStepMixin,
  [s.instruction]: InstructionStepMixin,
  [s.screw]: ScrewStepMixin,
  [s.material]: MaterialStepMixin,
  [s.check]: CheckStepMixin,
  [s.video]: videoStepMixin
};
