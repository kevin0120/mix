import React from 'react';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';

import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import * as Utils from '../../common/utils';

import styles from './styles';

class ConnectedHealthCheck extends React.Component {
  render() {
    const { classes, healthCheckResults } = this.props;

    const styleOptions = { disableGutters: false };

    return Utils.sortObj(healthCheckResults, 'displayOrder').map(
      ({ key, value: item }) => {
        const { isHealth, displayTitle } = item;
        const statusClassName = isHealth
          ? classes.infoSuccess
          : classes.infoError;

        return (
          <I18n ns="translations">
            {t => (
              <ListItem
                key={displayTitle}
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
  healthCheckResults: PropTypes.shape({}).isRequired
  // functions
};

export default withStyles(styles)(ConnectedHealthCheck);
