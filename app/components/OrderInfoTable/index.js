import React from 'react';
import { makeStyles, Typography } from '@material-ui/core';
import { tNS, withI18n } from '../../i18n';
import Table from '../Table/Table';
import type { IWorkStep } from '../../modules/step/interface/IWorkStep';
import { stepWorkingNS } from '../../containers/stepWorking/local';
import styles from './styles';

function OrderInfoRow(steps, classes) {
  return (steps &&
    steps.map((s: IWorkStep, idx) => [
      idx + 1,
      s.code,
      tNS(`${s.type}`, stepWorkingNS),
      <div dangerouslySetInnerHTML={{ __html: s.desc || '' }} className={classes.text}/>
    ])) || [];
}

export default function OrderInfoTable({ steps }) {
  const classes = makeStyles(styles)();
  return withI18n(t => <Table
    tableHeaderColor="info"
    tableHead={[
      t('Common.Idx'),
      t('Order.Step.name'),
      t('Order.Step.type'),
      t('Order.Step.desc')
    ]}
    tableData={OrderInfoRow(steps, classes)}
    colorsColls={['info']}
  />);
}