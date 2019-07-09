import InputStep from '../../components/InputStep';
import Working from '../working';

export default {
  input: {
    component: InputStep,
    props: (props) => ({
      onSubmit: (value) => props.pushState(value, props.parallelId),
      label: props.payload.label
    })
  }
};

