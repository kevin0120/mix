import InputStep from './InputStep';
import ScannerStep from './scannerStep';
import InstructionStep from './instructionStep';
import ScrewStep from './screwStep';

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
    genProps:props=>({
      points:props.payload.points,
      image:props.payload.image
    })
  }
};

