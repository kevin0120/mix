import InputStep from './InputStep';
import ScannerStep from './scannerStep';
import InstructionStep from './instructionStep';

export default {
  input: {
    component: InputStep,
    genProps: (props) => ({
      onSubmit: (value) => props.pushState(value, props.parallelId),
      label: props.payload.label
    })
  },
  scanner:{
    component:ScannerStep,
    genProps:props=>({
      label:props.payload.label
    })
  },
  instruction:{
    component:InstructionStep,
    genProps:props=>({
      instruction:props.payload.instruction
    })
  }
};

