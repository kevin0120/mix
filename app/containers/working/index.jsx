import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import List from '@material-ui/core/List';

import { I18n } from 'react-i18next';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import Paper from '@material-ui/core/Paper';
import FiberManualRecord from '@material-ui/icons/FiberManualRecord';
import Autorenew from '@material-ui/icons/Autorenew';
import Button from '../../components/CustomButtons/Button';
import ImageStick from '../../components/ImageStick/imageStick';

import ShutdownDiag from '../../components/ShutDownDiag';
import { ScannerNewData } from '../../modules/scanner/action';
import { switchWorkMode } from '../../modules/workmode/action';

import {
  switch2Timeout,
  operationBypassConfirm,
  operationBypassCancel,
  operationConflictConfirm,
  operationConflictCancel
} from '../../modules/operation/action';

import configs from '../../shared/config';


// components

import WorkingInfoBar from '../../components/WorkingInfoBar';

import {
  container,
  cardTitle,
  description,
  dangerColor,
  successColor
} from '../../common/jss/material-react-pro';
import ResultDialog from '../../components/ResultDialog';
import ManualDiag from '../../components/ManualDiag';

import TimeLine from '../../components/WorkPageTimeline';

import ProgressBar from '../../components/ProgressBar/Progress';
import { OPERATION_STATUS, OPERATION_SOURCE } from '../../modules/operation/model';
import withKeyboard from '../../components/Keyboard';

import toolSvg from './toolSvg';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  operations: state.operations,
  logo: state.logo,
  operationSettings: state.setting.operationSettings,
  workMode: state.workMode,
  timeline: state.timeline,
  reworkWorkCenter: state.setting.system.workcenter.rework_workcenter,
  enableFocus: state.setting.systemSettings.enableFocus,
  enableConflictOP: state.setting.systemSettings.enableConflictOP,
  tools: state.tools,
  ...ownProps
});

const mapDispatchToProps = {
  NewCar: ScannerNewData,
  switchWorkMode,
  switch2Timeout,
  doConfirmBypass: operationBypassConfirm,
  doCancelBypass: operationBypassCancel,
  doConfirmConflict: operationConflictConfirm,
  doCancelConflict: operationConflictCancel

};

// 与 style 里的变量相同
// css 覆盖不了的 放这里
const withstyles = theme => ({
  root: {
    margin: 0,
    padding:0,
    width:'100%',
    height:'100%',
    background: '#EFF4F7'
  },
  divider: {
    margin: '5px 10px'
  },
  progressWrap: {
    height: '100%',
    position: 'relative',
    padding: '0px'
  },
  fab: {
    position: 'absolute',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
    width: '150px'
  },
  fabOEE: {
    position: 'fixed',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
    width: '150px'
  },
  extendedIcon: {
    marginRight: theme.spacing()
  },
  cardDescription: {
    ...description,
    fontSize: '45px',
    fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
    textAlign: 'center'
  },
  LeftTopTab: {
    textAlign: 'left',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    width: '100%',
    height: '100%'
  },
  LeftBottomTab: {
    textAlign: 'left',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    height: '100%',
    width: '100%'
  },
  LeftTabContiner: {
    height: '100%',
    width: '100%',
    padding: '2px 10px',
    position: 'relative'
  },
  LeftTopDes: {
    marginBottom: '0',
    color: '#979797',
    '& p': {
      fontSize: '14px'
    }
  },
  CutDownPaper: {
    textAlign: 'center',
    position: 'relative',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100%',
    width: '100%',
    background: '#212121'
  },
  InfoTab: {
    textAlign: 'center',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    height: '100%',
    width: '100%'
  },
  InfoTabTimeLine: {
    textAlign: 'center',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    height: '100%',
    width: '100%'
  },
  CutDownContainer: {
    position: 'absolute',
    height: '100%',
    marginTop: '0px',
    width: '100%',
    color: '#ffeb3b',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  },
  CountDownContainer: {
    position: 'absolute',
    height: '90%',
    // marginTop: '10px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '100%',
    color: '#E5F0FA'
  },
  TurnPaper: {
    textAlign: 'center',
    height: '100%',
    width: '100%',
    position: 'relative',
    background: '#212121'
  },
  RetryPaper: {
    textAlign: 'center',
    height: '100%',
    width: '100%',
    position: 'relative',
    background: '#212121'
  },
  RightDescription: {
    fontSize: '2vh',
    marginBottom: '0'
  },
  RightNum: {
    fontSize: '6vh',
    fontWeight: '600',
    padding: 0,
    margin: 0
  },
  MarginTopBottom5: {
    margin: '0 0 5px'
  }
});

