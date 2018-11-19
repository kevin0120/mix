import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Button from '@material-ui/core/Button';
import Apps from '@material-ui/icons/Apps';
import { I18n } from 'react-i18next';

import styles from './styles';

/* eslint-disable react/prefer-stateless-function */
function AppBarBack(props) {
  const { classes } = props;

  function gotoWelcome() {
    window.location.href = '#/welcome';
  }

  return (
    <I18n ns="translations">
      {
        t => (
          <AppBar position="absolute" color="default" className={classes.root}>
            <Toolbar className={classes.appBarTool} variant="dense">
              <Button className={classes.backButton} onClick={() => gotoWelcome()}>
                <Apps className={classes.leftIcon} />
                {t('Back')}
              </Button>
              {props.children}
            </Toolbar>
          </AppBar>
        )
      }
    </I18n>
  );
}

AppBarBack.propTypes = {
  classes: PropTypes.shape({
  }).isRequired,
  children: PropTypes.shape({}),
};

AppBarBack.defaultProps = {
  children: null,
};

export default withStyles(styles)(AppBarBack);
