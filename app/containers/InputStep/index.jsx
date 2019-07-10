import React from 'react';
import { connect } from 'react-redux';
import { inputStepActions } from '../../modules/inputStep/action';

type Props = {
  label: string,
  isCurrent: boolean,
  submit: ()=>{}
};

class InputStep extends React.Component<Props> {
  constructor(props) {
    super(props);
    this.state = {
      value: ''
    };
  }

  onSubmit = (value) => {
    const { submit } = this.props;
    this.setState({
      value: ''
    });
    submit(value);

  };

  // eslint-disable-next-line flowtype/no-weak-types
  render(): React.ReactElement<any> | string | number | {} | React.ReactNodeArray | React.ReactPortal | boolean | null | undefined {
    const { label, isCurrent } = this.props;
    const { value } = this.state;
    return <div>
      {label}
      <input
        onChange={(e) => {
          this.setState({
            value: e.target.value
          });
        }}
        value={value}
      />
      <button
        type="button"
        onClick={() => this.onSubmit(value)}
        disabled={!isCurrent}
      >submit
      </button>
    </div>;
  }
}


const mapState = (state, props) => ({
  order: state.order,
  ...props
});

const mapDispatch = {
  submit: inputStepActions.submit
};

export default connect(mapState, mapDispatch)(InputStep);
