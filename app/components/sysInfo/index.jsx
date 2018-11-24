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

import * as Utils from '../../common/utils';

import styles from './styles';

const si = require('systeminformation');

class SysInfo extends React.Component {

  constructor(){
    super();
    this.timer = null;
    this.state = {
      cpustat: null,
      fsstat: null,
      memstat: null,
      batterystat: null,
    }
  }

  componentDidMount() {
    this.timer = setInterval(() => {
      si.currentLoad().then(data => {
        console.log(data);
      });
      si.fsSize().then(data => {
        console.log(data);
      });
      si.mem().then(data => {
        console.log(data);
      });
      si.battery().then(data => {
        console.log(data);
      });
    }, 2000);
  }

  componentWillUnmount() {
    clearInterval(this.timer);
  }

  render() {
    const { classes } = this.props;

    const styleOptions = { disableGutters: false };

    return (
      <I18n ns="translations">
        {t => (
          <ListItem
            key="1"
            disableGutters={styleOptions.disableGutters}
            className={classes.infoItem}
          >
            <ListItemText
              className={classes.infoText}
              primary="111"
              secondary='2222'
            />
          </ListItem>
        )}
      </I18n>
    );
  }
}

SysInfo.propTypes = {
  // styles
  classes: PropTypes.shape({}).isRequired,
  // functions
};


export default withStyles(styles)(SysInfo);
