// @flow
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Typography } from '@material-ui/core';
import { durationString, timeCost } from '../../common/utils';
import * as orderSelectors from '../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  times: orderSelectors.times(orderSelectors.viewingStep(state.order)) || null,
});

const mapDispatch = {};

type Props = {
  times: Array<Date>
};

function Timer({ times }: Props) {
  const [duration, setDuration] = useState(timeCost(times));
  useEffect(() => {
    setDuration(timeCost(times));
    const interval = setInterval(() => {
      setDuration(timeCost(times));
    }, 1000);
    return () => clearInterval(interval);
  }, [times]);


  return <Typography variant="h4">
    {durationString(duration)}
  </Typography>;
}

export default connect(mapState, mapDispatch)(Timer);
