// @flow
import React, { useEffect } from 'react';
import { connect } from 'react-redux';
import { Typography } from '@material-ui/core';
import materialStepActions from '../../../modules/step/materialStep/action';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/indexReducer';

type Props = {
  ready: Dispatch
};

function materialStep(props: Props & tStepProps) {
  const { step, bindAction, ready, isCurrent, bindDescription } = props;
  const { payload, description } = step;
  useEffect(() => {
    bindAction(
      isCurrent ? <Button onClick={() => ready()} color="primary">
        ready
      </Button> : null
    );
    bindDescription(
        <React.Fragment>
          <Typography variant="h5">{description}</Typography>
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
    style={{
      display:'flex',
      justifyContents:'center',
      alignItems:'center',
      width:'100%',
      height:'100%',
    }}>
    {payload.items.map(i => <img
      style={{
        maxWidth:'100%',
        maxHeight:'100%',
        margin:'auto',
        justifySelf:'center',
        alignSelf:'center',
resize:'both'
      }}
      key={i.name} src={i.image||''} alt={i.name}
    />)}
  </div>;
}

const mapState = (state, props) => ({
  ...props
});

const mapDispatch = {
  ready: materialStepActions.ready
};

export default connect(mapState, mapDispatch)(materialStep);
