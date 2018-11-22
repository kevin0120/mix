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
import * as AuthActions from '../../actions/userAuth';

import withLayout from '../../components/Layout/layout';
import ShutdownDiag from '../../components/ShutDownDiag';
import { routeConfigs } from '../../routes';
import styles from './index.scss';

const mapStateToProps = (state, ownProps) => ({
  userUUID: state.setting.base.userInfo.uuid,
  authEnable: state.setting.systemSettings.authEnable,
  ...ownProps
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ ...AuthActions, OpenShutdown }, dispatch);
}

const withstyles = theme => ({
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    minWidth: 300,
    width: '100%',
    height: '100%',
    backgroundColor: '#232c39',
    backgroundImage:
      'linear-gradient(45deg, rgba(0, 216, 255, 0.5) 10%, rgba(0, 1, 127, 0.7))'
  },
  container: {
    padding: '40px 10px',
    textAlign: 'center',
    marginTop: 0
  },
  card: {
    width: '15%',
    minWidth: '280px',
    height: '320px',
    borderRadius: '4px'
  },
  cardActionArea: {
    height: '100%'
  },
  cardGridItem: {
    paddingRight: '30px'
  },
  cardContent: {
    position: 'relative',
    padding: '50px 10px 0'
  },
  fabLeft: {
    position: 'absolute',
    bottom: theme.spacing.unit * 10,
    left: theme.spacing.unit * 2
  },
  fabRight: {
    position: 'absolute',
    bottom: theme.spacing.unit * 10,
    right: theme.spacing.unit * 2
  },
  fabMoveUp: {
    transform: 'translate3d(0, -46px, 0)',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.enteringScreen,
      easing: theme.transitions.easing.easeOut
    })
  },
  fabMoveDown: {
    transform: 'translate3d(0, 0, 0)',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.leavingScreen,
      easing: theme.transitions.easing.sharp
    })
  },
  media: {
    background: '#ddd',
    paddingTop: '50%', // 16:9
    width: '100%'
  },
  title: {
    fontSize: '34px',
    color: '#fff',
    marginBottom: '20px'
  },
  subTitle: {
    fontSize: '16px',
    color: '#fff'
  },
  btnWrap: {
    position: 'relative',
    height: 200,
    [theme.breakpoints.down('xs')]: {
      width: '100% !important', // Overrides inline-style
      height: 100
    },
    '&:hover, &$focusVisible': {
      zIndex: 1,
      opacity: 0.65,
      '& $titleMarked': {
        opacity: 0
      },
      '& $navTitle': {
        border: '4px solid currentColor'
      }
    }
  },
  focusVisible: {},
  extendedIcon: {
    marginRight: theme.spacing.unit
  },
  itemButton: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: theme.palette.common.white
  },
  navTitle: {
    position: 'relative',
    padding: `${theme.spacing.unit * 2}px ${theme.spacing.unit * 4}px ${theme
      .spacing.unit + 6}px`
  },
  titleMarked: {
    height: 3,
    width: 18,
    backgroundColor: theme.palette.common.white,
    position: 'absolute',
    bottom: -2,
    left: 'calc(50% - 9px)',
    transition: theme.transitions.create('opacity')
  }
});
/* eslint-disable no-unused-vars */
class ConnectedWelcome extends React.Component {
  componentDidMount() {
    const { doUserAuth, userUUID } = this.props;
    doUserAuth(userUUID);
  }

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
                        <div className={styles.iconWrap}>
                          <route.icon />
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
  userUUID: PropTypes.string.isRequired,
  // functions
  doUserAuth: PropTypes.func.isRequired,
  OpenShutdown: PropTypes.func.isRequired,
  userLogOut: PropTypes.func.isRequired
};

const Welcome = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWelcome);

export default withLayout(withStyles(withstyles)(Welcome));
