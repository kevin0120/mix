// @flow
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Typography } from '@material-ui/core';
import { durationString, timeCost } from '../../common/utils';
import * as orderSelectors from '../../modules/order/selector';
import type { tStep } from '../../modules/order/model';

// const mapState = (state, props) => ({
//   ...props,
//   step: orderSelectors.workingStep(orderSelectors.workingOrder(state.order)) || {},
// });
//
// const mapDispatch = {};

type Props = {
  step: tStep
};

export default function Timer({ step }: Props) {
  const times = orderSelectors.times(step);
  const [duration, setDuration] = useState(timeCost(times));
  useEffect(() => {
    if (times && (times.length > 0)) {
      setDuration(timeCost(times));
      const interval = setInterval(() => {
        setDuration(timeCost(times));
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [times]);


  return <Typography variant="h6">
    {durationString(duration)}
  </Typography>;
}

// export default connect(mapState, mapDispatch)(Timer);
