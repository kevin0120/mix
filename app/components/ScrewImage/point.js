// @flow
import React from 'react';
import { makeStyles } from '@material-ui/core';
import clsx from 'clsx';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
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
  scale?: number,
  info?: Object,
  menu?: boolean
};


function Point({ info, twinkle, x, y, status, label, scale, reworkModiBGColor, onClick, selected, menu = false, ...restProps }: Props) {
  const classes = makeStyles(styles.point(pointRadius, scale))();

  return <div style={{
    width: '0px', height: '0px', left: `calc(${x}% - ${pointRadius}px)`,
    top: `calc(${y}% - ${pointRadius}px)`, position: 'relative',
    overflow: 'visible'
  }}>
    <div
      className={clsx(
        classes.root, {}, {
          [classes.waiting]: isStatus.waiting(status),
          [classes.success]: isStatus.success(status),
          [classes.error]: isStatus.error(status),
          [classes.active]: twinkle || false,
          [classes.selected]: selected
        })}
      onClick={(e) => {
        if (onClick) {
          onClick(e);
        }
      }}
      {...(restProps: any)}
    >
      {label}


    </div>
    {menu && !!twinkle && info && Object.keys(info).length > 0 ? <List style={{
      width: 200,
      // minHeight: pointRadius * 2,
      left: x <= 80 ? pointRadius * 2 + 10 : -200 - 10,
      top: 5,
      position: 'absolute',
      borderRadius: 10,
      color: 'white',
      backgroundColor: 'rgba(113,114,118,0.8)',
      fontSize: '1.2rem'
    }}>
      {Object.keys(info || {}).map(i => <ListItem style={{
        display: 'flex',
        justifyContent: 'space-between'
      }}><span>{`${i} : `}</span><span>{info[i]}</span></ListItem>)}
    </List> : null}
  </div>;
}

Point.defaultProps = {
  scale: 1,
  reworkModiBGColor: false
};

export default Point;
