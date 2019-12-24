// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import checkStepActions from '../../../modules/step/checkStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { tStepPayload } from '../../../modules/step/interface/typeDef';
import styles from './styles';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  type: string,
  payload: ?tStepPayload
|};

type tDP = {|
  submit: Dispatch
|};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  type: stepPayload(viewingStep(state.order))?.type || '',
  payload: stepPayload(viewingStep(state.order))
});

const mapDispatch: tDP = {
  submit: checkStepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function CheckStep({ isCurrent, submit }: Props) {
  const classes = makeStyles(styles)();

  return (
    <div className={classes.root}>
      <Button
        variant="contained"
        color="danger"
        className={classes.button}
        disabled={!isCurrent}
        onClick={() => submit(false)}
      >
        Fail
      </Button>
      <Button
        variant="contained"
        color="info"
        className={classes.button}
        disabled={!isCurrent}
        onClick={() => submit(true)}
      >
        Pass
      </Button>
    </div>
  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(CheckStep);
