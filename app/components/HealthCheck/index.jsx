import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';


import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import * as Utils from '../../common/utils';

import {
  masterPCHealthCheck,
  controllerHealthCheck,
} from '../../actions/healthCheck';

import styles from './styles';


const mapStateToProps = (state, ownProps) => ({
  healthCheckResults: state.healthCheckResults,
  connInfo: state.connInfo,
  ...ownProps,
});

const mapDispatchToProps = {
  masterPCHealthCheck,
  controllerHealthCheck,
};

class ConnectedHealthCheck extends React.Component {
  constructor(props) {
    super(props);
    this.timerHealz = null;
  }

  componentDidMount() {
    console.log("health timer restart");
    this.restartTimer();
  }

  componentWillUnmount() {
    // console.log("health timer stop");
    // this.stopTimer();
  }

  stopTimer() {
    clearInterval(this.timerHealz);
    this.timerHealz = null;
  }

  restartTimer() {
    const { options } = this.props;
    const masterpcConn = this.props.connInfo.masterpc.connection;

    const Controllers = this.props.connInfo.controllers;

    if (options.enableDebugLog) {
      console.log('restart MasterPC HealthCheck. ' +
        `conn: ${masterpcConn}, ` +
        `interval: ${options.intervalMasterPC}ms`);
    }

    this.stopTimer();
    this.timerHealz = setInterval(() => {
      this.props.masterPCHealthCheck(masterpcConn);
      this.props.controllerHealthCheck(masterpcConn, Controllers);
    }, options.intervalMasterPC);
  }

  render() {
    const { classes, styleOptions, healthCheckResults } = this.props;

    return (Utils.sortObj(healthCheckResults, 'displayOrder').map(({ key, value: item }) => {
      const { health, displayTitle } = item;
      const statusClassName = health ? classes.infoSuccess : classes.infoError;

      return (
        <I18n ns="translations">
          {
            t => (
              <ListItem
                key={key}
                disableGutters={styleOptions.disableGutters}
                className={classes.infoItem}
              >
                <ListItemText
                  className={classes.infoText}
                  primary={t(displayTitle)}
                />
                <div className={`${classes.infoStatus} ${statusClassName}`}
                />
              </ListItem>
            )
          }
        </I18n>
      );
    }));
  }
}

ConnectedHealthCheck.propTypes = {
  // styles
  classes: PropTypes.shape({}).isRequired,
  connInfo: PropTypes.shape({
    masterpc: PropTypes.shape({
      connection: PropTypes.string,
    }),
    controllers: PropTypes.array,
  }),
  healthCheckResults: PropTypes.shape({}).isRequired,
  styleOptions: PropTypes.shape({
    disableGutters: PropTypes.bool,
  }),
  // Connection Options
  options: PropTypes.shape({
    intervalMasterPC: PropTypes.number,
    intervalControllerSN: PropTypes.number,
    enableDebugLog: PropTypes.bool,
  }),
  // functions
  masterPCHealthCheck: PropTypes.func.isRequired,
  controllerHealthCheck: PropTypes.func.isRequired,
};

ConnectedHealthCheck.defaultProps = {
  styleOptions: {
    disableGutters: false,
  },
  connInfo: {
    masterpc: {
      connection: null,
    },
    controllers: [],
  },
  options: {
    intervalMasterPC: 3000,
    intervalControllerSN: 3000,
    enableDebugLog: false,
  },
};

const HealthCheck = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ConnectedHealthCheck);

export default withStyles(styles)(HealthCheck);
