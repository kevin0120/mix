import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import styles from './style';
import { viewingStep, stepPayload, stepData, stepStatus } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';
import STEP_STATUS from '../../../modules/step/model';
import type { tStepProps } from '../types';
// import { staticScrewTool } from '../../../modules/tools/saga';

const mapState = (state, props) => {
  const vStep = viewingStep(state.order);
  console.log(vStep);
  return ({
    ...props,
    points: stepData(vStep)?.points || stepPayload(vStep).points || [],
    image: stepPayload(vStep)?.image || '',
    activeIndex: stepData(vStep)?.activeIndex,
    status: stepStatus(vStep)
  });
};

const mapDispatch = {
  result: screwStepAction.result,
  imageReady: screwStepAction.imageReady
};

function ScrewStep({ status, image, points, activeIndex, result }: tStepProps) {
  const classes = makeStyles(styles)();

  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      activeIndex={activeIndex}
      focus={status === STEP_STATUS.DOING ? 2 : 0}
      scale={1}
      onClick={() => result({ data: [{ result: 'ok' }] })}
    />
    <Paper
      square
      style={{ position: 'absolute', bottom: 10, right: 10, width: 'auto' }}
    >
      <ScrewImage
        style={{ width: '200px', height: '200px' }}

        image={image}
        points={points}
        activeIndex={activeIndex}
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

