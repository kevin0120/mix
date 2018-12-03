import React from 'react';
import Grid from '@material-ui/core/Grid/Grid';
import { withStyles } from '@material-ui/core';
import styles from './styles';
import Circle from './Circle';
// import { Line, Circle } from 'react-es6-progressbar.js';

class ProgressBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      timeRemaining: props.time,
      counterState: this.counterStates.ready,
      startAnimation: -1
    };
    this.timer = null;
  }

  componentWillUnmount() {
    if (this.timer !== null) {
      clearInterval(this.timer);
      this.timer = null;
    }
  }

  componentDidMount() {
    this.readyCounter();
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
      counterState: this.counterStates.ready,
      // timeRemaining: this.props.time,
      startAnimation: -1
    });
  };

  stopCounter = () => {
    console.log('counter stop');
    this.setState(
      {
        counterState: this.counterStates.stop,
        timeRemaining: 0,
        startAnimation: -1
      },
      () => {
        clearInterval(this.timer);
        this.timer = null;
        this.props.onStop && this.props.onStop();
      }
    );
  };

  finishCounter = () => {
    console.log('counter finish');
    this.setState(
      {
        counterState: this.counterStates.finish,
        // timeRemaining: this.props.time,
        startAnimation: -1
      },
      () => {
        clearInterval(this.timer);
        this.timer = null;
        this.props.onFinish && this.props.onFinish();
      }
    );
  };

  pauseCounter = () => {
    console.log('counter pause');
    this.setState({ counterState: this.counterStates.pause }, () => {
      clearInterval(this.timer);
      this.timer = null;
      this.props.onPause && this.props.onPause();
    });
  };

  startCounter = () => {
    console.log('counter start');
    this.setState(
      {
        counterState: this.counterStates.ticking,
        startAnimation: 0,
        timeRemaining: this.props.time
      },
      () => {
        this.timer = setInterval(() => {
          this.setState(
            {
              timeRemaining: this.state.timeRemaining - 1
            },
            () => {
              console.log(this.state.timeRemaining);
              if (this.state.timeRemaining <= 0) {
                this.finishCounter();
              }
            }
          );
        }, 1000);
        this.props.onStart && this.props.onStart();
      }
    );
  };

  componentWillReceiveProps(nextProps) {
    const {
      shouldCounterReady,
      shouldCounterStart,
      shouldCounterPause,
      shouldCounterStop
    } = nextProps;
    console.log('will recieve props');
    if (
      shouldCounterReady &&
      shouldCounterReady() &&
      this.state.counterState !== this.counterStates.ready &&
      this.state.counterState !== this.counterStates.start
    ) {
      this.readyCounter();
    }
    if (
      shouldCounterStart &&
      shouldCounterStart() &&
      this.state.counterState !== this.counterStates.ticking
    ) {
      console.log('should start');
      this.startCounter();
    }
    if (
      shouldCounterPause &&
      shouldCounterPause() &&
      this.state.counterState === this.counterStates.ticking
    ) {
      this.pauseCounter();
    }
    if (
      shouldCounterStop &&
      shouldCounterStop() &&
      this.state.counterState !== this.counterStates.stop
    ) {
      this.stopCounter();
    }
  }

  render() {
    const { time, gridClassName, contentClassName } = this.props;
    const { timeRemaining } = this.state;
    const textDisplay = timeRemaining;
    // const LineOptions = {
    //   strokeWidth: 4,
    //   easing: 'linear',
    //   duration: time * 1000,
    //   color: '#005AB5',
    //   trailColor: '#f0f0ff',
    //   trailWidth: 4,
    //   svgStyle: { width: '90%', height: '90%' },
    //   from: { color: '#f0f0ff' },
    //   to: { color: '#7257B8' },
    //   step: (state, bar) => {
    //     bar.path.setAttribute('stroke', state.color);
    //   }
    // };
    return (
      <Grid item className={gridClassName}>
        <Circle
          text={textDisplay}
          responsive
          animate
          animationDuration="1"
          size={150}
          lineWidth={50}
          progress={time > 0 ? (time - timeRemaining) / time : 0}
          bgColor="#FFFFFF"
          startColor="#F7FFA2"
          endColor="#F7E600"
          textColor="#F7E600"
          textStyle={{
            font: '30rem Helvetica, Arial, sans-serif'
          }}
          roundedStroke
        />
      </Grid>
    );
  }
}

export default withStyles(styles)(ProgressBar);
