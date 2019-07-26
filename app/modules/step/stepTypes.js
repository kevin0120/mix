import inputStep from './inputStep/saga';
import scannerStep from './scannerStep/saga';
import instructionStep from './instructionStep/saga';

export default {
  input: inputStep,
  scanner: scannerStep,
  instruction:instructionStep
};
