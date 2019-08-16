import React from 'react';
import Button from '../../../components/CustomButtons/Button';
import { Typography } from '@material-ui/core';
import checkStepAction from '../../../modules/step/checkStep/action';


const mapState = (state, props) => ({
  ...props
});

const mapDispatch = {
  submit:checkStepAction.submit,
  cancel:checkStepAction.cancel
};

export default function({ payload, submit, cancel }) {
  const { isCurrent, instruction } = payload;
  return (
    <div>
      <Typography>
        {instruction}
      </Typography>
      <div>

        <Button
          type="button"
          color="danger"
          onClick={() => {
            submit();
          }}
          disabled={!isCurrent}
        >
          否
        </Button>
        <Button
          type="button"
          color="primary"
          onClick={() => {
            cancel();
          }}
          disabled={!isCurrent}
        >
          是
        </Button>
      </div>
    </div>
  );
}