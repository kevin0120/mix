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
import Button from '@material-ui/core/Button';
import CardContent from '@material-ui/core/CardContent';
import PowerSettingIcon from '@material-ui/icons/PowerSettingsNew';
import ExitToApp from '@material-ui/icons/ExitToApp';

import { withStyles } from '@material-ui/core/styles';

import { openShutdown as OpenShutdown } from '../../actions/shutDownDiag';
// import * as AuthActions from '../../actions/userAuth';

import withLayout from '../../components/Layout/layout';
import ShutdownDiag from '../../components/ShutDownDiag';
import { routeConfigs } from '../../routes';
import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  authEnable: state.setting.systemSettings.authEnable,
  ...ownProps
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ OpenShutdown }, dispatch);
}


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
    const { classes, authEnable } = this.props;
    const fabRightClassName = classNames(classes.fabRight);
    const fabLeftClassName = classNames(classes.fabLeft);

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
                      component={Link}
                      to={route.url}
                      className={classes.cardActionArea}
                    >
                      <div
                        className={classes.media}
                        style={{ backgroundImage: `url(${route.image})` }}
                      />
                      <CardContent className={classes.cardContent}>
                        <div className={classes.iconWrap}>
                          <route.icon className={classes.icon}/>
                        </div>
                        <h1 className={classes.title}>{t(route.title)}</h1>
                        <p className={classes.subTitle}>{route.enName}</p>
                      </CardContent>
                    </CardActionArea>
                  </Card>
                </Grid>
              ))}
            </Grid>
            <Button
              variant="extendedFab"
              color="secondary"
              className={fabRightClassName}
              onClick={this.handleClickOpen}
            >
              <PowerSettingIcon className={classes.extendedIcon} />
              {t('Common.Shutdown')}
            </Button>
            {authEnable ? (
              <Button
                variant="extendedFab"
                color="primary"
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
  OpenShutdown: PropTypes.func.isRequired,
};

const Welcome = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWelcome);

export default withLayout(withStyles(styles)(Welcome));
