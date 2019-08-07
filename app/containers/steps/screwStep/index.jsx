import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import styles from './style';
import { viewingStep, stepPayload, stepData, stepStatus } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';
import STEP_STATUS from '../../../modules/step/model';
// import { staticScrewTool } from '../../../modules/tools/saga';

const mapState = (state, props) => ({
  ...props,
  points: stepData(viewingStep(state.order))?.points || [],
  image: stepPayload(viewingStep(state.order))?.image || {},
  activeIndex: stepData(viewingStep(state.order))?.activeIndex,
  status:stepStatus(viewingStep(state.order))
});

const mapDispatch = {
  result: screwStepAction.result,
  imageReady: screwStepAction.imageReady
};

function ScrewStep({ status,image, points, activeIndex, result }) {
  const classes = makeStyles(styles)();

  // useEffect(() => {
  //   staticScrewTool.Enable();
  //   return () => staticScrewTool.Disable();
  //   // eslint-disable-next-line react-hooks/exhaustive-deps
  // }, []);

  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      activeIndex={activeIndex}
      focus={status===STEP_STATUS.DOING?2:0}
      scale={1}
      onClick={() => result([{status:'success'}])}
    />
      <Paper
        square
        style={{ position: 'absolute', bottom: 10, right: 10 ,width:'auto'}}
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

