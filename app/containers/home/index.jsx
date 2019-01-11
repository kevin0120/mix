import React from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import { bindActionCreators } from 'redux';

import Grid from '@material-ui/core/Grid';

import { I18n } from 'react-i18next';

// import ButtonBase from '@material-ui/core/ButtonBase';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import PowerSettingIcon from '@material-ui/icons/PowerSettingsNew';
import ExitToApp from '@material-ui/icons/ExitToApp';

import { withStyles } from '@material-ui/core/styles';
import Button from '../../components/CustomButtons/Button';

import { openShutdown as OpenShutdown } from '../../actions/shutDownDiag';
// import * as AuthActions from '../../actions/userAuth';

import withLayout from '../../components/Layout/layout';
import ShutdownDiag from '../../components/ShutDownDiag';
import { routeConfigs } from '../../routes/index';
import styles from './styles';
import { push } from "connected-react-router";
import { setNewNotification } from '../../actions/notification';
const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  usersInfo: state.users,
  authEnable: state.setting.systemSettings.authEnable,
  ...ownProps
});

const mapDispatchToProps = {
  doPush: push,
  notification:setNewNotification,
  OpenShutdown
};

// function mapDispatchToProps(dispatch) {
//   return bindActionCreators({ OpenShutdown }, dispatch);
// }

/* eslint-disable no-unused-vars */
class ConnectedWelcome extends React.Component {
  // componentDidMount() {
  //   const { doUserAuth, userUUID } = this.props;
  //   doUserAuth(userUUID);
  // }

  handleClickOpen = () => {
    const { OpenShutdown } = this.props;
    OpenShutdown('shutdown');
  };

  handleLogOut = () => {
    const { userLogOut } = this.props;
    userLogOut();
  };

  render() {
    const { classes, authEnable,doPush,notification,usersInfo } = this.props;
    const fabRightClassName = classNames(classes.fabRight);
    const fabLeftClassName = classNames(classes.fabLeft);
    const { role } = usersInfo[0];
    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.root}>
            <Grid container className={classes.container} justify="center">
              {routeConfigs.slice(1, -1).map(route => (
                <Grid key={route.name} item className={classes.cardGridItem}>
                  <Card
                    key={route.name}
                    className={classes.card}
                    style={{ backgroundColor: route.color }}
                  >
                    <CardActionArea
                      // component={Link}
                      // to={route.url}
                      onClick={() => {
                        if (route.roles && lodash.includes(route.roles, role)) {
                          doPush(route.url);
                        } else {
                          notification('error','没有访问权限');
                        }
                      }}
                      className={classes.cardActionArea}
                    >
                      <div
                        className={classes.media}
                        style={{ backgroundImage: `url(${route.image})` }}
                      />
                      <CardContent className={classes.cardContent}>
                        <div className={classes.iconWrap}>
                          <route.icon className={classes.icon} />
                        </div>
                        <h1 className={classes.title}>{t(route.title)}</h1>
                        <p className={classes.subTitle}>{t(route.title,{lng:'en'})}</p>
                      </CardContent>
                    </CardActionArea>
                  </Card>
                </Grid>
              ))}
            </Grid>
            <Button
              round
              size="lg"
              color="danger"
              className={fabRightClassName}
              onClick={this.handleClickOpen}
            >
              <PowerSettingIcon className={classes.extendedIcon} />
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
                <ExitToApp className={classes.extendedIcon} />
                {t('Common.Logout')}
              </Button>
            ) : null}

            <ShutdownDiag />
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
  OpenShutdown: PropTypes.func.isRequired
};

const Welcome = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWelcome);

export default withStyles(styles)(Welcome);
