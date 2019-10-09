import InputStepMixin from './inputStep/InputStep';
import ScannerStepMixin from './scannerStep/ScannerStep';
import InstructionStepMixin from './instructionStep/InstructionStep';
import ScrewStepMixin from './screwStep/ScrewStep';
import MaterialStepMixin from './materialStep/MaterialStep';
import CheckStepMixin from './checkStep/CheckStep';
import videoStepMixin from './videoStep/videoStep';

export default {
  input: InputStepMixin,
  scanner: ScannerStepMixin,
  instruction: InstructionStepMixin,
  screw: ScrewStepMixin,
  material: MaterialStepMixin,
  check: CheckStepMixin,
  video: videoStepMixin
};
