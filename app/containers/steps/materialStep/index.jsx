// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { makeStyles, Typography } from '@material-ui/core';
import materialStepActions from '../../../modules/step/materialStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { IMaterialStep } from '../../../modules/step/materialStep/interface/IMaterialStep';
import styles from './styles';

const mapState = (state, props: ownProps): Props => ({
  ...props
});

const mapDispatch = {
  ready: materialStepActions.ready
};

type Props = {|
  ...ownProps
|};

type ownProps = {|
  ...tStepProps,
  ready: Dispatch,
  step: IMaterialStep
|};

function materialStep(props: Props) {
  const { step, bindAction, ready, isCurrent, bindDescription } = props;
  const { payload, desc, consumeProduct } = step;
  // eslint-disable-next-line react-hooks/rules-of-hooks
  useEffect(() => {
    bindAction(
      isCurrent ? <Button onClick={() => ready()} color="primary">
        完成
      </Button> : null
    );
    // bindDescription(
    //   <React.Fragment>
    //
    //   </React.Fragment>
    // );
    return () => {
      bindAction(null);
      bindDescription(null);
    };
  }, [bindAction, bindDescription, consumeProduct, desc, isCurrent, payload, ready]);
  const classes = makeStyles(styles)();

  return <div className={classes.root}>
    <span className={classes.desc}>

    {desc ? (
      <Typography variant="h4"  dangerouslySetInnerHTML={{ __html: desc || '' }} />
    ) : null}
      </span>

    {[consumeProduct].map(i =>
      <span>
        <Typography key={i} variant="h5">{i}</Typography>
      </span>
    )}

  </div>;
}

export default connect<Props, ownProps, _, _, _, _>(
  mapState,
  mapDispatch
)(materialStep);
