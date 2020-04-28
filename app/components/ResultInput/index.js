import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import Radio from '@material-ui/core/Radio';
import RadioGroup from '@material-ui/core/RadioGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import { Typography } from '@material-ui/core';
import {
  clearStories,
  close,
  didMount,
  manualTightening,
  newData,
  resultInput,
  start
} from '../../modules/manual/action';
import withKeyboard from '../Keyboard';
import styles from './style';

const mapStateToProps = (state, ownProps) => ({
  logo: state.logo,
  scanner: state.manual.scanner,
  tool: state.manual.tool,
  pset: state.manual.pset,
  result: state.manual.result,
  timeline: state.manual.timeline,
  ...ownProps
});


const mapDispatchToProps = {
  newData,
  manualTightening,
  start,
  close,
  didMount,
  clearStories,
  resultInput
};


type Props = {
  resultInput: () => void,
  newData: () => void
};

type State = {
  niu: number,
  jao: number,
  ok: string
};

type Result = {
  result: State,
  sucess: boolean
};

function InputLabel({ labelText, classes }) {
  return <Grid item xs={4}>
    <Paper
      className={classes.inputLabel}
      component="div"
      elevation={0}
    >
      <Typography variant="h5">
        {labelText}
      </Typography>
    </Paper>
  </Grid>;
}

function InputItem({ classes, onClick, disabled, content }) {
  return <Grid item xs={8}>
    <Paper
      className={classes.inputContainer}
      component="button"
      onClick={onClick}
      disabled={disabled || false}
    >
      {content}
    </Paper>
  </Grid>;
}

class ConnectedWorking extends React.Component {

  props: Props;

  v1: State;

  r: Result;

  inputTitles = {
    niu: '请输入扭矩值(N·m)：',
    jao: '请输入角度值(°)：'
  };

  constructor(props) {
    super(props);
    this.keyboard = null;
    this.v1 = {
      niu: 0,
      jao: 0,
      ok: ''
    };

    this.r = {
      result: this.v1
    };
  }

  componentDidMount(): void {

  }

  componentWillUnmount(): void {
    const a1 = Number(this.v1.niu);
    const a2 = Number(this.v1.jao);
    this.r.result = this.v1;
    if (isNaN(a1) || isNaN(a2) || a1 <= 0 || a2 < 0 || this.v1.ok === '') {
      this.r.sucess = false;
    } else {
      this.r.sucess = true;
    }
    this.props.resultInput(this.r);
  }

  openManualDiag = (e, inputKey) => {
    e.preventDefault();
    const { keyboardInput } = this.props;
    keyboardInput({
      onSubmit:
        text => {
          this.v1[inputKey] = text;
        },
      text: e.target.value,
      title: this.inputTitles[inputKey],
      label: '请精确到小数点后三位'
    });
  };


  render() {
    const { classes } = this.props;


    return (
      <I18n ns="translations">
        {t => (
          <Grid container spacing={0} className={classes.root}>
            <InputLabel labelText="扭矩值(N·m)：" classes={classes}/>
            <InputItem
              onClick={e => this.openManualDiag(e, 'niu')}
              classes={classes}
              content={<Typography variant="h4">
                {this.v1.niu}
              </Typography>}
            />

            <InputLabel labelText="角度值(°)：" classes={classes}/>

            <InputItem
              onClick={e => this.openManualDiag(e, 'jao')}
              classes={classes}
              content={<Typography variant="h4">
                {this.v1.jao}
              </Typography>}
            />

            <InputLabel labelText="OK/NOK：" classes={classes}/>

            <InputItem
              disabled
              classes={classes}
              content={
                <RadioGroup
                  className={classes.radioGroup}
                  onChange={event => {
                    this.v1.ok = event.target.value;
                  }}>
                  <FormControlLabel
                    value="ok"
                    control={<Radio/>}
                    label={<Typography variant="h4">OK</Typography>}
                  />
                  <FormControlLabel
                    value="nok"
                    control={<Radio/>}
                    label={<Typography variant="h4">NOK</Typography>}
                  />
                </RadioGroup>
              }
            />
          </Grid>
        )}
      </I18n>
    );
  }
}

ConnectedWorking.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  keyboardInput: PropTypes.func.isRequired
};

const Working = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWorking);

// export default Index;
export default withKeyboard(withStyles(styles)(Working), 'numeric');
