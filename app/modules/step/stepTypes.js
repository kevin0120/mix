import inputStep from './inputStep/saga';
import scannerStep from './scannerStep/saga';
import instructionStep from './instructionStep/saga';
import screwStep from './screwStep/saga';
import materialStep from './materialStep/saga';

export default {
  input: inputStep,
  scanner: scannerStep,
  instruction: instructionStep,
  screw: screwStep,
  material: materialStep
};
