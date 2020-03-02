import { connect } from 'react-redux';
import React from 'react';
import { Link } from 'react-router-dom';
import { bindActionCreators } from 'redux';
import Keyboard from 'react-simple-keyboard';
import styles from './Counter.css';
import * as CounterActions from '../../modules/manual/action';

type Props = {
  start: () => void,
  close: () => void,
  manualTightening: () => void,
  manualCancel: () => void,
  newData: () => void,
  didMount: () => void,
  scanner: string,
  tool: string,
  pset: number,
  result: Array<string>
};
class Manualpage extends React.Component<Props>{

  props: Props;

  state = {
    layoutName: "default",
    inputName:"",
    input: ""
  };

  componentDidMount(): void {
    this.props.start();
    this.props.didMount();
  }

  componentWillUnmount(): void {
    this.props.close();
  }


  onChange = input => {
    this.setState({ input });
  };

  onKeyPress = button => {

    /**
     * If you want to handle the shift and caps lock buttons
     */
    if (button === "{enter}") {
      this.props.newData(this.state.input);
      this.removeActiveInput();
    }

    if (button === "{shift}" || button === "{lock}"){
      this.handleShift();
    }

  };

  handleShift = () => {
    const {layoutName} = this.state;

    this.setState({
      layoutName: layoutName === "default" ? "shift" : "default"
    });
  };

  onFocusInput = event => {
    this.setState({
      inputName: event.target.id
    });
  };


  removeActiveInput = () => {
    this.setState({
      inputName: "",
      // input:"",
    });
  };


  onChangeInput = event => {
    const input = event.target.value;
    this.setState({ input });
    this.keyboard.keyboard.setInput(input);
  };

  render() {
  const {
    manualCancel,
    manualTightening,
    scanner,
    tool,
    pset,
    result
  } = this.props;

  let result1: Array<string>;

  if (result.length>15){
    result1  = result.slice(result.length-16,result.length)
  } else {
    result1  = result
  }

    return (
      <div className={styles.backGround}>
        <div className={styles.backButton} data-tid="backButton">
          <Link to="/app">
            <i className="fa fa-arrow-left fa-3x"/>
          </Link>
        </div>

        <div className={`counter ${styles.counter}`} data-tid="counter">
          {JSON.stringify(result1)}
        </div>

        <div className={`counter ${styles.counter1}`} data-tid="counter">
          拧紧工具:{tool}
        </div>
        <div className={`counter ${styles.counter2}`} data-tid="counter">
          程序号:{pset}
        </div>
        <div className={`counter ${styles.counter3}`} data-tid="counter">
          追溯码:{scanner}
        </div>
        <div className={`counter ${styles.counter4}`} data-tid="counter">手动输入:
          <input value={this.state.input}
                 id="input"
                 className={styles.input}
                 onFocus={this.onFocusInput}
                 onChange={this.onChangeInput} />
        </div>
        {this.state.inputName.length !== 0 ? (
        <div className={`counter ${styles.counter5}`} data-tid="counter">
          <Keyboard
            ref={r => (this.keyboard = r)}
            onChange={this.onChange}
            onKeyPress={this.onKeyPress}
            layoutName={this.state.layoutName}
          />
        </div>
          ) : null}
        <div className={styles.btnGroup}>

          <button
            className={styles.btn}
            onClick={()=>{
              manualCancel();
              this.removeActiveInput();}
            }
            data-tclass="btn"
            type="button"
          >
            选择工具
          </button>
          <button
            className={styles.btn}
            onClick={()=>{
              manualTightening();
              this.removeActiveInput();
            }
            }
            data-tclass="btn"
            type="button"
          >
            开始作业
          </button>
        </div>
      </div>
    );
  }
}

const mapState = (state, props) => ({
  counter: state.manual.counter,
  scanner: state.manual.scanner,
  tool: state.manual.tool,
  pset: state.manual.pset,
  result:state.manual.result,
  ...props
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators(CounterActions, dispatch);
}

export default (connect(mapState, mapDispatchToProps)(Manualpage));
