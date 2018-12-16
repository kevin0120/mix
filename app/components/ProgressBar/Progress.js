import React from 'react';
import Grid from '@material-ui/core/Grid/Grid';
import { withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import styles from './styles';
import Circle from './Circle';

// import { Line, Circle } from 'react-es6-progressbar.js';

class ProgressBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      timeRemaining: props.time,
      counterState: this.counterStates.ready
    };
    this.timer = null;
  }

  componentDidMount() {
    this.readyCounter();
  }

  componentWillReceiveProps(nextProps) {
    const {
      shouldCounterReady,
      shouldCounterStart,
      shouldCounterStop
    } = nextProps;
    const { counterState } = this.state;
    // console.log('will receive props');
    if (
      shouldCounterReady &&
      shouldCounterReady() &&
      counterState !== this.counterStates.ready &&
      counterState !== this.counterStates.start
    ) {
      this.readyCounter(nextProps);
    }
    if (
      shouldCounterStart &&
      shouldCounterStart() &&
      counterState !== this.counterStates.ticking
    ) {
      // console.log('should start');
      console.log(nextProps);
      this.startCounter(nextProps);
    }
    if (
      shouldCounterStop &&
      shouldCounterStop() &&
      counterState !== this.counterStates.stop
    ) {
      this.stopCounter(nextProps);
    }
  }

  componentWillUnmount() {
    if (this.timer !== null) {
      clearInterval(this.timer);
      this.timer = null;
    }
  }

  counterStates = {
    ready: 'ready',
    ticking: 'ticking',
    pause: 'pause',
    stop: 'stop',
    finish: 'finish'
  };

  readyCounter = () => {
    this.setState({
      counterState: this.counterStates.ready
    });
  };

  stopCounter = () => {
    console.log('counter stop');
    this.setState(
      {
        counterState: this.counterStates.stop,
        timeRemaining: 0
      },
      () => {
        clearInterval(this.timer);
        this.timer = null;
      }
    );
  };

  finishCounter = () => {
    const { onFinish } = this.props;
    console.log('counter finish');
    this.setState(
      {
        counterState: this.counterStates.finish
      },
      () => {
        clearInterval(this.timer);
        this.timer = null;
        onFinish();
      }
    );
  };

  startCounter = props => {
    console.log('counter start');

    this.setState(
      {
        counterState: this.counterStates.ticking,
        timeRemaining: props.time
      },
      () => {
        console.log('set state:', this.state);
        this.timer = setInterval(() => {
          const { timeRemaining } = this.state;
          console.log(timeRemaining);
          if (timeRemaining <= 0) {
            this.finishCounter();
          } else {
            this.setState({
              timeRemaining: timeRemaining - 1
            });
          }
        }, 1000);
      }
    );
  };

  render() {
    const { time, gridClassName } = this.props;
    const { timeRemaining } = this.state;
    const circleProps = {
      text: `${timeRemaining}`,
      responsive: true,
      animate: true,
      animationDuration: 1,
      size: 150,
      lineWidth: 50,
      progress: time > 0 ? (time - timeRemaining) / time : 0,
      bgColor: '#FFFFFF',
      startColor: '#F7FFA2',
      endColor: '#F7E600',
      textColor: '#F7E600',
      textStyle: {
        font: '30rem Helvetica, Arial, sans-serif'
      },
      roundedStroke: true
    };
    return (
      <Grid item className={gridClassName}>
        <Circle
          {...circleProps}
          // text={timeRemaining}
          // responsive
          // animate
          // animationDuration="1"
          // size={150}
          // lineWidth={50}
          // progress={time > 0 ? (time - timeRemaining) / time : 0}
          // bgColor="#FFFFFF"
          // startColor="#F7FFA2"
          // endColor="#F7E600"
          // textColor="#F7E600"
          // textStyle={{
          //   font: '20rem Helvetica, Arial, sans-serif'
          // }}
          // roundedStroke
        />
      </Grid>
    );
  }
}

ProgressBar.propTypes = {
  time: PropTypes.number.isRequired,
  gridClassName: PropTypes.string.isRequired,
  shouldCounterReady: PropTypes.func,
  shouldCounterStart: PropTypes.func,
  shouldCounterStop: PropTypes.func,
  onFinish: PropTypes.func
};

ProgressBar.defaultProps = {
  shouldCounterReady: () => {},
  shouldCounterStart: () => {},
  shouldCounterStop: () => {},
  onFinish: () => {}
};

export default withStyles(styles)(ProgressBar);
