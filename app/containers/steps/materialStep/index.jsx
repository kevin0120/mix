// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { Typography } from '@material-ui/core';
import materialStepActions from '../../../modules/step/materialStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { IMaterialStep } from '../../../modules/step/materialStep/interface/IMaterialStep';

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
  const { payload, description } = step;
  // eslint-disable-next-line react-hooks/rules-of-hooks
  useEffect(() => {
    bindAction(
      isCurrent ? <Button onClick={() => ready()} color="primary">
        完成
      </Button> : null
    );
    bindDescription(
      <React.Fragment>
        <Typography variant="h5" style={{ paddingBottom: '10px' }}>{description}</Typography>
        {(payload?.items || []).map(i =>
          <Typography key={i.name} variant="body1">{i.name}</Typography>
        )}
      </React.Fragment>
    );
    return () => {
      bindAction(null);
      bindDescription(null);
    };
  }, [bindAction, bindDescription, description, isCurrent, payload, ready]);


  return <div
    style={
      {
        display: 'flex',
        justifyContents: 'center',
        alignItems: 'center',
        width: '100%',
        height: '100%'

      }
    }>
    {payload.items.map(i => <img
      style={
        {
          maxWidth: '100%',
          maxHeight: '100%',
          margin: 'auto',
          justifySelf: 'center',
          alignSelf: 'center',
          resize: 'both'
        }
      }
      key={i.name} src={i.image || ''} alt={i.name}
    />)}
  </div>;
}

export default connect<Props, ownProps, _, _, _, _>(
  mapState,
  mapDispatch
)(materialStep);
