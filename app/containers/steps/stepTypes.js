import InputStep from './InputStep';
import ScannerStep from './scannerStep';
import InstructionStep from './instructionStep';
import ScrewStep from './screwStep';
import MaterialStep from './materialStep';
import CheckStep from './CheckStep';
import VideoStep from './VideoStep';
import {stepTypeKeys as s} from '../../modules/step/constants';

export default {
   [s.input]: {
    component: InputStep
  },
  [s.scanner]:{
    component:ScannerStep
  },
  [s.instruction]:{
    component:InstructionStep,
  },
  [s.screw]:{
    component:ScrewStep,
  },
  [s.material]:{
    component:MaterialStep,
  },
  [s.check]:{
    component:CheckStep,
  },
  [s.video]:{
    component:VideoStep,
  }
};

