// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import styles from './style';
import { viewingStep, stepPayload, stepData, stepStatus } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';
import STEP_STATUS from '../../../modules/step/constants';
import type { tStepProps } from '../types';

const mapState = (state, props) => {
  const vStep = viewingStep(state.order);

  return ({
    ...props,
    points: stepData(vStep)?.points || stepPayload(vStep)?.points || [],
    image: stepPayload(vStep)?.image || '',
    activeIndex: stepData(vStep)?.activeIndex,
    status: stepStatus(vStep)
  });
};

const mapDispatch = {
  result: screwStepAction.result,
  redoPoint: screwStepAction.redoPoint
};

type Props = {
};


function ScrewStep({ isCurrent, status, image, points, activeIndex, result, redoPoint }: Props & tStepProps) {
  const classes = makeStyles(styles)();

  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      activeIndex={isCurrent ? activeIndex : -1}
      focus={status === STEP_STATUS.DOING ? 2 : 0}
      scale={1}
      twinkle={isCurrent}
      // onClick={() => result({ data: [{ result: 'ok' }] })}
      onPointClick={(point) => {
        // console.log('on point click', point);

        // redoPoint(point);
      }}
    />
    <Paper
      square
      style={{ position: 'absolute', bottom: 10, right: 10, width: 'auto' }}
    >
      <ScrewImage
        style={{ width: '200px', height: '200px' }}
        image={image}
        points={points}
        twinkle={isCurrent}
        activeIndex={isCurrent ? activeIndex : -1}
        focus={0}
        pointScale={0.5}
      />
    </Paper>
  </div>;
}


export default connect(
  mapState,
  mapDispatch
)(ScrewStep);

