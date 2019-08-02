import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import Grid from '@material-ui/core/Grid';

import { I18n } from 'react-i18next';

// import ButtonBase from '@material-ui/core/ButtonBase';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import PowerSettingIcon from '@material-ui/icons/PowerSettingsNew';
import ExitToApp from '@material-ui/icons/ExitToApp';

import { withStyles } from '@material-ui/core/styles';
import { push } from 'connected-react-router';
import Button from '../../components/CustomButtons/Button';

// import * as AuthActions from '../../actions/userAuth';

import ShutdownDiag from '../../components/ShutDownDiag';
import styles from './styles';
import { shutDownAction } from '../../modules/power/action';
import { setNewNotification } from '../../modules/notification/action';
import { logoutRequest } from '../../modules/user/action';
import PageEntrance from '../../components/pageEntrance';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  users: state.users,
  authEnable: state.setting.systemSettings.authEnable,
  pagesConfig: state.setting.pages,
  ...ownProps
});

const mapDispatchToProps = {
  doPush: push,
  notification: setNewNotification,
  logout: logoutRequest,
  // OpenShutdown
  shutDown: shutDownAction
};

// function mapDispatchToProps(dispatch) {
//   return bindActionCreators({ OpenShutdown }, dispatch);
// }

/* eslint-disable no-unused-vars */
class ConnectedWelcome extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      showDiag: false
    };
  }

  handleClickOpen = () => {
    this.setState({ showDiag: true });
  };

  handleShutDown = () => {
    const { shutDown } = this.props;
    this.setState({ showDiag: false });
    shutDown();
  };

  handleCloseDiag = () => {
    this.setState({ showDiag: false });
  };

  handleLogOut = () => {
    const { logout } = this.props;
    logout();
  };

  render() {
    const { classes, authEnable, doPush, notification, users, childRoutes } = this.props;
    const { showDiag } = this.state;
    const fabRightClassName = classNames(classes.fabRight);
    const fabLeftClassName = classNames(classes.fabLeft);
    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.root}>
            <PageEntrance
              type="card"
              routes={childRoutes}
              onItemClick={(route) => {
                if (!route.role || route.role.length === 0 || users.some((u) => lodash.includes(route.role, u.role))) {
                  doPush(route.url);
                } else {
                  notification('Error', '没有访问权限');
                }
              }}
              navigationClassName={classes.BottomNavigation}
              ActionClassName={classes.BottomNavigationIcon}
            />
            <Button
              round
              size="lg"
              color="danger"
              className={fabRightClassName}
              onClick={() => {
                this.handleClickOpen();
              }}
            >
              <PowerSettingIcon className={classes.extendedIcon}/>
              {t('Common.Shutdown')}
            </Button>
            {authEnable ? (
              <Button
                round
                size="lg"
                color="info"
                className={fabLeftClassName}
                component={Link}
                to="/pages/login"
              >
                <ExitToApp className={classes.extendedIcon}/>
                {t('Common.Logout')}
              </Button>
            ) : null}

            <ShutdownDiag
              show={showDiag}
              title={t('Common.Shutdown')}
              onConfirm={this.handleShutDown}
              onCancel={this.handleCloseDiag}
              content={t('Common.QuestShutdown')}
              showCancel
            />
          </div>
        )}
      </I18n>
    );
  }
}

/* eslint-enable no-unused-vars */
ConnectedWelcome.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  authEnable: PropTypes.bool.isRequired,
  // functions
  shutDown: PropTypes.func.isRequired
};

const Welcome = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWelcome);

export default withStyles(styles)(Welcome);
