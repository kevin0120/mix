// @flow
import React, { useState } from 'react';
import { makeStyles } from '@material-ui/core';
import clsx from 'clsx';
import styles from './style';
import { POINT_STATUS } from '../../modules/step/screwStep/constants';
import type { tPointStatus } from '../../modules/step/screwStep/interface/typeDef';


const pointRadius = 30;

const isStatus = {
  waiting: status => !status || status === POINT_STATUS.WAITING,
  success: status => status === POINT_STATUS.SUCCESS,
  error: status => status === POINT_STATUS.ERROR
};

type Props = {
  twinkle: boolean,
  x: number,
  y: number,
  onClick: (Object) => boolean,
  reworkModiBGColor?: boolean,
  status: tPointStatus,
  label: string | number,
  scale?: number
};


function Point({ twinkle, x, y, status, label, scale, reworkModiBGColor, onClick, selected, ...restProps }: Props) {
  const classes = makeStyles(styles.point(pointRadius, scale))();

  return <div
    className={clsx(
      classes.root, {
        [classes.waiting]: isStatus.waiting(status),
        [classes.success]: isStatus.success(status),
        [classes.error]: isStatus.error(status),
        [classes.active]: twinkle || false,
        [classes.selected]: selected
      })}
    style={{
      position: 'absolute',
      left: `calc(${x}% - ${pointRadius}px)`,
      top: `calc(${y}% - ${pointRadius}px)`
    }}
    onClick={(e) => {
      if (onClick) {
        onClick(e)
      }
    }}
    {...(restProps: any)}
  >
    {label}
  </div>;
}

Point.defaultProps = {
  scale: 1,
  reworkModiBGColor: false
};

export default Point;
