// @flow
import React, { useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import Paper from '@material-ui/core/Paper';
import styles from './style';
import { viewingStep, stepPayload, stepData, stepStatus } from '../../../modules/order/selector';
import ScrewImage from '../../../components/ScrewImage';
import screwStepAction from '../../../modules/step/screwStep/action';
import reworkActions from '../../../modules/reworkPattern/action';
import type { tStepProps } from '../types';
import type { tAnyStatus } from '../../../modules/step/interface/typeDef';
import type { Dispatch, tAction } from '../../../modules/typeDef';
import { ClsOperationPoint } from '../../../modules/step/screwStep/classes/ClsOperationPoint';
import notifierActions from '../../../modules/Notifier/action';
import type { CommonLogLvl } from '../../../common/utils';
import Alert from '../../../components/Alert';
import { withI18n } from '../../../i18n';
import { translation as Navtrans } from '../../../components/NavBar/local';
import {screwStepNS, translation as trans} from './local'


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
  redoPointSpecPoint: Dispatch,
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
  redoPointSpecPoint: reworkActions.aReworkSpecialScrewPoint,
  notify: notifierActions.enqueueSnackbar
};



function ScrewStep({ step, isCurrent, image, points, workCenterMode, notify, redoPointSpecPoint, result }: Props) {
  
  const screwImgClasses = makeStyles(styles)();
  
  const [showRedoScrewPointDiag, setShowRedoScrewPointDiag] = useState(false);
  
  const [redoScrewPoint, setRedoScrewPoint] = useState(null);
  
  const enableReWork = workCenterMode === Navtrans.reworkWorkCenterMode;
  
  return withI18n(
    t => (<div className={screwImgClasses.layout}>
      <ScrewImage
        image={image}
        points={points}
        // activeIndex={isCurrent ? activeIndex : -1}
        focus={0}
        scale={1}
        twinkle={isCurrent}
        enableReWork={enableReWork}
        notifyInfo={notify}
        onPointClick={(point: ClsOperationPoint) => {
          if (!point.canRedo) {
            notify('Error', '此拧紧点不具备返工条件!');
            return false;
          }
          setShowRedoScrewPointDiag(true);
          setRedoScrewPoint(point);
          return true
        }}
      />
      <Alert
        warning
        show={showRedoScrewPointDiag}
        title={t(trans.redoSpecScrewPointTitle)}
        onConfirm={() => {
          redoPointSpecPoint(step, redoScrewPoint);
          setShowRedoScrewPointDiag(false);
        }}
        onCancel={() => {
          setShowRedoScrewPointDiag(false);
        }}
        confirmBtnCssClass={`${screwImgClasses.button} ${
          screwImgClasses.success
          } ${screwImgClasses.buttonTxt}`}
        cancelBtnCssClass={`${screwImgClasses.button} ${
          screwImgClasses.danger
          } ${screwImgClasses.buttonTxt}`}
        confirmBtnText={t(trans.confirm)}
        cancelBtnText={t(trans.cancel)}
        showCancel
      >
        {`${t(trans.redoSpecScrewPointContent)} ${redoScrewPoint?.toString() || ""}`}
      </Alert>
      <Paper
        square
        className={screwImgClasses.thumbPaper}
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
    </div>), screwStepNS);
}


export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(ScrewStep);

