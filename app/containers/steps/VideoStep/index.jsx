// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Button from '../../../components/CustomButtons/Button';
import { videoStepActions } from '../../../modules/step/videoStep/action';
import type { tStepProps } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import styles from './styles';
import type { Dispatch } from '../../../modules/typeDef';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  video: string
|};

type tDP = {|
  submit: Dispatch
|};
const mapState = (state, props: tOP): tSP => {
  // console.log(stepPayload(viewingStep(state.order)));
  return ({
    ...props,
    video: stepPayload(viewingStep(state.order))?.url || ''
  });
};

const mapDispatch: tDP = {
  submit: videoStepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function InstructionStep({ step, isCurrent, submit, bindAction, video }: Props) {
  const classes = makeStyles(styles)();
  useEffect(() => {
    bindAction(
      <Button
        type="button"
        color="primary"
        onClick={() => {
          submit();
        }}
        disabled={!isCurrent}
      >
        完成
      </Button>
    );
    return () => bindAction(null);
  }, [step, bindAction, isCurrent, submit]);

  return (
    <div className={classes.container}>
      <video src={video} className={classes.video} autoPlay loop controls/>
    </div>

  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(mapState, mapDispatch)(InstructionStep);
