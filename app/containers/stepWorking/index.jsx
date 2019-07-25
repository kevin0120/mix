import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core/styles';
import { Typography, Paper, createMuiTheme, MuiThemeProvider } from '@material-ui/core';
import { StylesProvider } from '@material-ui/styles';
import moment from 'moment';
import ButtonsContainer from './ButtonsContainer';
import * as orderSelectors from '../../modules/order/selector';
import styles from './styles';
import Dialog from '../../components/Dialog';
import StepperContainer from './StepperContainer';
import StepPageContainer from './StepPageContainer';

const format = t => `${t}`.length <= 2 ? `00${t}`.slice(-2) : t;

const renderTimer = (duration) => {
  const h = moment.duration(duration).hours();
  const m = moment.duration(duration).minutes();
  const s = moment.duration(duration).seconds();
  return <Typography variant="h4">
    {`${format(h)}:${format(m)}:${format(s)}`}
  </Typography>;
};

type StepWorkingProps = {
  currentOrder: {},
  startTime: Date
};

const theme = createMuiTheme({
  overrides: {
    MuiPaper: {
      root: {
        backgroundColor: 'inherit'
      }
    }
  }
});

function StepWorking({ currentOrder, startTime }: StepWorkingProps) {
  const classes = makeStyles(styles.layout)();
  const [action, bindAction] = useState(null);

  const [duration, setDuration] = useState(startTime ? new Date() - startTime : 0);


  useEffect(() => {
    const interval = setInterval(() => {
      const current = new Date();
      setDuration(startTime ? current - startTime : 0);
    }, 500);
    return () => clearInterval(interval);
  }, [setDuration, startTime]);

  return (
    <StylesProvider injectFirst>
      <MuiThemeProvider theme={theme}>
        <div className={classes.root}>
          <Dialog/>
          <Paper
            square
            className={classes.leftContainer}
            classes={{ root: classes.leftContainer }}
          >
            <Paper square className={classes.orderInfoContainer}>
              <Typography variant="h5">{currentOrder?.name}</Typography>
            </Paper>
            <ButtonsContainer action={action}/>
            <StepPageContainer bindAction={bindAction}/>
          </Paper>
          <div className={classes.rightContainer}>
            <Paper square className={classes.timerContainer}>
              {renderTimer(duration)}
            </Paper>
            <Paper square className={classes.stepperContainer}>
              <StepperContainer/>
            </Paper>
          </div>
        </div>
      </MuiThemeProvider>
    </StylesProvider>

  );
}

const mapState = (state, props) => ({
  ...props,
  currentOrder: orderSelectors.currentOrder(state.order) || {},
  startTime: orderSelectors.startTime(orderSelectors.processingStep(state.order)) || null
});

const mapDispatch = {};

export default connect(
  mapState,
  mapDispatch
)(StepWorking);
