import React from 'react';
import PropTypes from 'prop-types';
import AppBarBack from '../../components/AppBarBack';
import withLayout from '../../components/Layout/layout';
import Paper from '@material-ui/core/Paper';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Divider from '@material-ui/core/Divider';

import { withStyles } from '@material-ui/core/styles';

import styles from './styles';

const fakeHelp = [
  {
    key: 'Version',
    displayTitle: '版本',
    value: 'v1.0.1-beta'
  },
  {
    key: 'ReleaseDate',
    displayTitle: '发布时间',
    value: '2019-01-23'
  },
  {
    key: 'Contact',
    displayTitle: '售后电话',
    value: '400-810-2333'
  },
  {
    key: 'Company',
    displayTitle: '公司',
    value: '上海途泰工具技术有限公司'
  },
  {
    key: 'ContactPerson',
    displayTitle: '联系人',
    value: 'Henry Cai'
  },
  {
    key: 'ContactCellPhone',
    displayTitle: '联系人电话',
    value: '18701858083'
  }
];

class Help extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      // activeMenu: '网络配置',
    };
  }

  render() {
    const { classes } = this.props;
    const helpList = fakeHelp.map((h, idx) => {
      if (idx >= fakeHelp.length - 1) {
        return (
          <div key={h.key}>
            <ListItem className={classes.item}>
              {h.displayTitle}: {h.value}
            </ListItem>
          </div>
        );
      }
      return (
        <div key={h.key}>
          <ListItem className={classes.item}>
            {h.displayTitle}: {h.value}
          </ListItem>
          <li>
            <Divider />
          </li>
        </div>
      );
    });

    return (
      <div className={classes.root}>
        <AppBarBack />
        <main className={classes.content}>
          <div className={classes.toolbar} />
          <div className={classes.wrap}>
            <h2 className={classes.title}>帮助信息</h2>
            <Paper elevation={1}>
              <List>{helpList}</List>
            </Paper>
          </div>
        </main>
      </div>
    );
  }
}

Help.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

export default withStyles(styles)(Help);
