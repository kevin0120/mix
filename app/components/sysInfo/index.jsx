/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

import React from 'react';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';

import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import MenuList from '@material-ui/core/MenuList';

import * as Utils from '../../common/utils';

import styles from './styles';

const si = require('systeminformation');

class SysInfo extends React.Component {
  constructor() {
    super();
    this.timer = null;
    this.state = {
      cpustat: { avgload: 100.0 },
      fsstat: { used: 100, size: 100 },
      memstat: { used: 100, total: 100 },
      batterystat: { currentcapacity: 0, maxcapacity: 100 }
    };
  }

  componentDidMount() {
    this.timer = setInterval(() => {
      si.currentLoad().then(data => {
        this.setState({
          cpustat: { avgload: data.avgload }
        });
      });
      si.fsSize().then(data => {
        console.log(data);
      });
      si.mem().then(data => {
        this.setState({
          memstat: { used: data.used, total: data.total }
        });
      });
      si.battery().then(data => {
        console.log(data);
      });
      si.networkInterfaces().then(data => {
        console.log(data);
      });
    }, 5000);
  }

  shouldComponentUpdate(nextProps, nextState) {
    console.log(nextProps, nextState);
    return this.state !== nextState || this.props !== nextProps;
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  render() {
    const { classes } = this.props;

    const styleOptions = { disableGutters: false };

    return Utils.normalSortObj(this.state).map(({ key, value: item }) => {
      const { isHealth, displayTitle } = item;
      const statusClassName = isHealth
        ? classes.infoSuccess
        : classes.infoError;

      return (
        <I18n ns="translations">
          {t => (
            <MenuList key={key}>
              <ListItem
                key={displayTitle}
                disableGutters={styleOptions.disableGutters}
                className={classes.infoItem}
              >
                <ListItemText
                  key={displayTitle}
                  className={classes.infoText}
                  primary={t(displayTitle)}
                />
                <div className={`${classes.infoStatus} ${statusClassName}`} />
              </ListItem>
            </MenuList>
          )}
        </I18n>
      );
    });
  }
}

SysInfo.propTypes = {
  // styles
  classes: PropTypes.shape({}).isRequired
  // functions
};

export default withStyles(styles)(SysInfo);
