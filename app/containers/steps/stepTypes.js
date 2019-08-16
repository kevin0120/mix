import InputStep from './InputStep';
import ScannerStep from './scannerStep';
import InstructionStep from './instructionStep';
import ScrewStep from './screwStep';
import MaterialStep from './materialStep';
import CheckStep from './CheckStep';

export default {
  input: {
    component: InputStep
  },
  scanner:{
    component:ScannerStep
  },
  instruction:{
    component:InstructionStep,
  },
  screw:{
    component:ScrewStep,
  },
  material:{
    component:MaterialStep,
  },
  check:{
    component:CheckStep,
  }
};

