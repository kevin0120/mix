import React from 'react';
import { makeStyles } from '@material-ui/core';
import clsx from 'clsx';
import styles from './style';

const isStatus = {
  waiting: status => status === 'waitingActive' || status === 'waiting',
  active: status => status === 'waitingActive' || status === 'errorActive',
  success: status => status === 'success',
  error: status => status === 'error' || status === 'errorActive'
};


function Point({ x, y, status, label }) {
  const classes = makeStyles(styles.point)();
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
      left: `${x}%`,
      top: `${y}%`
    }}
  >
    {label}
  </div>;
}

export default Point;
