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


class ConnectedWorking extends React.Component {

  props: Props;

  v1: State;

  r: Result;

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

  openManualDiag = (e, butt) => {
    e.preventDefault();
    const { keyboardInput } = this.props;
    if (butt === 1) {
      keyboardInput({
        onSubmit:
          text => {
            this.v1.niu = text;
          },
        text: e.target.value,
        title: '请输入扭矩值：',
        label: '请精确到小数点后三位'
      });
    } else {
      keyboardInput({
        onSubmit:
          text => {
            this.v1.jao = text;
          },
        text: e.target.value,
        title: '请输入角度值：',
        label: '请精确到小数点后三位'
      });

    }
  };


  render() {
    const {
      classes
    } = this.props;


    return (
      <I18n ns="translations">
        {t => (
          <Grid container spacing={0} className={classes.root} justify="center">
            <Grid item xs={3} container style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 1)}
                disabled
              >
                <div className={classes.LeftTabContiner}>
                  <h4 className={classes.LeftTopDes}>
                    <p className={classes.MarginTopBottom5}>
                      扭矩值:
                    </p>
                  </h4>
                </div>
              </Paper>
            </Grid>

            <Grid item xs={9} container style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 1)}
                disabled={false}
              >
                <div className={classes.LeftTabContiner}>
                  <p className={classes.cardDescription}>
                    {this.v1.niu}
                  </p>
                </div>
              </Paper>
            </Grid>


            <Grid item xs={3} container style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 2)}
                disabled
              >
                <div className={classes.LeftTabContiner}>
                  <h4 className={classes.LeftTopDes}>
                    <p className={classes.MarginTopBottom5}>
                      角度值：
                    </p>
                  </h4>
                </div>
              </Paper>
            </Grid>

            <Grid item xs={9} container style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 2)}
                disabled={false}
              >
                <div className={classes.LeftTabContiner}>
                  <p className={classes.cardDescription}>
                    {this.v1.jao}
                  </p>
                </div>
              </Paper>
            </Grid>


            <Grid item xs={3} container style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 2)}
                disabled
              >
                <div className={classes.LeftTabContiner}>
                  <h4 className={classes.LeftTopDes}>
                    <p className={classes.MarginTopBottom5}>
                      OK/NOk：
                    </p>
                  </h4>

                </div>
              </Paper>
            </Grid>

            <Grid item xs={9} style={{ height: '30%' }}>
              <Paper
                className={classes.LeftTopTab}
                component="button"
                onClick={e => this.openManualDiag(e, 2)}
                disabled
              >
                <div className={classes.LeftTabContiner}>

                  <RadioGroup aria-label="gender" name="gender1" className={classes.cardDescription}
                              onChange={event => {
                                this.v1.ok = event.target.value;
                              }}>
                    <FormControlLabel value="ok" control={<Radio/>} label="OK" className={classes.cardDescription}/>
                    <FormControlLabel value="nok" control={<Radio/>} label="NOK"/>
                  </RadioGroup>
                </div>
              </Paper>
            </Grid>

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
