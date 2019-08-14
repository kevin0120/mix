import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Document, Page } from 'react-pdf/dist/entry.webpack';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Button from '../../../components/CustomButtons/Button';
import { instructionStepActions } from '../../../modules/step/instructionStep/action';
import { tStepProps } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import demoPDF from '../../../../resources/demoPDF.pdf';
import styles from './styles';

const mapState = (state, props) => ({
  ...props,
  instruction: stepPayload(viewingStep(state.order))?.instruction || null
});

const mapDispatch = {
  submit: instructionStepActions.submit
};

type Props = {
  submit: () => {}
};

function InstructionStep({ step, isCurrent, submit, bindAction, instruction }: Props & tStepProps) {
  const classes=makeStyles(styles)();
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
        submit
      </Button>
    );
    return () => bindAction(null);
  }, [step, bindAction, isCurrent, submit]);

  return (
    <div className={classes.container}>
      <Document
        className={classes.document}
        file={demoPDF}
        style={{
          width: '100%',
          height: '100%'
        }}
      >
        <Page
          className={classes.page}
          scale={1}
          pageIndex={0}
          renderAnnotationLayer={false}
        />
      </Document>
    </div>

  );
}

export default connect(mapState, mapDispatch)(InstructionStep);
