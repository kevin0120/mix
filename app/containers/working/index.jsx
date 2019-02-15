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
import { NewCar } from '../../actions/scannerDevice';
import { switchWorkMode } from '../../actions/workMode';

import {
  switch2Timeout,
  operationBypassConfirm,
  operationBypassCancel,
  operationConflictConfirm,
  operationConflictCancel
} from '../../actions/operation';

import configs from '../../shared/config/index';


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
import { OPERATION_STATUS, OPERATION_SOURCE } from '../../reducers/operations';
import withKeyboard from '../../components/Keyboard';

import toolSvg from './toolSvg';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  operations: state.operations,
  logo: state.logo,
  operationSettings: state.setting.operationSettings,
  workMode: state.workMode,
  timeline: state.timeline,
  reworkWorkCenter: state.connections.rework_workcenter,
  enableFocus: state.setting.systemSettings.enableFocus,
  enableConflictOP: state.setting.systemSettings.enableConflictOP,
  tools: state.tools,
  ...ownProps
});

const TOPHEIGHT = '150px';

const mapDispatchToProps = {
  NewCar,
  switchWorkMode,
  switch2Timeout,
  doConfirmBypass: operationBypassConfirm,
  doCancelBypass: operationBypassCancel,
  doConfirmConflict: operationConflictConfirm,
  doCancelConflict: operationConflictCancel

};

