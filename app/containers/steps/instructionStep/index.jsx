// @flow
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Document, Page } from 'react-pdf/dist/entry.webpack';
import makeStyles from '@material-ui/core/styles/makeStyles';
import Button from '../../../components/CustomButtons/Button';
import { instructionStepActions } from '../../../modules/step/instructionStep/action';
import type { tStepProps } from '../types';
import { stepPayload, viewingStep } from '../../../modules/order/selector';
import styles from './styles';
import { defaultClient } from '../../../common/utils';
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
  submit: instructionStepActions.submit
};

type Props = {|
  ...tSP,
  ...tDP
|};

function InstructionStep({ step, isCurrent, submit, bindAction, url, page }: Props) {
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
  const [pdfPage, setPage] = useState(0);
  const [pdfUrl, setPdfUrl] = useState('');

  useEffect(() => {
    const response = defaultClient
      .get(url)
      .then(resp => {
        setPdfUrl(resp.request._redirectable._currentUrl);
      });
  }, [url]);

  return (
    <div className={classes.container}>
      <Document
        className={classes.document}
        file={pdfUrl}
        // file={pdf}
        onLoadSuccess={() => setPage(page)}
      >
        <Page
          className={classes.page}
          scale={1.3}
          pageIndex={pdfPage}
          renderAnnotationLayer={false}
        />
      </Document>
    </div>

  );
}

export default connect<Props, tOP, tSP, tDP, _, _>(mapState, mapDispatch)(InstructionStep);
