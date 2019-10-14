// @flow
import React, { useEffect, useState } from 'react';
import { Typography } from '@material-ui/core';
import { durationString } from '../../common/utils';
import type { IWorkStep } from '../../modules/step/interface/IWorkStep';

type Props = {
  step: IWorkStep
};

export default function Timer({ step }: Props) {
  const [duration, setDuration] = useState(step.timeCost());
  useEffect(() => {
    if (step.timeCost) {
      const interval = setInterval(() => {
        setDuration(step.timeCost());
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [step]);


  return <Typography variant="h6">
    {durationString(duration)}
  </Typography>;
}
