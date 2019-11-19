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
import type { tAnyStatus } from '../../../modules/step/interface/typeDef';
import type { Dispatch, tAction } from '../../../modules/typeDef';
import { ClsOperationPoint } from '../../../modules/step/screwStep/classes/ClsOperationPoint';
import { translation as trans } from '../../../components/NavBar/local';
import notifierActions from '../../../modules/Notifier/action';
import type { CommonLogLvl } from '../../../common/utils';


type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  points: Array<ClsOperationPoint>,
  image: string,
  workCenterMode: string,
  status: ?tAnyStatus
|};

type tDP = {|
  redoPoint: Dispatch,
  result: Dispatch,
  notify: (variant: CommonLogLvl, message: string, meta: Object)=>tAction<any, any>
|};

type Props = {|
  ...tOP,
  ...tSP,
  ...tDP
|};

const mapState = (state, props: tOP): tSP => {
  const vStep = viewingStep(state.order);
  return {
    ...props,
    points: stepData(vStep)?.tightening_points || stepPayload(vStep)?.tightening_points || [],
    image: vStep?.image || '',
    workCenterMode: state.workCenterMode,
    status: stepStatus(vStep)
  };
};

const mapDispatch: tDP = {
  result: screwStepAction.result,
  redoPoint: screwStepAction.redoPoint,
  notify: notifierActions.enqueueSnackbar
};


function ScrewStep({ isCurrent, image, points, workCenterMode, notify, redoPoint, result }: Props) {
  const classes = makeStyles(styles)();
  return <div className={classes.layout}>
    <ScrewImage
      image={image}
      points={points}
      // activeIndex={isCurrent ? activeIndex : -1}
      focus={0}
      scale={1}
      twinkle={isCurrent}
      enableReWork={workCenterMode === trans.reworkWorkCenterMode}
      notifyInfo={notify}
      onPointClick={(point: ClsOperationPoint) => {
        // result({
        //   data: [{
        //     ...point._point,
        //     result: RESULT_STATUS.ok
        //   }]
        // });
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