class ConnectedWorking extends React.Component {
  constructor(props) {
    super(props);

    this.autoCancel = null;
    this.keyboard = null;

  }

  shouldComponentUpdate() {
    return true;
  }

  componentDidUpdate(prevProps) {
    this.prevOperationStatus = prevProps.operations.operationStatus;
  }

  openManualDiag = (e, input, t) => {
    e.preventDefault();
    const { keyboardInput, NewCar } = this.props;
    keyboardInput({
      onSubmit: text => {
        NewCar(text, OPERATION_SOURCE.MANUAL);
      },
      text: e.target.value,
      title: t(`Operation.Input.${input}`),
      label: t(`Operation.Input.${input}`)
    });
  };

  workSiteInfo = t => {
    const { reworkWorkCenter } = this.props;
    return [
      {
        key: 'rework',
        value: reworkWorkCenter,
        displayTitle: t('Operation.Info.ReworkWorkCenter')
      }
    ];
  };

  toggleAutoScannerMode = e => {
    const { workMode, switchWorkMode } = this.props;
    let mode = 'auto';
    if (workMode.workMode === 'auto') {
      mode = 'scanner';
    }
    switchWorkMode(mode);
  };

  orderInfo = t => {
    const { operations, workMode, tools } = this.props;

    // const showButtonInfo = isAutoMode? 'Common.Auto':'Common.Manual';
    let programme = operations.jobID.toString();
    if (workMode.controllerMode === 'pset') {
      if (operations.results.length > 0) {
        programme = operations.results[
          operations.activeResultIndex
          ].pset.toString();
      }
    }

    return [
      {
        key: 'WorkMode',
        value: t(`Operation.Info.WorkMode.${workMode.workMode}`),
        displayTitle: t('Operation.Info.WorkMode.Title')
      },
      {
        key: 'WorkStatus',
        value: t(`Operation.Display.${operations.operationStatus}`),
        displayTitle: t('Operation.Info.WorkStatus.Title')
      },
      // {
      //   key: '控制模式',
      //   value: ControlMode,
      //   displayTitle: '控制模式',
      // },
      {
        key: 'Job',
        value: programme,
        displayTitle: t('Operation.Info.Job.Title')
      },
      {
        key: 'Gun',
        value: toolSvg(tools.status === 'connected' ? successColor : dangerColor, { width: '4vh', height: '4vh' }),
        displayTitle: t('Operation.Info.Tool.Title')
      }
    ];
  };

  confirmBypass = () => {
    const { doConfirmBypass } = this.props;
    doConfirmBypass();
  };

  cancelBypass = () => {
    const { doCancelBypass } = this.props;
    doCancelBypass();
  };

  bypassDiag = (t) => ({
    show: true,
    title: t('Common.Bypass'),
    onConfirm: this.confirmBypass,
    onCancel: this.cancelBypass,
    content: (
      <div>
        <h3>{t('Common.QuestBypass')} </h3>
        <h4 style={{ textAlign: 'center' }}> OR </h4>
        {/* <Button onClick={() => this.handleReset} color="danger" round> */}
        {/* {t('Common.Reset')} */}
        {/* </Button> */}
      </div>
    ),
    showCancel: true
  });

  confirmConflict = () => {
    const { doConfirmConflict, operations } = this.props;
    if (this.autoCancel) {
      clearTimeout(this.autoCancel);
      this.autoCancel = null;
    }
    doConfirmConflict(operations.conflict.data);
  };

  cancelConflict = () => {
    const { doCancelConflict } = this.props;
    if (this.autoCancel) {
      clearTimeout(this.autoCancel);
      this.autoCancel = null;
    }
    doCancelConflict();
  };

  conflictDiag = (t) => {
    const { enableConflictOP, operations, doCancelConflict } = this.props;
    if (!this.autoCancel) {
      this.autoCancel = setTimeout(() => {
        doCancelConflict();
      }, 5000);
    }

    if (enableConflictOP) {
      return {
        show: true,
        title: t('Common.VerifyCar'),
        onConfirm: this.confirmConflict,
        onCancel: this.cancelConflict,
        content: t('Common.QuestConfirmCarInfo'),
        showCancel: true
      };
    }
    return {
      show: true,
      title: t('Common.VerifyCar'),
      onConfirm: this.cancelConflict,
      onCancel: this.cancelConflict,
      content: `设定为不允许重复拧紧同一张工单 VIN:${operations.conflict.data.vin}`,
      showCancel: false
    };
  };

