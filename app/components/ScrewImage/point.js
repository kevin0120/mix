import React from 'react';
import { makeStyles } from '@material-ui/core';
import clsx from 'clsx';
import styles from './style';
import { POINT_STATUS } from '../../modules/step/screwStep/constents';


const pointRadius = 30;

const isStatus = {
  waiting: status => !status || status === POINT_STATUS.WAITING || status === POINT_STATUS.WAITING_ACTIVE,
  success: status => status === POINT_STATUS.SUCCESS,
  error: status => status === POINT_STATUS.ERROR || status === POINT_STATUS.ERROR_ACTIVE,
  active: (status, twinkle) => twinkle ? status === POINT_STATUS.WAITING_ACTIVE || status === POINT_STATUS.ERROR_ACTIVE : false
};


function Point({ twinkle, x, y, status, label, scale, ...restProps }) {
  const classes = makeStyles(styles.point(pointRadius, scale))();
  return <div
    className={clsx(
      classes.root, {
        [classes.waiting]: isStatus.waiting(status),
        [classes.success]: isStatus.success(status),
        [classes.error]: isStatus.error(status),
        [classes.active]: isStatus.active(status, twinkle || false)
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
