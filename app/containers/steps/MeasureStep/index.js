// @flow
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import checkStepActions from '../../../modules/step/checkStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { tStepPayload } from '../../../modules/step/interface/typeDef';

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

function CheckStep({
                     step,
                     isCurrent,
                     submit,
                     bindAction,
                   }: Props) {
  const [value, setValue] = useState('');

  useEffect(
    () => {
      const onSubmit = v => {
        submit(v);
      };
      bindAction(
        <Button
          type="button"
          color="primary"
          onClick={() => {
            onSubmit(value);
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

  return <div>
    Measure
  </div>;
}


export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(CheckStep);
