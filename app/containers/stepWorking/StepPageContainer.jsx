// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles, Paper, Grid, Typography } from '@material-ui/core';
import { CardTravel } from '@material-ui/icons';
import stepTypes from '../steps/stepTypes';
import * as oSel from '../../modules/order/selector';
import styles from './styles';
import type { tStep } from '../../modules/order/model';
import TimeLine from '../../components/WorkPageTimeline';

const mapState = (state, props) => {
  const vStep = oSel.viewingStep(state.order);
  return {
    ...props,
    step: vStep,
    timeLine: oSel.stepData(vStep)?.timeLine || [],
    workingStep: oSel.workingStep(oSel.workingOrder(state.order)) || {}
  };
};

const mapDispatch = {};

type Props = {
  step: tStep,
  workingStep: tStep,
  bindAction: Function,
  timeLine: Array<any>
};

const StepPageContainer = ({
  step,
  workingStep,
  bindAction,
  timeLine
}: Props) => {
  const classes = makeStyles(styles.stepPageContainer)();
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;

    return (
      <Grid container spacing={1} className={classes.root}>
        <Grid item container className={classes.left} spacing={1}>
          <Grid item className={classes.left}>
            <Paper square className={classes.image}>
              {(StepComponent && (
                <StepComponent
                  step={step}
                  isCurrent={step === workingStep}
                  bindAction={bindAction}
                />
              )) ||
                null}
            </Paper>
          </Grid>
        </Grid>
        <Grid
          item
          container
          spacing={1}
          className={classes.right}
          direction="column"
        >
          <Grid item className={classes.description}>
            <Paper square className={classes.Paper}>
              <Typography>{step.description}</Typography>
            </Paper>
          </Grid>
          <Grid item className={classes.result}>
            <Paper square className={classes.Paper}>
              <TimeLine
                simple
                stories={timeLine.map(t => ({
                  simple: true,
                  inverted: true,
                  badgeColor: t.color,
                  titleColor: t.color,
                  title: t.title,
                  badgeIcon: CardTravel,
                  body: t.body,
                  footerTitle: t.footerTitle
                }))}
              />
            </Paper>
          </Grid>
        </Grid>
      </Grid>
    );
  }
  return null;
};

export default connect(
  mapState,
  mapDispatch
)(StepPageContainer);
