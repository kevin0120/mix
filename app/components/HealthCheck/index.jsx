import React from 'react';
import PropTypes from 'prop-types';
import { I18n } from 'react-i18next';
import { withStyles } from '@material-ui/core/styles';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuList from '@material-ui/core/MenuList';
import styles from './styles';


class HealthCheck extends React.Component {
  render() {
    const { status, classes } = this.props;


    const styleOptions = { disableGutters: false };

    return Object.keys(status).sort((a, b) => a.localeCompare(b, 'zh-CN')).map(
      (k) => {
        const isHealth = status[k];
        const displayTitle = k;
        const statusClassName = isHealth
          ? classes.infoSuccess
          : classes.infoError;

        return (
          <I18n ns="translations" key={`health-check-${k}`}>
            {t => (
              <ListItem
                disableGutters={styleOptions.disableGutters}
                className={classes.infoItem}
              >
                <ListItemText
                  className={classes.infoText}
                  primary={displayTitle}
                />
                <div className={`${classes.infoStatus} ${statusClassName}`}/>
              </ListItem>
            )}
          </I18n>
        );
      }
    );
  }

}

HealthCheck.propTypes = {
  // styles
  classes: PropTypes.shape({}).isRequired,
  status: PropTypes.shape({}).isRequired
  // functions
};

export default withStyles(styles)(HealthCheck);
