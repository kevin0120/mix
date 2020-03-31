// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles, Typography } from '@material-ui/core';
import { Check, Close } from '@material-ui/icons';
import stepActions from '../../../modules/step/actions';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import styles from './styles';
import { stepData, viewingStep } from '../../../modules/order/selector';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP
|};

type tDP = {|
  submit: Dispatch
|};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  value: stepData(viewingStep(state.order)).result
});

const mapDispatch: tDP = {
  submit: stepActions.submit,
  passFailInput: stepActions.input
};

type Props = {|
  ...tSP,
  ...tDP
|};

function PassFailStep({ bindAction, step, isCurrent, submit, passFailInput, value }: Props) {
  const classes = makeStyles(styles)();
  const { text } = step;
  useEffect(
    () => {
      bindAction(
        <Button
          type="button"
          color="primary"
          onClick={() => {
            submit(value);
          }}
          disabled={!isCurrent}
        >
          完成
        </Button>
      );
      return () => bindAction(null);
    },
    [step, isCurrent, value, submit, bindAction]
  );

  return (
    <div className={classes.root}>
      {text ? (
        <div className={classes.textContainer}>
          <Typography variant="h4">
            {text}
          </Typography>
        </div>
      ) : null}
      <div className={classes.status}>
        {value ? <Check style={{ fontSize: 40 }} color="primary"/> : <Close style={{ fontSize: 40 }} color="error"/>}
        <Typography variant="h4">
        {value ? 'Pass' : 'Fail'}
        </Typography>
      </div>
      <div className={classes.buttonsContainer}>
        <Button
          variant="contained"
          color="danger"
          className={classes.button}
          disabled={!isCurrent}
          onClick={() => passFailInput({
            data: false,
            time: new Date(),
            source: 'input'
          })}
        >
          <Close style={{ fontSize: 40 }}/>
          <Typography variant="h4">
            Fail
          </Typography>
        </Button>
        <Button
          variant="contained"
          color="info"
          className={classes.button}
          disabled={!isCurrent}
          onClick={() => passFailInput({
            data: true,
            time: new Date(),
            source: 'input'
          })}
        >
          <Check style={{ fontSize: 40 }}/>
          <Typography variant="h4">
            Pass
          </Typography>
        </Button>
      </div>
    </div>
  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(PassFailStep);
