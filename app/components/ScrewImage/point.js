import React from 'react';
import { makeStyles } from '@material-ui/core';
import clsx from 'clsx';
import styles from './style';


const pointRadius = 30;

const isStatus = {
  waiting: status => !status || status === 'waiting' || status === 'waitingActive',
  success: status => status === 'success',
  error: status => status === 'error' || status === 'failActive',
  active: status => status === 'waitingActive' || status === 'failActive'
};


function Point({ x, y, status, label, scale, ...restProps }) {
  const classes = makeStyles(styles.point(pointRadius, scale))();
  return <div
    className={clsx(
      classes.root, {
        [classes.waiting]: isStatus.waiting(status),
        [classes.success]: isStatus.success(status),
        [classes.error]: isStatus.error(status),
        [classes.active]: isStatus.active(status)
      })}
    style={{
      position: 'absolute',
      left: `calc(${x}% - ${pointRadius}px)`,
      top: `calc(${y}% - ${pointRadius}px)`
    }}
    {...restProps}
  >
    {label}
  </div>;
}

export default Point;
