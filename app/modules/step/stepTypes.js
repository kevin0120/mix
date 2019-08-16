import InputStep from './inputStep/InputStep';
import ScannerStep from './scannerStep/ScannerStep';
import InstructionStep from './instructionStep/InstructionStep';
import ScrewStep from './screwStep/ScrewStep';
import MaterialStep from './materialStep/MaterialStep';
import CheckStep from './checkStep/CheckStep';

export default {
  input: InputStep,
  scanner: ScannerStep,
  instruction: InstructionStep,
  screw: ScrewStep,
  material: MaterialStep,
  check: CheckStep
};
