// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Button from '../../../components/CustomButtons/Button';
import stepActions  from '../../../modules/step/actions';
import type { tStepProps } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import styles from './styles';
import type { Dispatch } from '../../../modules/typeDef';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  url: string,
  page: number
|};

type tDP = {|
  submit: Dispatch
|};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  url: stepPayload(viewingStep(state.order))?.url || '',
  page: stepPayload(viewingStep(state.order))?.page || 0
});

const mapDispatch: tDP = {
  submit: stepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function InstructionStep({ step, isCurrent, submit, bindAction }: Props) {
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
      text
    </div>

  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(mapState, mapDispatch)(InstructionStep);
