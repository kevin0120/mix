import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import styles from './style';
import { viewingStep, stepPayload, stepData } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';

const mapState = (state, props) => ({
  ...props,
  points: stepData(viewingStep(state.order))?.points || [],
  image: stepPayload(viewingStep(state.order))?.image || {},
  activeIndex: stepData(viewingStep(state.order))?.activeIndex,
});

const mapDispatch = {
  result: screwStepAction.result,
  imageReady: screwStepAction.imageReady
};

function ScrewStep({ image, points, activeIndex, result, imageReady }) {
  const classes = makeStyles(styles)();
  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      activeIndex={activeIndex}
      focus={2}
      scale={1}
      onClick={() => result([{status:'success'}])}
      onReady={imageReady}
      pointScale={0.5}

    />
    <ScrewImage
      style={{ width: '200px', height: '200px', position: 'absolute', bottom: 0, right: 0 }}
      image={image}
      points={points}
      activeIndex={activeIndex}
      focus={0}
      onReady={imageReady}
      pointScale={0.5}
    />
  </div>;
}


export default connect(
  mapState,
  mapDispatch
)(ScrewStep);

