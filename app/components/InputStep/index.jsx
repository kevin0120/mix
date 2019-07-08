import React from 'react';

export default class InputStep extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      value: ''
    };
  }

  render(): React.ReactElement<any> | string | number | {} | React.ReactNodeArray | React.ReactPortal | boolean | null | undefined {
    const { label, onSubmit, isCurrent } = this.props;
    const { value } = this.state;
    return <div>
      {label}
      <input onChange={(e) => {
        this.setState({
          value: e.target.value
        });
      }}
             value={value}
      />
      <button
        type="button"
        onClick={() => onSubmit(value)}
        disabled={!isCurrent}
      >submit
      </button>
    </div>;
  }
}
