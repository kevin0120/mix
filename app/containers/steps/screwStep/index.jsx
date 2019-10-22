// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import styles from './style';
import { viewingStep, stepPayload, stepData, stepStatus } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';
import type { tStepProps } from '../types';
import type { tAnyStepStatus } from '../../../modules/step/interface/typeDef';
import type { Dispatch } from '../../../modules/typeDef';
import { ClsOperationPoint } from '../../../modules/step/screwStep/classes/ClsOperationPoint';
import { RESULT_STATUS } from '../../../modules/step/screwStep/constants';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  points: Array<ClsOperationPoint>,
  image: string,
  // activeIndex: number,
  status: ?tAnyStepStatus
|};

type tDP = {|
  redoPoint: Dispatch,
  result: Dispatch
|};

type Props = {|
  ...tOP,
  ...tSP,
  ...tDP
|};

const mapState = (state, props: tOP): tSP => {
  const vStep = viewingStep(state.order);

  return ({
    ...props,
    points: stepData(vStep)?.points || stepPayload(vStep)?.points || [],
    image: stepPayload(vStep)?.image || '',
    // activeIndex: stepData(vStep)?.activeIndex,
    status: stepStatus(vStep)
  });
};

const mapDispatch: tDP = {
  result: screwStepAction.result,
  redoPoint: screwStepAction.redoPoint
};


function ScrewStep({ isCurrent, image, points, redoPoint, result }: Props) {
  const classes = makeStyles(styles)();

  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      // activeIndex={isCurrent ? activeIndex : -1}
      focus={0}
      scale={1}
      twinkle={isCurrent}
      onPointClick={(point) => {
        result({
          data: [{
            ...point._point,
            result: RESULT_STATUS.ok
          }]
        });
        // redoPoint(point);
      }}
    />
    <Paper
      square
      className={classes.thumbPaper}
    >
      <ScrewImage
        style={{
          width: '200px',
          height: '200px'
        }}
        image={image}
        points={points}
        twinkle={isCurrent}
        focus={0}
        pointScale={0.5}
      />
    </Paper>
  </div>;
}


export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(ScrewStep);

