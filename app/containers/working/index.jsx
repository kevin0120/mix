import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import Button from '../../components/CustomButtons/Button.jsx';
import classNames from 'classnames';

import List from '@material-ui/core/List';
import Avatar from '@material-ui/core/Avatar';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import NoteAddOutlined from '@material-ui/icons/NoteAddOutlined';
import PlayArrowOutlined from '@material-ui/icons/PlayArrowOutlined';
import FiberManualRecord from '@material-ui/icons/FiberManualRecord';
import Autorenew from '@material-ui/icons/Autorenew.js';

import WorkProgressBar from '../../components/ProgressBar/ProgressBar';
// import WorkStatusLed from '../../components/StatusLED/StatusLED';
import ImageStick from '../../components/ImageStick/imageStick';
// import OeeDiag from '../../components/OeeDiag';

// import ManualBar from '../../components/ManualBar';
// import SwitchModeBar from '../../components/SwitchModeBar';
import ShutdownDiag from '../../components/ShutDownDiag';

import Clock from 'react-live-clock';

import CardTravel from '@material-ui/icons/CardTravel';

import { I18n } from 'react-i18next';

// actions
// import { IOSet } from "../../actions/controllerIO";

// import {
//   openOeeDiag,
//   setAutoMode,
//   setCarByPass,
// } from '../../actions/commonActions';

// import { setOrderStatus } from "../../actions/ongoingWorkOrder";

//components
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';

import Divider from '@material-ui/core/Divider';
import withLayout from '../../components/Layout/layout';

import WorkingInfoBar from '../../components/WorkingInfoBar';

import HealthCheck from '../../components/HealthCheck';

import Keyboard from 'react-simple-keyboard';

const lodash = require('lodash');

import {
  container,
  cardTitle,
  description
} from '../../common/jss/material-react-pro';
import ResultDialog from '../../components/ResultDialog';
import ManualDiag from '../../components/ManualDiag';

import TimeLine from '../../components/WorkPageTimeline';
import Paper from '@material-ui/core/Paper';

import ProgressBar from '../../components/ProgressBar/Progress';
import { OPERATION_STATUS } from '../../reducers/operations';
// import { progressCountingStarted, progressCountingStopped } from '../../actions/progressCounting';

const mapStateToProps = (state, ownProps) => ({
  operations: state.operations,
  operationSettings: state.setting.operationSettings,
  timeline: state.timeline,
  // odooUrl: state.userConfigs.odooConnection.odooUrl.value,
  // hmiSn: state.userConfigs.odooConnection.hmiSn.value,
  // masterpcUrl: state.connInfo.masterpc.connection,
  // controllerSN: state.workingPage.infoTools.controllerSN.value,
  // workcenterId: state.workingPage.infoUser.workcenter.id,
  // activeResultIdIdx: state.orderProgress.activeResultIdIdx,
  // failCnt: state.orderProgress.failCnt,
  // results: state.ongoingWorkOrder.results,
  // userid: state.userInfo.id,
  // name: state.userInfo.name,
  // secondaryInfo: state.userInfo.login, // email or role etc
  // avatarImg: state.userInfo.image_small,
  // infos: state.workingPage, // stories在其中
  // orderStatus: state.orderProgress.orderStatus,
  // showSwitchMode: state.userConfigs.showSwitchMode,
  // switchAutoManual: state.userConfigs.switchAutoManual,
  // oeeFuncEnable: state.userConfigs.oeeFuncEnable,
  // isAutoMode: state.isAutoMode,
  // workFlow: state.userConfigs.workFlow,
  // carID: state.orderProgress.carID,
  // carType: state.orderProgress.carType,
  // lnr: state.orderProgress.lnr,
  // isEmptyCar: state.orderProgress.isEmptyCar,
  // jobID: state.ongoingWorkOrder.jobID,
  // byPassJob: state.userConfigs.byPassJob,
  // bypassInfo: state.userConfigs.bypass,
  // ControlMode: state.ControlMode,
  // productID: state.ongoingWorkOrder.productID,
  // maxOpTime: state.ongoingWorkOrder.maxOpTime,
  ...ownProps
});

const TOPHEIGHT = '150px';

