import React from 'react';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';

import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuList from '@material-ui/core/MenuList';

import * as Utils from '../../common/utils';

import styles from './styles';

const lodash = require('lodash');

class ConnectedHealthCheck extends React.Component {
  render() {
    const { classes, healthCheckResults } = this.props;

    const results = lodash.filter(healthCheckResults, 'enable');

    const styleOptions = { disableGutters: false };

    return Utils.sortObj(results, 'displayOrder').map(
      ({ key, value: item }) => {
        const { isHealth, displayTitle } = item;
        const statusClassName = isHealth
          ? classes.infoSuccess
          : classes.infoError;

        return (
          <I18n ns="translations">
            {t => (
              <MenuList key={key}>
                <ListItem
                  key={key}
                  disableGutters={styleOptions.disableGutters}
                  className={classes.infoItem}
                >
                  <ListItemText
                    key={key}
                    className={classes.infoText}
                    primary={t(displayTitle)}
                  />
                  <div className={`${classes.infoStatus} ${statusClassName}`} />
                </ListItem>
              </MenuList>
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
