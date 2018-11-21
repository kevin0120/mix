import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import Grid from '@material-ui/core/Grid';
import Paper from '@material-ui/core/Paper';
import { withStyles } from '@material-ui/core/styles';
import { I18n } from 'react-i18next';

import { Line } from 'react-es6-progressbar.js'
// import { WorkingStyle } from "../../actions";
// import { WorkMode} from "../../actions/commonActions";
// import {
//   progressCountingStarted,
//   progressCountingStopped,
// } from '../../actions/progressCounting';

import { Warn } from '../../logger';

import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  // jobStatus: state.orderProgress.jobStatus,
  orderStatus: state.orderProgress.orderStatus,
  hintMsg: state.orderProgress.hintMsg,
  // isProcessing: state.orderProgress.isProcessing,
  // shoudStartCounting: state.orderProgress.shoudStartCounting,
  maxOpTime: state.ongoingWorkOrder.maxOpTime,
  carID: state.orderProgress.carID,
  carType: state.orderProgress.carType,
  isAutoMode: state.isAutoMode,
  ...ownProps,
});

const mapDispatchToProps = {
  // progressCountingStarted,
  // progressCountingStopped,
};


class ConnectedProgressBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      // completed: 100,
      startAnimation: 1,
    };
    // this.timer = null;
    this.initCount = this.initCount.bind(this);
    this.startCount = this.startCount.bind(this);
    this.stopCount = this.stopCount.bind(this);
    this.finishCount = this.finishCount.bind(this);
    this.setCompleted = this.setCompleted.bind(this);
  }

  shouldComponentUpdate(nextProps, nextState) {

    let rt = true;

    if (this.state.startAnimation !== nextState.startAnimation ) {
      return true;
    }

    if (this.props.orderStatus === nextProps.orderStatus) {
      rt = false;
    }

    if (this.props.orderStatus === 'Done' && (nextProps.orderStatus === 'Timeout' || nextProps.orderStatus === 'Fail')) {
      rt = false;
    }

    if (this.props.orderStatus === 'Timeout' && (nextProps.orderStatus === 'Fail')) {
      rt = false;
    }

    if (this.props.orderStatus === 'Fail' && (nextProps.orderStatus === 'Timeout')) {
      rt = false;
    }

    return rt;
  }


  componentWillUpdate(nextProps,  nextState) {
    const {orderStatus, isAutoMode, progressCountingStarted, carType, carID} = this.props;

    const nextorderStatus = nextProps.orderStatus;

    switch(nextorderStatus) {
      case 'Doing':
        // if (userConfigs.workFlow === 'General' && !isAutoMode) {
        //   break;
        // }

        this.startCount();
        // progressCountingStarted();
        break;
      case 'Ready':
        this.finishCount();
        break;

      case 'Timeout':
        this.stopCount();
        break;

      case 'Done':
        console.log("progress done");
        // this.finishCount();
        break;

      default:
        this.initCount();
    }
  }

  componentDidUpdate(preProps, preState) {

    // const {orderStatus, isAutoMode, progressCountingStarted} = this.props;
    //
    // switch(orderStatus) {
    //   case 'Doing':
    //     // if (userConfigs.workFlow === 'General' && !isAutoMode) {
    //     //   break;
    //     // }
    //
    //     this.startCount();
    //     // progressCountingStarted();
    //     break;
    //
    //   case 'Timeout':
    //     this.stopCount();
    //     break;
    //
    //   case 'Done':
    //     console.log("progress done");
    //     this.finishCount();
    //     break;
    //
    //   default:
    //     this.initCount();
    // }
  }

  componentWillUnmount() {
    // clearInterval(this.timer);
  }

  setCompleted() {
    const {carID} = this.props;
    Warn("工单已超时 车辆ID:" + carID);
    this.props.progressCountingStopped();
  }

  initCount() {
    this.setState( {
      startAnimation: -1
    })
  }

  startCount() {
    this.setState( {
      startAnimation: 0
    })
  }

  stopCount() {
    this.setState( {
      startAnimation: 1
    })
  }

  finishCount() {
    this.setState( {
      startAnimation: 2
    })
  }


  render() {

    const {
      classes, carID, hintMsg, orderStatus
    } = this.props;

    const { startAnimation } = this.state;
    console.log("progress startAnimation:", startAnimation);

    // let progressHint = WorkingStyle[orderStatus].hint; // default
    // if (hintMsg && hintMsg.length > 0) {
    //   progressHint = hintMsg;
    // } else if (!progressHint) {
    //   progressHint = carID;
    // }

    const LineOptions = {
      strokeWidth: 4,
      easing: 'linear',
      duration: this.props.maxOpTime * 1000,
      color: '#1ca552',
      trailColor: '#eee',
      trailWidth: '100%',
      svgStyle: {width: '100%', height: '100%'},
      from: {color: '#1ca552'},
      to: {color: '#F44336'},
      step: (state, bar) => {
        bar.path.setAttribute('stroke', state.color);
      }
    };

    return (
      <I18n ns="translations">
        {
          t => (
            <Grid item className={classes.progressWrap}>
              <Line
                progress={1.0}
                options={LineOptions}
                container_class={classes.progressMote}
                container_style={{height: '100%', width: '100%'}}
                startAnimate={startAnimation}
                onStop={() => this.setCompleted()}
              />
              <Paper className={classes.progressText}>
                {t('')}
              </Paper>
            </Grid>
          )
        }
      </I18n>
    );
  }
}


ConnectedProgressBar.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  maxOpTime: PropTypes.number,
  carID: PropTypes.string,
  carType: PropTypes.string,
  progressCountingStarted: PropTypes.func.isRequired,
  progressCountingStopped: PropTypes.func.isRequired,
  // jobStatus: PropTypes.shape({
  //   key: PropTypes.string.isRequired,
  //   hint: PropTypes.string,
  // }).isRequired,
  orderStatus: PropTypes.string.isRequired,
  hintMsg: PropTypes.string,
  isAutoMode: PropTypes.bool.isRequired,
};

ConnectedProgressBar.defaultProps = {
  maxOpTime: 30,
  carID: '',
  carType: '',
  hintMsg: '',
};

const WorkProgressBar = connect(mapStateToProps, mapDispatchToProps)(ConnectedProgressBar);

export default withStyles(styles)(WorkProgressBar);
