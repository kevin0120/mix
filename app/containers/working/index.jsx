import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';
import { withStyles } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { newData,manualCancel,manualTightening,start,close,didMount,clearStories} from '../../modules/manual/action';


import {
  container,
  cardTitle,
  description,
} from '../../common/jss/material-react-pro';
// import ResultDialog from '../../components/ResultDialog';

import TimeLine from '../../components/WorkPageTimeline';

import withKeyboard from '../../components/Keyboard';
import styles from './Counter.css';


const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  logo: state.logo,
  scanner: state.manual.scanner,
  tool: state.manual.tool,
  pset: state.manual.pset,
  result:state.manual.result,
  timeline: state.manual.timeline,
  ...ownProps
});

const TOPHEIGHT = '150px';

const mapDispatchToProps = {
  newData,
  manualCancel,
  manualTightening,
  start,
  close,
  didMount,
  clearStories,
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
    bottom: theme.spacing(2),
    right: theme.spacing(2),
    width: '150px'
},
fabOEE: {
  position: 'fixed',
    bottom: theme.spacing(15),
    right: theme.spacing(2),
    width: '150px'
},
fabResume: {
  position: 'fixed',
    bottom: theme.spacing(15),
    right: theme.spacing(2),
    width: '200px'
},
extendedIcon: {
  marginRight: theme.spacing(1)
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
  textAlign: 'center',
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


type Props = {
  start: () => void,
  close: () => void,
  manualTightening: () => void,
  manualCancel: () => void,
  newData: () => void,
  didMount: () => void,
  clearStories: () => void,
  scanner: string,
  tool: string,
  pset: number,
  result: Array<string>
};

class ConnectedWorking extends React.Component {

  props: Props;

  constructor(props) {
    super(props);
    this.keyboard = null;
  }

  componentDidMount(): void {
    this.props.start();
    this.props.didMount();
  }

  componentWillUnmount(): void {
    this.props.close();
    this.props.clearStories();
  }



  openManualDiag = (e, input, t) => {
    e.preventDefault();
    const { keyboardInput, newData:dispatchNewData} = this.props;
    keyboardInput({
      onSubmit: text => {
        dispatchNewData(text);
      },
      text: e.target.value,
      title: t(`Operation.Input.${input}`),
      label: t(`Operation.Input.${input}`)
    });
  };


  render() {
    const {
      classes,
      timeline,
      scanner,
      logo,
      tool,
      pset,
      result,
      manualCancel,
      manualTightening
    } = this.props;

    const batch = '0/0';
    const redoBatch = '0/0';

    return (
      <I18n ns="translations">
        {t => (
          <Grid container spacing={1} className={classes.root} justify="center">
            <Grid item xs={9} container spacing={1}>

              <Grid item xs={5} style={{ height: '13%' }}>
                <Paper
                  className={classes.LeftTopTab}
                  component="button"
                  onClick={e => this.openManualDiag(e, 'vin', t)}
                  disabled={true}
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        {t('Operation.Input.carType')}:
                      </p>
                    </h4>
                    <p className={classes.cardDescription}>
                      {}
                    </p>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={7} style={{ height: '13%' }}>
                <Paper
                  className={classes.LeftTopTab}
                  component="button"
                  onClick={e => this.openManualDiag(e, 'vin', t)}
                  disabled={false}
                >
                  <div className={classes.LeftTabContiner}>
                    <h4 className={classes.LeftTopDes}>
                      <p className={classes.MarginTopBottom5}>
                        {t('Operation.Input.vin')}:
                      </p>
                    </h4>
                    <p className={classes.cardDescription}>
                      {scanner}
                    </p>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '87%' }}>
                <Paper className={classes.LeftBottomTab}>
                  <img src={logo} alt="alt"/>
                </Paper>
              </Grid>
            </Grid>
            <Grid item xs={3} container spacing={1} alignContent="flex-start" alignItems="flex-start">
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
                      <p className={classes.RightDescription} onClick={manualCancel}>
                        选择工具
                      </p>
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
                      <p className={classes.RightDescription}　onClick={manualTightening}>
                        开始作业
                      </p>
                    </div>
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '10%' }}>
                <Paper className={classes.InfoTabTimeLine}>
                  <div className={`counter ${styles.counter1}`} data-tid="counter">
                    拧紧工具:{tool}
                  </div>
                </Paper>
              </Grid>
              <Grid item xs={12} style={{ height: '10%' }}>
                <Paper className={classes.InfoTabTimeLine}>
                  <div className={`counter ${styles.counter2}`} data-tid="counter">
                    程序号:{pset}
                  </div>
                </Paper>
              </Grid>

              <Grid item xs={12} style={{ height: '65%' }}>
                <Paper className={classes.InfoTabTimeLine}>
                  <TimeLine simple stories={timeline}/>
                  {/* <TimeLine simple stories={teststory} /> */}
                </Paper>
              </Grid>
            </Grid>

            {/*{<ResultDialog/>}*/}
          </Grid>
        )}
      </I18n>
    );
  }
}

ConnectedWorking.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  keyboardInput: PropTypes.func.isRequired,
};

const Working = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWorking);

// export default Index;
export default withKeyboard(withStyles(withstyles)(Working));
