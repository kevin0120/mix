// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles, Paper, Grid, Typography } from '@material-ui/core';
import stepTypes from '../steps/stepTypes';
import * as oSel from '../../modules/order/selector';
import styles from './styles';
import type { tStep } from '../../modules/order/model';
import TimeLine from '../../components/WorkPageTimeline';
import {CardTravel} from '@material-ui/icons';

const mapState = (state, props) => ({
  ...props,
  step: oSel.viewingStep(state.order) || {},
  workingStep: oSel.workingStep(oSel.workingOrder(state.order)) || {},
  result: oSel.stepData(oSel.viewingStep(state.order))?.result
});
const mapDispatch = {};

type Props = {
  step: tStep,
  workingStep: tStep,
  bindAction: ()=>{},
  result: Object
};

const StepPageContainer = ({ step, workingStep, bindAction, result }: Props) => {
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
              )) || null}
            </Paper>
          </Grid>
        </Grid>
        <Grid item container spacing={1} className={classes.right} direction="column">

          <Grid item className={classes.description}>
            <Paper square className={classes.Paper}>
              <Typography>
                {step.description}
              </Typography>
            </Paper>
          </Grid>
          <Grid item className={classes.result}>
            <Paper square className={classes.Paper}>

              {JSON.stringify(result)}
              <TimeLine
                simple
                stories={[
                {
                  // First story
                  simple:true,
                  inverted: true,
                  badgeColor: "danger",
                  badgeIcon: CardTravel,
                  title: "Some Title",
                  titleColor: "danger",
                  body: (
                    <p>
                      Wifey made the best Father's Day meal ever. So thankful so happy so
                      blessed. Thank you for making my family We just had fun with the
                      “future” theme !!! It was a fun night all together ... The always rude
                      Kanye Show at 2am Sold Out Famous viewing @ Figueroa and 12th in
                      downtown.
                    </p>
                  ),
                  footerTitle: "11 hours ago via Twitter"
                },
                {
                  // First story
                  simple:true,
                  inverted: true,
                  badgeColor: "danger",
                  badgeIcon: CardTravel,
                  title: "Some Title",
                  titleColor: "danger",
                  body: (
                    <p>
                      Wifey made the best Father's Day meal ever. So thankful so happy so
                      blessed. Thank you for making my family We just had fun with the
                      “future” theme !!! It was a fun night all together ... The always rude
                      Kanye Show at 2am Sold Out Famous viewing @ Figueroa and 12th in
                      downtown.
                    </p>
                  ),
                  footerTitle: "11 hours ago via Twitter"
                },
              ]}
              />
            </Paper>
          </Grid>
        </Grid>
      </Grid>
    );
  }
  return null;
};

export default connect(mapState, mapDispatch)(StepPageContainer);
