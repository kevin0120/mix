import { Typography } from '@material-ui/core';
import Table from '../../../components/Table/Table';
import React from 'react';

function Title({ title, classes }) {
  return <div className={classes.titleContainer}>
    <Typography variant="h6" color="inherit">{title}</Typography>
  </div>;
}

export default function Panel({ title, cols = [], classes, data = [] }) {
  return (
    <div className={classes.panelRoot}>
      <Title classes={classes} title={title}/>
      <Table
        classes={{
          tableResponsive: classes.table,
          tableCell: classes.tableCell
        }}
        tableHeaderColor="info"
        tableHead={cols.map(c => c.label) || []}
        tableData={data.map(d => cols.map(c => d[c.name] || '')) || []}
        colorsColls={['info']}
      />
    </div>
  );
}
