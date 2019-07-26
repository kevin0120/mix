import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { Typography } from '@material-ui/core';
import { durationString } from '../../common/utils';
import * as orderSelectors from '../../modules/order/selector';

const mapState = (state, props) => ({
  ...props,
  start: orderSelectors.startTime(orderSelectors.viewingStep(state.order)) || null,
  end: orderSelectors.endTime(orderSelectors.viewingStep(state.order)) || null
});

const mapDispatch = {};

const getDuration = (start, end) => {
  if (!start) {
    return 0;
  }
  if (!end) {
    return new Date() - start;
  }
  return end - start;
};

function Timer({ start, end }) {
  const [duration, setDuration] = useState(getDuration(start, end));

  useEffect(() => {
    setDuration(getDuration(start, end));
    const interval = setInterval(() => {
      setDuration(getDuration(start, end));
    }, 1000);
    return () => clearInterval(interval);
  }, [end, start]);


  return <Typography variant="h4">
    {durationString(duration)}
  </Typography>;
}

export default connect(mapState, mapDispatch)(Timer);
