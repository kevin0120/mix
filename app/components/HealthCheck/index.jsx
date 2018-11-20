import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';
import { bindActionCreators } from 'redux';

import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import * as Utils from '../../common/utils';

import  * as healthzCheckActions  from '../../actions/healthCheck';

import styles from './styles';

function mapDispatchToProps(dispatch){
  return bindActionCreators(healthzCheckActions, dispatch);
}

class ConnectedHealthCheck extends React.Component {
  constructor(props) {
    super(props);
  }

  componentDidMount() {
    const { startHealthzCheck } = this.props;
    startHealthzCheck();
  }

  componentWillUnmount() {
    const { stopHealthzCheck } = this.props;
    stopHealthzCheck();
  }

  render() {
    const { classes, healthCheckResults } = this.props;

    const styleOptions = {disableGutters: false};

    return Utils.sortObj(healthCheckResults, 'displayOrder').map(
      ({ key, value: item }) => {
        const { health, displayTitle } = item;
        const statusClassName = health
          ? classes.infoSuccess
          : classes.infoError;

        return (
          <I18n ns="translations">
            {t => (
              <ListItem
                key={key}
                disableGutters={styleOptions.disableGutters}
                className={classes.infoItem}
              >
                <ListItemText
                  className={classes.infoText}
                  primary={t(displayTitle)}
                />
                <div className={`${classes.infoStatus} ${statusClassName}`} />
              </ListItem>
            )}
          </I18n>
        );
      }
    );
  }
}

ConnectedHealthCheck.propTypes = {
  // styles
  classes: PropTypes.shape({}).isRequired,
  healthCheckResults: PropTypes.shape({}).isRequired,
  // functions
  stopHealthzCheck: PropTypes.func.isRequired,
  startHealthzCheck: PropTypes.func.isRequired
};

const HealthCheck = connect(mapDispatchToProps)(ConnectedHealthCheck);

export default withStyles(styles)(HealthCheck);
