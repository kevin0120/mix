import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import styles from './style';
import { viewingStep, stepPayload, stepData } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';

const mapState = (state, props) => ({
  ...props,
  points: stepData(viewingStep(state.order))?.points || [],
  image: stepPayload(viewingStep(state.order))?.image || {},
  activePoint: stepData(viewingStep(state.order))?.activePoint
});

const mapDispatch = {};

function ScrewStep({ image, points, activePoint }) {
  const classes = makeStyles(styles)();
console.log(activePoint);
  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      activePoint={activePoint}
      focus={2}
    />
  </div>;
}


export default connect(
  mapState,
  mapDispatch
)(ScrewStep);

