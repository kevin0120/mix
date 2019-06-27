// @flow
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
/* eslint-disable no-unused-vars */
import SwipeableDrawer from '@material-ui/core/SwipeableDrawer';
import Slide from '@material-ui/core/Slide';
import { I18n, Trans } from 'react-i18next';

import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Clock from 'react-live-clock';

import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';

import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import Language from '@material-ui/icons/LanguageRounded';
import Fade from '@material-ui/core/Fade';

import MenuIcon from '@material-ui/icons/Menu';

import Divider from '@material-ui/core/Divider';

import Flag from 'react-flags';

import BottomNavigation from '@material-ui/core/BottomNavigation';
import BottomNavigationAction from '@material-ui/core/BottomNavigationAction';
import { push } from 'connected-react-router';
import Notify from '../Notify';
import SysInfo from '../sysInfo';

import styles from './styles';

import i18n from '../../i18n';
import HealthCheck from '../HealthCheck';
import Button from '../CustomButtons/Button';

import { setNewNotification } from '../../modules/notification/action';

const lodash = require('lodash');

/* eslint-disable react/prefer-stateless-function */

// export default function withLayout(SubCompontents, showTop = true) {
class ConnectedLayout extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      isMenuOpen: false,
      anchorEl: null,
      value: 'welcome',
      showStatus: null,
      showSysInfo: null
    };
    this.toggleMenu = this.toggleMenu.bind(this);
    this.handleMenu = this.handleMenu.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.handleStatus = this.handleStatus.bind(this);
    this.handleCloseStatus = this.handleCloseStatus.bind(this);
    this.handleSysInfo = this.handleSysInfo.bind(this);
    this.handleCloseSysInfo = this.handleCloseSysInfo.bind(this);
  }

  // shouldComponentUpdate(nextProps, nextState) {
  //   return true;
  // }

  handleChange = (event, value) => {
    this.setState({ value });
  };

  toggleMenu(open, e) {
    let shouldProcessing = true;
    const { orderStatus, workMode } = this.props;
    const isAutoMode = workMode === 'auto';
    if (
      lodash.includes(['Ready', 'PreDoing', 'Timeout', 'Init'], orderStatus) ||
      !isAutoMode
    ) {
      shouldProcessing = false;
    }

    if (!shouldProcessing) {
      this.setState({
        isMenuOpen: open
      });
    }
  }

  handleMenu(event) {
    this.setState({ anchorEl: event.currentTarget });
  }

  handleClose() {
    this.setState({ anchorEl: null });
  }

  handleChangeLng(lng) {
    i18n.changeLanguage(lng);
    this.setState({ anchorEl: null });
  }

  handleStatus(event) {
    this.setState({
      showStatus: event.currentTarget
    });
  }

  handleSysInfo(event) {
    this.setState({
      showSysInfo: event.currentTarget
    });
  }

  handleCloseStatus() {
    this.setState({ showStatus: null });
  }

  handleCloseSysInfo() {
    this.setState({ showSysInfo: null });
  }

  HealthCheckOk() {
    const { healthCheckResults } = this.props;
    const results = lodash.filter(healthCheckResults, 'enable');
    return lodash.every(results, ['isHealth', true]);
  }

  handleRouterSwitch = e => {
    console.log(e.target);
  };

  render() {
    let disabled = true;
    const {
      orderStatus,
      classes,
      workMode,
      healthCheckResults,
      usersInfo,
      doPush,
      notification,
      path,
      shouldRender,
      children,
      pages
    } = this.props;
    const isAutoMode = workMode === 'auto';
    const { name, avatar, role } = usersInfo[0];
    if (
      lodash.includes(['Ready', 'PreDoing', 'Timeout', 'Init'], orderStatus)
      // || !isAutoMode
    ) {
      disabled = false;
    }

    const { anchorEl, value, showStatus, isMenuOpen, showSysInfo } = this.state;
    const open = Boolean(anchorEl);

    const openStatusMenu = Boolean(showStatus);

    const openSysInfo = Boolean(showSysInfo);

    const disableSwipeToOpen = false;

    const statusClassName = this.HealthCheckOk()
      ? classes.menuStatusOK
      : classes.menuStatusFail;
    console.log('render');
    return (
      <I18n ns="translations">
        {t => (
          <React.Fragment>
            {/* <ClickAwayListener onClickAway={() => this.toggleMenu(false)}> */}
            {/* <SwipeableDrawer */}
            {/* anchor="right" */}
            {/* open={isMenuOpen} */}
            {/* disableSwipeToOpen={disableSwipeToOpen} */}
            {/* onClose={() => this.toggleMenu(false)} */}
            {/* onOpen={() => this.toggleMenu(true)} */}
            {/* > */}
            {/* <div */}
            {/* tabIndex={0} */}
            {/* role="button" */}
            {/* aria-hidden */}
            {/* onClick={() => this.toggleMenu(false)} */}
            {/* > */}
            {/* <NavBar /> */}
            {/* </div> */}
            {/* </SwipeableDrawer> */}
            {/* </ClickAwayListener> */}
            {/* <SubCompontents /> */}
            <Notify />
            {/*<div style={{height:'calc(100% - 64px)'}}>*/}
            {/*  {children}*/}
            {/*</div>*/}
            <AppBar className={classes.appBar}>
              <Toolbar className={classes.topBar}>
                <div className={classes.menuBtnWrapAvatar}>
                  <img
                    alt={name}
                    src={avatar}
                    className={`${classes.imgRaised} ${
                      classes.imgRoundedCircle
                    } ${classes.imgFluid}`}
                    style={{ height: '100%' }}
                  />
                </div>
                <div className={classes.menuUserName}>
                  <p>{name}</p>
                </div>
                <div className={classes.menuClock}>
                  <Clock
                    className={classes.timeContent}
                    format="HH:mm:ss"
                    ticking
                    timezone="Asia/Shanghai"
                  />
                </div>
                <div className={classes.menuBtnWrapLeft}>
                  <BottomNavigation
                    value={path}
                    // onChange={this.handleChange}
                    showLabels
                    className={classes.BottomNavigation}
                  >
                    {pages.slice(1, -1).map(route => (
                      <BottomNavigationAction
                        key={route.name}
                        value={route.url}
                        // component={Link}
                        // to={route.url}
                        onClick={() => {
                          if(disabled){
                            return;
                          }
                          if (
                            route.roles &&
                            lodash.includes(route.roles, role)
                          ) {
                            doPush(route.url);
                          } else {
                            notification('error', '没有访问权限');
                          }
                        }}
                        label={t(route.title)}
                        icon={<route.icon />}
                        className={classes.BottomNavigationIcon}
                        // disabled={disabled}
                      />
                    ))}
                  </BottomNavigation>
                </div>
                <div className={classes.menuBtnWrapRight}>
                  {/* <Button */}
                  {/* onClick={this.handleSysInfo} */}
                  {/* className={`${statusClassName}`} */}
                  {/* > */}
                  {/* {'系统'} */}
                  {/* </Button> */}
                  <Button
                    onClick={this.handleStatus}
                    className={`${statusClassName}`}
                  >
                    {'连接'}
                  </Button>

                  <IconButton
                    aria-owns={open ? 'menu-appbar' : null}
                    aria-haspopup="true"
                    onClick={this.handleMenu}
                    color="inherit"
                    disabled={disabled}
                  >
                    <Language />
                  </IconButton>
                  {/* <Menu */}
                  {/* id="menu-sysInfo" */}
                  {/* anchorEl={showSysInfo} */}
                  {/* anchorOrigin={{ */}
                  {/* vertical: 'top', */}
                  {/* horizontal: 'left' */}
                  {/* }} */}
                  {/* transformOrigin={{ */}
                  {/* vertical: 'bottom', */}
                  {/* horizontal: 'left' */}
                  {/* }} */}
                  {/* open={openSysInfo} */}
                  {/* onClose={this.handleCloseSysInfo} */}
                  {/* TransitionComponent={Fade} */}
                  {/* classes={{ */}
                  {/* paper: classes.popover */}
                  {/* }} */}
                  {/* > */}
                  {/* <SysInfo /> */}
                  {/* </Menu> */}
                  <Menu
                    id="menu-healthz"
                    anchorEl={showStatus}
                    anchorOrigin={{
                      vertical: 'top',
                      horizontal: 'left'
                    }}
                    transformOrigin={{
                      vertical: 'bottom',
                      horizontal: 'left'
                    }}
                    open={openStatusMenu}
                    onClose={this.handleCloseStatus}
                    TransitionComponent={Fade}
                    classes={{
                      paper: classes.popover
                    }}
                  >
                    <HealthCheck healthCheckResults={healthCheckResults} />
                  </Menu>
                  <Menu
                    id="menu-i18n"
                    anchorEl={anchorEl}
                    anchorOrigin={{
                      vertical: 'top',
                      horizontal: 'left'
                    }}
                    transformOrigin={{
                      vertical: 'bottom',
                      horizontal: 'left'
                    }}
                    open={open}
                    onClose={this.handleClose}
                    TransitionComponent={Fade}
                  >
                    <MenuItem
                      className={classes.menuItem}
                      onClick={() => this.handleChangeLng('en')}
                    >
                      <ListItemIcon className={classes.icon}>
                        <Flag
                          name="GB"
                          format="png"
                          pngSize={24}
                          basePath="./flags"
                        />
                      </ListItemIcon>
                      <ListItemText
                        classes={{ primary: classes.primary }}
                        inset
                        primary={t('Language.en')}
                      />
                    </MenuItem>
                    <Divider />
                    <MenuItem
                      className={classes.menuItem}
                      onClick={() => this.handleChangeLng('zh_CN')}
                    >
                      <ListItemIcon className={classes.icon}>
                        <Flag
                          name="CN"
                          format="png"
                          pngSize={24}
                          basePath="./flags"
                        />
                      </ListItemIcon>
                      <ListItemText
                        classes={{ primary: classes.primary }}
                        inset
                        primary={t('Language.zh_CN')}
                      />
                    </MenuItem>
                  </Menu>
                </div>
              </Toolbar>
            </AppBar>
          </React.Fragment>
        )}
      </I18n>
    );
  }
}

/* eslint-disable react/forbid-prop-types */

ConnectedLayout.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  usersInfo: PropTypes.array.isRequired,
  orderStatus: PropTypes.string.isRequired,
  workMode: PropTypes.string.isRequired,
  healthCheckResults: PropTypes.shape({}).isRequired
};

const mapStateToProps = (state, ownProps) => ({
  usersInfo: state.users,
  orderStatus: state.operations.operationStatus,
  workMode: state.workMode.workMode,
  healthCheckResults: state.healthCheckResults,
  path:state.router.location.pathname,
  ...ownProps
});

const mapDispatchToProps = {
  doPush: push,
  notification: setNewNotification
};

export default withStyles(styles, { withTheme: true })(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(ConnectedLayout)
);
