// @flow
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles } from '@material-ui/core';
import TextField from '@material-ui/core/TextField';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import stepActions from '../../../modules/step/actions';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { tStepPayload } from '../../../modules/step/interface/typeDef';
import withKeyboard from '../../../components/Keyboard';
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
  submit: stepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function MeasureStep({
                       step,
                       isCurrent,
                       submit,
                       bindAction,
                       keyboardInput,
                       target,
                       max,
                       min
                     }: Props) {
  const [value, setValue] = useState('');
  const classes = makeStyles(styles)();
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
    [step, bindAction, isCurrent, value, submit]
  );

  useEffect(
    () => {
      setValue('');
    },
    [step]
  );

  return <div className={classes.root}>
    <TextField
      label="目标值"
      disabled
      margin="normal"
      value={target || '未指定'}
      variant="outlined"
    />
    <TextField
      label="最小值"
      disabled
      margin="normal"
      value={min || '未指定'}
      variant="outlined"
    />
    <TextField
      label="最大值"
      disabled
      margin="normal"
      value={max || '未指定'}
      variant="outlined"
    />
    <TextField
      label="测量值"
      disabled={!isCurrent}
      // color="primary"
      margin="normal"
      value={value || ''}
      variant="outlined"
      onClick={() => {
        if (isCurrent) {
          keyboardInput({
            onSubmit: text => {
              setValue(text);
            },
            text: value || '',
            title: '请输入测量值',
            label: '请输入测量值'
          });
        }
      }}
    />
  </div>;
}


export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(withKeyboard(MeasureStep));
