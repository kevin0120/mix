import InputStep from '../InputStep';
import ScannerStep from '../ScannerStep';

export default {
  input: {
    component: InputStep,
    props: (props) => ({
      onSubmit: (value) => props.pushState(value, props.parallelId),
      label: props.payload.label
    })
  },
  scanner:{
    component:ScannerStep
  }
};

