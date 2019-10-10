// @flow
import React from 'react';
import type { Node } from 'react';
import { connect } from 'react-redux';
import { makeStyles, Paper, Grid, Typography } from '@material-ui/core';
import { InfoOutlined } from '@material-ui/icons';
import clsx from 'clsx';
import stepTypes from '../steps/stepTypes';
import * as oSel from '../../modules/order/selector';
import styles from './styles';
import TimeLine from '../../components/WorkPageTimeline';
import type { tClsStep } from '../../modules/step/Step';

const mapState = (state, props) => {
  const vStep: ?tClsStep = oSel.viewingStep(state.order);
  return {
    ...props,
    step: vStep,
    timeLine: oSel.stepData(vStep)?.timeLine || [],
    workingStep: oSel.workingStep(oSel.workingOrder(state.order)) || {}
  };
};

const mapDispatch = {};

type Props = {
  step: tClsStep,
  workingStep: tClsStep,
  bindParentAction: Node => any,
  timeLine: Array<any>,
  bindParentDescription: Node => any,
  description: Node
};

// 工步展示内容
const StepPageContainer = ({
  step,
  workingStep,
  bindParentAction,
  timeLine,
  description,
  bindParentDescription
}: Props) => {
  const classes = makeStyles(styles.stepPageContainer)();
  if (stepTypes?.[step?.type]?.component) {
    const StepComponent = stepTypes[step.type].component;
    return (
      <Grid container spacing={1} className={classes.root}>
        <Grid item container className={classes.leftContainer} spacing={1}>
          <Grid item className={classes.leftContent}>
            <Paper square className={classes.image}>
              {(StepComponent && (
                <StepComponent
                  step={step}
                  isCurrent={step === workingStep}
                  bindAction={bindParentAction}
                  bindDescription={bindParentDescription}
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
          <Grid item className={classes.descriptionContainer}>
            <Paper square className={clsx(classes.Paper, classes.Description)}>
              {description || (
                <Typography variant="h5">{step.description}</Typography>
              )}
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
                  badgeIcon: t.icon || InfoOutlined,
                  body: <Typography>{t.body}</Typography>,
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

export default connect<Props, *, _, _, _, _>(
  mapState,
  mapDispatch
)(StepPageContainer);