const mapDispatchToProps = {
  // setAutoMode,
  // openOeeDiag,
  // // job,
  // // pset,
  // setCarByPass,
  // // IOSet,
  // // setOrderStatus,
  // progressCountingStarted,
  // progressCountingStopped
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
    position: 'relative',
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
    position: 'fixed',
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
    height: '100%',
    padding: '20px 5px 0 20px!important'
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
    padding: '0 5px 0px 12px!important'
  },
  LeftTop2: {
    padding: '0 5px 0 5px!important'
  },
  LeftTop3: {
    padding: '0 12px 0 5px!important'
  },
  LeftTopTab: {
    textAlign: 'left',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    width: '100%',
    height: '100px'
  },
  LeftBottomTab: {
    textAlign: 'left',
    position: 'relative',
    background: '#FFFFFF',
    borderRadius: '0',
    height: '100%'
  },
  LeftTabContiner: {
    height: '100%',
    width: '100%',
    padding: '2px 10px',
    position: 'relative'
  },
  ImgTabContiner: {
    height: '100%',
    width: '100%',
    padding: '0px 0px',
    position: 'absolute'
  },
  LeftTopDes: {
    marginBottom: '0',
    color: '#979797',
    '& p': {
      fontSize: '14px'
    }
  },
  RightWrapper: {
    height: '100%',
    padding: '20px 20px 0 5px!important'
  },
  CutDownPaper: {
    textAlign: 'center',
    position: 'relative',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '15%',
    background: '#2A2A2A'
  },
  InfoTab: {
    textAlign: 'center',
    position: 'relative',
    background: '#FFFFFF',
    marginTop: '17px',
    borderRadius: '0'
  },
  InfoTabTimeLine: {
    textAlign: 'center',
    position: 'relative',
    background: '#FFFFFF',
    marginTop: '10px',
    borderRadius: '0',
    height: 'calc(100% - 505px)'
  },
  InfoTabContiner: {
    height: '100%',
    width: '100%'
  },
  CutDownContainer: {
    position: 'absolute',
    height: '75px',
    marginTop: '0px',
    width: '100%',
    color: '#E5F0FA'
  },
  CountDownContainer: {
    position: 'absolute',
    height: '90%',
    //marginTop: '10px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    width: '100%',
    color: '#E5F0FA'
    //top: '15px',
    //bottom: '15px',
  },
  TurnPaper: {
    textAlign: 'center',
    height: '120px',
    position: 'relative',
    background: '#66AECE'
  },
  RetryPaper: {
    textAlign: 'center',
    height: '120px',
    position: 'relative',
    background: '#6B70BD'
  },
  RightDescription: {
    fontSize: '20px',
    marginBottom: '0'
  },
  RightNum: {
    fontSize: '60px',
    fontWeight: '600'
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

function sortInfoToDisplay(infos, t) {
  return Object.keys(infos)
    .sort((a, b) => infos[a].displayOrder - infos[b].displayOrder)
    .map(name => {
      const { value, displayTitle } = infos[name];
      const d = t(displayTitle);
      return {
        key: name,
        value,
        displayTitle: d
      };
    });
}

class ConnectedWorking extends React.Component {
  constructor(props) {
    super(props);

    this.keyboard = null;

    this.state = {
      inputName: '',
      layoutName: 'default',
      vehicle: '',
      vin: '',
      manualDiagShow: false,
      resultShow: false
    };

    this.toggleOPMode = this.toggleOPMode.bind(this);
  }

  componentDidMount() {
    // this.props.setCarByPass(false);
  }

  shouldComponentUpdate() {
    return true;
  }

  componentDidUpdate(prevProps) {
    const { operations, operationSettings } = this.props;
    this.switchResultDiagShow(
      prevProps.operations.operationStatus,
      operations.operationStatus,
      operationSettings.enableResultDialog
    );

    // const {
    //   masterpcUrl,
    //   orderStatus, userid, workcenterId,
    //   results, activeResultIdIdx, failCnt, pset, workFlow, controllerSN, carType, carID, job, jobID, hmiSn, ControlMode, isAutoMode, IOSet, productID, isEmptyCar,byPassJob, bypassInfo
    // } = this.props;
    //
    // const prevOrderstatus = prevProps.orderStatus;
    //
    //
    // const resultInfo = results[activeResultIdIdx];
    //
    // const manualJob = !isAutoMode;
    // let isLast = false;
    // console.log('isEmptyCar',isEmptyCar);
    // if (isEmptyCar !== prevProps.isEmptyCar && isEmptyCar) {
    //   //发送空车信息
    //   job(masterpcUrl, controllerSN, carType, carID, userid, byPassJob, results, hmiSn, productID, workcenterId, true, manualJob);
    //   return;
    // }
    // if (lodash.includes(['Doing' ,'Continue','Done'], orderStatus) && results.length) {
    //   if (activeResultIdIdx + 1 >= results.length) {
    //     isLast = true;
    //   }
    //
    //   if (workFlow === 'VW') {
    //     pset(masterpcUrl, resultInfo, failCnt, userid, isLast);
    //   } else {
    //     if (orderStatus === 'Doing') {
    //
    //       if (ControlMode === 'job') {
    //         job(masterpcUrl, controllerSN, carType, carID, userid, jobID, results, hmiSn, productID, workcenterId, false, manualJob);
    //       }
    //     }
    //   }
    // }
    // if (orderStatus === 'PreDoing' && bypassInfo.enable) {
    //   IOSet(masterpcUrl, controllerSN, bypassInfo.output, 'off')
    // }
  }

  componentWillUnmount() {
    // this.props.setCarByPass(true);
  }

  handleClickOpenOee = e => {
    // this.props.openOeeDiag(true);
  };

  handleClickResume = e => {
    // this.props.setOrderStatus('Doing');
  };

  switchResultDiagShow = (prevOperationStatus, operationStatus, enable) => {
    if (
      operationStatus === OPERATION_STATUS.READY &&
      prevOperationStatus === OPERATION_STATUS.DOING &&
      enable
    ) {
      this.setState({
        resultShow: true
      });
    } else if (
      operationStatus === OPERATION_STATUS.DOING &&
      prevOperationStatus === OPERATION_STATUS.PREDOING
    ) {
      this.setState({
        resultShow: false
      });
    }
  };

  toggleOPMode() {
    const { isAutoMode } = this.props;
    // this.props.setAutoMode(!isAutoMode);
  }

  handleChangeInput = input => {
    const { inputName } = this.state;
    if (inputName === 'vehicle') {
      this.setState({
        vehicle: input
      });
    }
    if (inputName === 'vin') {
      this.setState({
        vin: input
      });
    }
  };

  handleChangeInputValue = event => {
    const { inputName } = this.state;
    if (inputName === 'vehicle') {
      this.setState({
        vehicle: event.target.value
      });
    }
    if (inputName === 'vin') {
      this.setState({
        vin: event.target.value
      });
    }
    this.keyboard.setInput(event.target.value, inputName);
  };

  handlePress = press => {
    const { vehicle, vin } = this.state;
    if (
      lodash.isEqual(press, '{enter}') &&
      (vehicle.length !== 0 || vin.length !== 0)
    ) {
      this.setState({
        inputName: ''
      });
    }
    if (lodash.isEqual(press, '{shift}')) {
      const layoutName = this.state.layoutName;

      this.setState({
        layoutName: layoutName === 'default' ? 'shift' : 'default'
      });
    }
  };

  handleInputName = ref => {
    this.setState({
      inputName: ref
    });
  };

  closeManualDiag = () => {
    this.setState({
      manualDiagShow: false
    });
  };

  openManualDiag = e => {
    e.preventDefault();
    this.setState({
      manualDiagShow: true
    });
  };

  orderInfo = t => {
    const { operations, operationSettings } = this.props;

    // const showButtonInfo = isAutoMode? 'Common.Auto':'Common.Manual';
    let programme = operations.jobID.toString();
    if (operationSettings.controllerMode === 'pset') {
      if (operations.results.length > 0) {
        programme = operations.results[
          operations.activeResultIndex
        ].pset.toString();
      }
    }

    return [
      {
        key: '工作模式',
        value: t(operationSettings.workMode),
        displayTitle: '工作模式'
      },
      {
        key: '工作状态',
        value: operations.operationStatus,
        displayTitle: '工作状态'
      },
      // {
      //   key: '控制模式',
      //   value: ControlMode,
      //   displayTitle: '控制模式',
      // },
      {
        key: '程序号',
        value: programme,
        displayTitle: '程序号'
      }
    ];
  };

  render() {
    const { classes, operations, timeline } = this.props;
    //
    //
    const { inputName, resultShow } = this.state;
    //
    let batch = '0/0';
    let redoBatch = '0/0';
    let maxRedoTimes = 0;
    if (operations.results.length > 0) {
      maxRedoTimes =
        operations.results[operations.activeResultIndex].max_redo_times;
      batch =
        (operations.activeResultIndex + 1).toString() +
        '/' +
        operations.results.length.toString();
      redoBatch =
        (maxRedoTimes - operations.failCount).toString() +
        '/' +
        maxRedoTimes.toString();
    }
    //
    // const fabClassName = classNames(classes.fab);
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
    // const showButtonInfo = isAutoMode? 'Common.Manual':'Common.Auto';
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
          <Grid container spacing={24} className={classes.root}>
            <Grid item xs={9} className={classes.LeftWrapper}>
              <Grid container spacing={24} className={classes.LeftTopWrapper}>
                <Grid item xs={2} className={classes.LeftTop1}>
                  <Paper className={classes.LeftTopTab}>
                    <div className={classes.LeftTabContiner}>
                      <h4 className={classes.LeftTopDes}>
                        <p
                          href="#pablo"
                          className={classes.MarginTopBottom5}
                          onClick={e => e.preventDefault()}
                        >
                          车序:
                        </p>
                      </h4>
                      <p className={classes.cardDescription}>
                        {operations.lnr}
                      </p>
                    </div>
                  </Paper>
                </Grid>
                <Grid item xs={3} className={classes.LeftTop2}>
                  <Paper className={classes.LeftTopTab}>
                    <div className={classes.LeftTabContiner}>
                      <h4 className={classes.LeftTopDes}>
                        <p
                          href="#pablo"
                          className={classes.MarginTopBottom5}
                          onClick={e => e.preventDefault()}
                        >
                          车型:
                        </p>
                      </h4>
                      <p className={classes.cardDescription}>
                        {operations.carType}
                      </p>
                    </div>
                  </Paper>
                </Grid>
                <Grid item xs={7} className={classes.LeftTop3}>
                  <Paper
                    className={classes.LeftTopTab}
                    component="button"
                    onClick={e => this.openManualDiag(e)}
                    disabled={false}
                  >
                    <div className={classes.LeftTabContiner}>
                      <h4 className={classes.LeftTopDes}>
                        <p
                          href="#pablo"
                          className={classes.MarginTopBottom5}
                          onClick={e => e.preventDefault()}
                        >
                          VIN/KNR:
                        </p>
                      </h4>
                      <p className={classes.cardDescription}>
                        {operations.carID}
                      </p>
                    </div>
                  </Paper>
                </Grid>
              </Grid>
              <Grid
                container
                spacing={24}
                className={classes.LeftBottomWrapper}
              >
                <Grid item xs={12} className={classes.MainWrapper}>
                  <Paper className={classes.LeftBottomTab}>
                    <div className={classes.ImgTabContiner}>
                      <ImageStick />
                    </div>
                  </Paper>
                </Grid>
              </Grid>
            </Grid>
            <Grid item xs={3} className={classes.RightWrapper}>
              <Paper className={classes.CutDownPaper}>
                <div className={classes.CountDownContainer}>
                  <ProgressBar
                    time={operations.maxOpTimes}
                    shouldCounterStart={() =>
                      lodash.includes(
                        [OPERATION_STATUS.DOING],
                        operations.operationStatus
                      )
                    }
                    shouldCounterStop={() =>
                      lodash.includes(
                        [OPERATION_STATUS.READY],
                        operations.operationStatus
                      )
                    }
                    // onStart={()=>this.props.progressCountingStarted()}
                    // onFinish={()=>this.props.progressCountingStopped()}
                    gridClassName={classes.progressWrap}
                  />
                </div>
              </Paper>
              <Grid container spacing={24} className={classes.MarginTop5}>
                <Grid item xs={6} className={classes.LeftPadding}>
                  <Paper className={classes.TurnPaper}>
                    <div className={classes.CutDownContainer}>
                      <h4>
                        <p
                          href="#pablo"
                          className={classes.RightDescription}
                          onClick={e => e.preventDefault()}
                        >
                          拧紧点
                        </p>
                      </h4>
                      <p className={classes.RightNum}>{batch}</p>
                    </div>
                  </Paper>
                </Grid>
                <Grid item xs={6} className={classes.RightPadding}>
                  <Paper className={classes.RetryPaper}>
                    <div className={classes.CutDownContainer}>
                      <h4>
                        <p
                          href="#pablo"
                          className={classes.RightDescription}
                          onClick={e => e.preventDefault()}
                        >
                          可拧紧次数
                        </p>
                      </h4>
                      <p className={classes.RightNum}>{redoBatch}</p>
                    </div>
                  </Paper>
                </Grid>
              </Grid>
              <Paper className={classes.InfoTab}>
                <div className={classes.InfoTabContiner}>
                  <List>
                    {/*<WorkingInfoBar*/}
                    {/*key="infoUser"*/}
                    {/*infos={sortInfoToDisplay(infos.infoUser, t)}*/}
                    {/*/>*/}
                    <Divider
                      className={classes.divider}
                      key="divider-infoUser"
                      light
                    />
                    <WorkingInfoBar key="infoOrder" infos={this.orderInfo(t)} />
                  </List>
                </div>
              </Paper>
              <Paper className={classes.InfoTabTimeLine}>
                <div className={classes.InfoTabContiner}>
                  <div className={classes.InfoTabContiner}>
                    <TimeLine simple stories={timeline.stories} />
                    {/*<TimeLine simple stories={teststory} />*/}
                  </div>
                </div>
              </Paper>
              <ShutdownDiag />
              <ResultDialog show={resultShow} />
              <ManualDiag
                show={this.state.manualDiagShow}
                close={this.closeManualDiag}
              />
            </Grid>
          </Grid>
        )}
      </I18n>
    );
  }
}