// 与 style 里的变量相同
// const TOPHEIGHT = '150px';
// css 覆盖不了的 放这里
const withstyles = theme => ({
  container: {
    ...container,
    zIndex: '4',
    [theme.breakpoints.down('sm')]: {
      paddingBottom: '100px'
    }
  },
  cardTitle,
  content: {
    flex: 1,
    overflow: 'hidden',
    flexWrap: 'nowrap'
  },
  root: {
    // position: 'relative',
    height: '100%',
    width: '100%',
    margin: '0',
    background: '#EFF4F7'
  },
  infoWrap: {
    fontSize: 14,
    color: '#333',
    position: 'relative',
    transition: 'all 225ms cubic-bezier(0, 0, 0.2, 1) 0ms',
    width: 200,
    overflowY: 'auto'
  },
  transfromInfo: {
    width: '0'
  },
  divider: {
    margin: '5px 10px'
  },
  drawerPaper: {
    position: 'relative'
  },
  toolbar: theme.mixins.toolbar,
  row: {
    display: 'flex',
    margin: '10px 0'
  },
  avatar: {
    marginRight: 10,
    width: 50,
    height: 50
  },
  userInfo: {
    color: '#333',
    fontSize: 12,
    padding: 0
  },
  userText: {
    fontSize: 12,
    paddingLeft: '10px'
  },
  timeWrap: {
    padding: '10px 5px'
  },
  timeContent: {
    margin: '10px 10px 0px',
    fontSize: 20
  },
  baseInfo: {
    boxSizing: 'border-box',
    position: 'absolute',
    width: '100%',
    padding: '10px 20px 20px',
    background: 'transparent',
    bottom: 0,
    left: 0
  },
  progressWrap: {
    height: '100%',
    position: 'relative',
    padding: '0px'
  },
  topWrap: {
    boxShadow:
      '0px 2px 4px -1px rgba(0, 0, 0, 0.2), 0px 4px 5px 0px rgba(0, 0, 0, 0.14), 0px 1px 10px 0px rgba(0, 0, 0, 0.12)',

    height: TOPHEIGHT
  },
  zoomBtn: {
    position: 'absolute',
    height: 30,
    width: 30,
    minHeight: 30,
    zIndex: 99,
    left: 0,
    bottom: 30,
    boxShadow: '0 2px 7px rgba(0, 0, 0, .8)'
    // background: '#fff',
  },
  exitIcon: {
    fontSize: 26
  },
  zoomOut: {
    fontSize: 18
  },
  fab: {
    position: 'absolute',
    bottom: theme.spacing.unit * 2,
    right: theme.spacing.unit * 2,
    width: '150px'
  },
  fabOEE: {
    position: 'fixed',
    bottom: theme.spacing.unit * 15,
    right: theme.spacing.unit * 2,
    width: '150px'
  },
  fabResume: {
    position: 'fixed',
    bottom: theme.spacing.unit * 15,
    right: theme.spacing.unit * 2,
    width: '200px'
  },
  extendedIcon: {
    marginRight: theme.spacing.unit
  },
  cardVehicleSeq: {
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    marginLeft: '10px',
    marginBottom: '0px',
    width: '120px'
  },
  cardVehicleType: {
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    marginBottom: '0px',
    width: '310px',
    marginLeft: '60px'
  },
  cardVehicleVIN: {
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    marginBottom: '0px',
    width: '480px',
    marginLeft: '145px'
  },
  cardNormal: {
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    marginBottom: '0px'
  },
  cardCountdown: {
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    marginBottom: '0px',
    marginLeft: '5px',
    width: '310px'
  },
  cardBodyNormal: {
    padding: '0',
    margin: '0',
    height: '100%'
  },
  LeftContainer: {
    height: '100%',
    width: '75%'
  },
  RightContainer: {
    height: '100%',
    width: '25%'
  },
  RightContent: {
    height: '100%',
    marginTop: '10px'
  },
  InfoBarGrid: {
    marginLeft: '20px',
    height: '100px'
  },
  InfoBarGridContainer: {
    marginTop: '10px',
    width: '100%',
    height: '100px'
  },
  ImageStickGrid: {
    height: '700px'
  },
  ImageStickGridContainer: {
    height: 'calc(100% - 100px)',
    marginTop: '0'
  },
  ImageStickGridItem: {
    width: '100%'
  },
  keyboard: {
    margin: '300px auto',
    '& span': {
      color: '#000'
    }
  },
  cardCategorySocialWhite: {
    marginTop: '10px',
    color: 'rgba(255, 255, 255, 0.8)',
    '& .fab,& .fas,& .far,& .fal,& .material-icons': {
      fontSize: '22px',
      position: 'relative',
      marginTop: '-4px',
      top: '2px',
      marginRight: '5px'
    },
    '& svg': {
      position: 'relative',
      top: '5px'
    }
  },
  InfoBarGridItem: {},
  cardDescription: {
    ...description,
    fontSize: '45px',
    fontFamily: 'Roboto, Helvetica, Arial, sans-serif',
    textAlign: 'center'
  },
  cardCategoryWhite: {
    marginTop: '10px',
    color: 'rgba(255, 255, 255, 0.7)'
  },
  cardDescriptionWhite: {
    color: 'rgba(255, 255, 255, 0.8)'
  },
  CountDownItem: {
    height: '120px',
    padding: '0'
  },
  RightCommonItem: {
    marginLeft: '0',
    marginRight: '0',
    marginTop: '10px'
  },

  // InfoTab: {
  //   height: '150px',
  // },
  InfoWorkContainer: {
    height: '100px'
  },
  InfoWorkMarginContainer: {
    height: '100px',
    padding: '0',
    margin: '0',
    marginLeft: '10px'
  },
  InfoWorkItem: {
    padding: '0'
  },
  InfoWorkMarginItem: {
    marginLeft: '10px',
    marginTop: '0px',
    height: '100%',
    borderRadius: '0',
    paddingLeft: '-15px',
    marginBottom: '0px'
  },
  TimeLine: {
    height: '500px'
  },
  LeftWrapper: {
    height: '100%'
    // padding: '20px 5px 0 20px!important'
  },
  LeftTopWrapper: {
    marginTop: '0'
  },
  MainWrapper: {
    height: '100%'
  },
  LeftBottomWrapper: {
    marginTop: '11px',
    height: 'calc(100% - 160px)'
  },
  LeftTop1: {
    // padding: '0 5px 0px 12px!important',
  },
  LeftTop2: {
    // padding: '0 5px 0 5px!important'
  },
  LeftTop3: {
    // padding: '0 12px 0 5px!important'
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
  RightWrapper: {
    height: '100%'
    // padding: '20px 20px 0 5px!important'
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
  MarginTop5: {
    marginTop: '5px'
  },
  MarginTopBottom5: {
    margin: '0 0 5px'
  },
  LeftPadding: {
    padding: '5px 5px 5px 12px!important'
  },
  RightPadding: {
    padding: '5px 12px 5px 5px!important'
  }
});

class ConnectedWorking extends React.Component {
  constructor(props) {
    super(props);

    this.autoCancel = null;
    this.keyboard = null;

    this.state = {
      manualDiagShow: false
    };

    this.toggleOPMode = this.toggleOPMode.bind(this);
  }

  shouldComponentUpdate() {
    return true;
  }

  componentDidUpdate(prevProps) {
    this.prevOperationStatus = prevProps.operations.operationStatus;
  }

  toggleOPMode() {
    const { isAutoMode } = this.props;
    // this.props.setAutoMode(!isAutoMode);
  }

  closeManualDiag = () => {
    this.setState({
      manualDiagShow: false
    });
  };

  openManualDiag = (e, input,t) => {
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
        value: toolSvg(tools.status === 'connected' ? successColor : dangerColor,{ width: '4vh', height: '4vh', }),
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

    const { manualDiagShow } = this.state;

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
    // const teststory = [
    //   {
    //     // First story
    //     inverted: true,
    //     badgeColor: "danger",
    //     badgeIcon: CardTravel,
    //     title: "Some Title",
    //     titleColor: "danger",
    //     body: (
    //       <p>
    //         Wifey made the best Father's Day meal ever. So thankful so happy so
    //         blessed. Thank you for making my family We just had fun with the
    //         “future” theme !!! It was a fun night all together ... The always rude
    //         Kanye Show at 2am Sold Out Famous viewing @ Figueroa and 12th in
    //         downtown.
    //       </p>
    //     ),
    //     footerTitle: "11 hours ago via Twitter"
    //   }
    // ];
    //
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
          <Grid container spacing={16} className={classes.root} justify="center">
            <Grid item xs={9} container spacing={8}>
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
                  onClick={e => this.openManualDiag(e, 'carType',t)}
                  disabled={
                    !manualEnable ||
                    configs.operationSettings.opMode === 'order'
                  }
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        车型:
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
                  onClick={e => this.openManualDiag(e, 'vin',t)}
                  disabled={!manualEnable}
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        VIN/KNR:
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
            <Grid item xs={3} container spacing={8} alignContent="flex-start" alignItems="flex-start">
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
            <ManualDiag show={manualDiagShow} close={this.closeManualDiag}/>
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

// export default Index;
export default withKeyboard(withStyles(withstyles)(Working));