  render() {
    const {
      classes,
      operations,
      timeline,
      workMode,
      enableFocus,
      logo
    } = this.props;

    let batch = '0/0';
    let redoBatch = '0/0';
    let maxRedoTimes = 0;
    if (operations.results.length > 0) {
      maxRedoTimes =
        operations.results[operations.activeResultIndex].max_redo_times;
      batch = `${(
        operations.activeResultIndex + 1
      ).toString()}/${operations.results[
      operations.results.length - 1
        ].group_sequence.toString()}`;
      redoBatch = `${(
        maxRedoTimes - operations.failCount
      ).toString()}/${maxRedoTimes.toString()}`;
    }

    if (operations.operationStatus === 'Ready') {
      redoBatch = '0';
    }

    const { enableResultDialog } = configs.operationSettings;

    const showResultDiag = configs.operationSettings.opMode === 'order';

    // const showManualCarType = configs.operationSettings.opMode === 'op';

    const carTypeSize = configs.operationSettings.opMode === 'op' ? 5 : 3;

    const manualEnable =
      (configs.operationSettings.opMode === 'op' &&
        workMode.workMode === 'manual') ||
      configs.operationSettings.opMode === 'order';


    const fabClassName = classNames(classes.fab);
    // const fabOEEClassName = switchAutoManual? classNames(classes.fabOEE) :classNames(classes.fab);
    // let statusShow = lodash.includes(['Init', 'Ready', 'PreDoing', 'Timeout'],orderStatus);
    //
    // const enableSwitchMode = workFlow === 'General';
    //
    // if (!isAutoMode) {
    //   statusShow = true;
    // }
    //
    // const oeeShow = oeeFuncEnable && lodash.includes(['Doing'],orderStatus);
    //
    // const resumeShow = lodash.includes(['Pending'],orderStatus);
    //
    // const showAutoManual = statusShow && switchAutoManual;
    //
    const showButtonInfo =
      workMode.workMode === 'auto' ? 'Common.Scanner' : 'Common.Auto';
    //
    // const showManualInfo = !isAutoMode || workFlow === 'General';
    //
    // const number = maxOpTime;
    //
    // let shouldProcessing = true;
    // if (this.props.orderStatus === 'Ready' || this.props.orderStatus === 'PreDoing' || this.props.orderStatus === 'Timeout' || this.props.orderStatus === 'Init' || (!this.props.isAutoMode)){
    //   shouldProcessing = false;
    // }
    return (
      <I18n ns="translations">
        {t => (
          <Grid container spacing={1} className={classes.root} justify="center">
            <Grid item xs={9} container spacing={1}>
              {showResultDiag ? (
                <Grid item xs={2} style={{ height: '13%' }}>
                  <Paper
                    className={classes.LeftTopTab}
                    component="button"
                    disabled
                  >
                    <div className={classes.LeftTabContiner}>
                      <h4 className={classes.LeftTopDes}>
                        <p className={classes.MarginTopBottom5}>
                          车序:
                        </p>
                      </h4>
                      <p className={classes.cardDescription}>
                        {operations.lnr}
                      </p>
                    </div>
                  </Paper>
                </Grid>
              ) : null}
              <Grid item xs={carTypeSize} style={{ height: '13%' }}>
                <Paper
                  className={classes.LeftTopTab}
                  component="button"
                  onClick={e => this.openManualDiag(e, 'carType', t)}
                  disabled={
                    configs.operationSettings.opMode === 'order' ||
                    configs.operationSettings.opMode === 'op'
                  }
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        {t('Operation.Input.carType')}:
                      </p>
                    </h4>
                    <p className={classes.cardDescription}>
                      {operations.carType}
                    </p>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={7} style={{ height: '13%' }}>
                <Paper
                  className={classes.LeftTopTab}
                  component="button"
                  onClick={e => this.openManualDiag(e, 'vin', t)}
                  disabled={!manualEnable}
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        {t('Operation.Input.vin')}:
                      </p>
                    </h4>
                    <p className={classes.cardDescription}>
                      {operations.carID}
                    </p>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '87%' }}>
                <Paper className={classes.LeftBottomTab}>
                  <ImageStick
                    logo={logo}
                    operations={operations}
                    enableFocus={enableFocus}
                  >
                    {!manualEnable && configs.operationSettings.opMode === 'op' ? (
                      <Button
                        color="rose"
                        disabled={manualEnable}
                        round
                        className={fabClassName}
                        size="lg"
                        onClick={this.toggleAutoScannerMode}
                      >
                        {workMode.workMode === 'auto' ? (
                          <FiberManualRecord className={classes.extendedIcon}/>
                        ) : (
                          <Autorenew className={classes.extendedIcon}/>
                        )}
                        {t(showButtonInfo)}
                      </Button>
                    ) : null}
                  </ImageStick>
                </Paper>
              </Grid>
            </Grid>
            <Grid item xs={3} container spacing={1} alignContent="flex-start" alignItems="flex-start">
              <Grid item xs={12} style={{ height: '18%' }}>
                <Paper className={classes.CutDownPaper}>
                  <div className={classes.CountDownContainer}>
                    <ProgressBar
                      time={operations.maxOpTimes}
                      shouldCounterStart={() =>
                        // predoing -> doing
                        operations.operationStatus === OPERATION_STATUS.DOING &&
                        this.prevOperationStatus === OPERATION_STATUS.PREDOING
                      }
                      shouldCounterStop={() =>
                        lodash.includes(
                          [OPERATION_STATUS.READY],
                          operations.operationStatus
                        )
                      }
                      shouldCounterReady={() =>
                        operations.operationStatus === OPERATION_STATUS.PREDOING
                      }
                      // onFinish={() => this.props.switch2Timeout()}
                      gridClassName={classes.progressWrap}
                    />
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={6} style={{ height: '15%' }}>
                <Paper className={classes.TurnPaper}>
                  <div className={classes.CutDownContainer}>
                    <div style={{
                      height: '20%',
                      justifyContent: 'center',
                      alignItems: 'center',
                      display: 'flex',
                      flexDirection: 'column'
                    }}>
                      <p className={classes.RightDescription}>
                        拧紧点
                      </p>
                    </div>
                    <div style={{
                      height: '60%',
                      justifyContent: 'center',
                      alignItems: 'center',
                      display: 'flex',
                      flexDirection: 'column'
                    }}>

                      <p className={classes.RightNum}>{batch}</p>
                    </div>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={6} style={{ height: '15%' }}>
                <Paper className={classes.RetryPaper}>
                  <div className={classes.CutDownContainer}>
                    <div style={{
                      height: '20%',
                      justifyContent: 'center',
                      alignItems: 'center',
                      display: 'flex',
                      flexDirection: 'column'
                    }}>
                      <p className={classes.RightDescription}>
                        可拧紧次数
                      </p>
                    </div>
                    <div style={{
                      height: '60%',
                      justifyContent: 'center',
                      alignItems: 'center',
                      display: 'flex',
                      flexDirection: 'column'
                    }}>

                      <p className={classes.RightNum}>{redoBatch}</p>
                    </div>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '27%' }}>
                <Paper className={classes.InfoTab}>
                  <List style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
                    <WorkingInfoBar
                      key="infoUser"
                      infos={this.workSiteInfo(t)}
                    />
                    <Divider
                      className={classes.divider}
                      key="divider-infoUser"
                      light
                    />
                    <WorkingInfoBar key="infoOrder" infos={this.orderInfo(t)}/>
                  </List>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '40%' }}>
                <Paper className={classes.InfoTabTimeLine}>
                  <TimeLine simple stories={timeline}/>
                  {/* <TimeLine simple stories={teststory} /> */}
                </Paper>
              </Grid>
            </Grid>

            {operations.bypassToConfirm ?
              <ShutdownDiag
                {...this.bypassDiag(t)}
              /> : null
            }
            {operations.conflict.toConfirm ?
              <ShutdownDiag
                {...this.conflictDiag(t)}
              /> : null
            }
            {enableResultDialog ? <ResultDialog/> : null}
          </Grid>
        )}
      </I18n>
    );
  }
}

ConnectedWorking.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  operations: PropTypes.shape({}).isRequired,
  operationSettings: PropTypes.shape({}).isRequired,
  workMode: PropTypes.shape({}).isRequired,
  timeline: PropTypes.arrayOf(PropTypes.shape({})).isRequired,
  reworkWorkCenter: PropTypes.string.isRequired,
  keyboardInput: PropTypes.func.isRequired,
  switchWorkMode: PropTypes.func.isRequired,
  switch2Timeout: PropTypes.func.isRequired,
  NewCar: PropTypes.func.isRequired
};

const Working = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWorking);

export default withKeyboard(withStyles(withstyles)(Working));