ConnectedWorking.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  // orderStatus: PropTypes.string.isRequired,
  // odooUrl: PropTypes.string.isRequired,
  // hmiSn: PropTypes.string.isRequired,
  // masterpcUrl: PropTypes.string.isRequired,
  // // workFlow: PropTypes.string.isRequired,
  // // shouldPset: PropTypes.bool.isRequired,
  // activeResultIdIdx: PropTypes.number.isRequired,
  // failCnt: PropTypes.number.isRequired,
  // userid: PropTypes.number.isRequired,
  // showSwitchMode: PropTypes.bool.isRequired,
  // results: PropTypes.arrayOf(PropTypes.shape({
  //   id: PropTypes.number,
  // }).isRequired).isRequired,
  // infos: PropTypes.shape({
  //   infoTools: PropTypes.shape({
  //     controllerSN: PropTypes.shape({}).isRequired,
  //   }).isRequired,
  //   workorder: PropTypes.shape({
  //     pset: PropTypes.number,
  //     result_ids: PropTypes.arrayOf(PropTypes.number),
  //   }),
  //   carID: PropTypes.string.isRequired,
  //   progress: PropTypes.shape({
  //     resultIdIdx: PropTypes.number.isRequired,
  //     status: PropTypes.string.isRequired,
  //   }).isRequired,
  // }).isRequired,
  // openOeeDiag: PropTypes.func.isRequired,
  // pset: PropTypes.func.isRequired,
  // job: PropTypes.func.isRequired,
  // IOSet: PropTypes.func.isRequired,
  // name: PropTypes.string.isRequired,
  // secondaryInfo: PropTypes.string.isRequired,
  // avatarImg: PropTypes.string.isRequired,
  // workFlow: PropTypes.string.isRequired,
  // controllerSN: PropTypes.string,
  // carType: PropTypes.string,
  // carID: PropTypes.string,
  // ControlMode: PropTypes.string,
  // productID: PropTypes.number,
  // lnr: PropTypes.string,
  // maxOpTime: PropTypes.number,
  operations: PropTypes.shape({}).isRequired,
  operationSettings: PropTypes.shape({}).isRequired,
  timeline: PropTypes.shape({}).isRequired
};

const Working = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWorking);

// export default Index;
export default withLayout(withStyles(withstyles)(Working));
