import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import checkStepActions from '../../../modules/step/checkStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import Qualitative from './Qualitative';
import Quantitative from './Quantitative';
import { checkStepTypes } from '../../../modules/step/checkStep/constants';

const Contents = {
  [checkStepTypes.measurement]: Quantitative,
  [checkStepTypes.passFail]: Qualitative
};

const mapState = (state, props) => ({
  ...props,
  type: stepPayload(viewingStep(state.order))?.type || '',
  payload: stepPayload(viewingStep(state.order))
});

const mapDispatch = {
  submit: checkStepActions.submit
};

type Props = {
  type: string,
  submit: () => {}
};

function CheckStep({
                     step,
                     type,
                     isCurrent,
                     submit,
                     bindAction,
                     payload
                   }: Props & tStepProps) {
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
  const RenderContent = Contents[type];

  return (
    <div>
      {
        RenderContent ? <RenderContent payload={payload}/> : null
      }
    </div>
  );
}


export default connect(
  mapState,
  mapDispatch
)(CheckStep);
