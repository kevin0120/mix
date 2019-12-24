// @flow
import React from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import stepActions from '../../../modules/step/actions';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import styles from './styles';

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
});

const mapDispatch: tDP = {
  submit: stepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function PassFailStep({ isCurrent, submit }: Props) {
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
)(PassFailStep);
