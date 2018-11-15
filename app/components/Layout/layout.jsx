// @flow
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
/* eslint-disable no-unused-vars */
import SwipeableDrawer from '@material-ui/core/SwipeableDrawer';
import Slide from '@material-ui/core/Slide';
import { I18n, Trans } from 'react-i18next';

import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import ListItemIcon from '@material-ui/core/ListItemIcon'
import ListItemText from '@material-ui/core/ListItemText'
import ClickAwayListener from '@material-ui/core/ClickAwayListener';

import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';

import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import Language from '@material-ui/icons/LanguageRounded';
import Fade from '@material-ui/core/Fade';

import MenuIcon from '@material-ui/icons/Menu';


import NavBar from '../../components/NavBar';
import Notify from '../../components/Notify';

import Divider from '@material-ui/core/Divider';

import Flag from "react-flags";



import BottomNavigation from '@material-ui/core/BottomNavigation';
import BottomNavigationAction from '@material-ui/core/BottomNavigationAction';


import styles from './styles';

import { routeConfigs } from '../../routes';

import i18n from '../../i18n';
import HealthCheck from '../HealthCheck';
import Button from '../CustomButtons/Button';

/* eslint-disable react/prefer-stateless-function */
export default function withLayout(SubCompontents, showTop = true) {
  class ConnectedLayout extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        isMenuOpen: false,
        anchorEl: null,
        value: "recents",
        showStatus: null,
      };
      this.toggleMenu = this.toggleMenu.bind(this);
      this.handleMenu = this.handleMenu.bind(this);
      this.handleClose = this.handleClose.bind(this);
      this.handleStatus = this.handleStatus.bind(this);
      this.handleCloseStatus = this.handleCloseStatus.bind(this);
    }

    handleChange = (event, value) => {
      this.setState({ value });
    };

    toggleMenu(open,e) {
      let shouldProcessing = true;
      if (this.props.orderStatus === 'Ready' || this.props.orderStatus === 'PreDoing' || this.props.orderStatus === 'Timeout' || this.props.orderStatus === 'Init' || (!this.props.isAutoMode)){
        shouldProcessing = false;
      }

      if (!shouldProcessing) {
        this.setState({
          isMenuOpen: open,
        });
      }
    }

    shouldComponentUpdate(nextProps, nextState) {
      return true;
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

    handleCloseStatus() {
      this.setState({ showStatus: null });
    }

    HealthCheckOk() {
      const { healthCheckResults } = this.props;
      for (const key in healthCheckResults) {
        if (!healthCheckResults[key].health) {
          return false;
        }
      }

      return false;
    }

    render() {

      let shouldProcessing = true;
      if (this.props.orderStatus === 'Ready' || this.props.orderStatus === 'PreDoing' || this.props.orderStatus === 'Timeout' || this.props.orderStatus === 'Init' || (!this.props.isAutoMode)){
        shouldProcessing = false;
      }

      const { classes, theme , avatarImg, Username} = this.props;
      const { anchorEl, value, showStatus } = this.state;
      const open = Boolean(anchorEl);

      const openStatusMenu = Boolean(showStatus);

      const disableSwipeToOpen = false;

      const statusClassName = this.HealthCheckOk() ? classes.menuStatusOK : classes.menuStatusFail;

      return (
        <I18n ns="translations">
          {
            t => (
              <div className={classes.layout}>
                <ClickAwayListener onClickAway={() => this.toggleMenu(false)}>
                  <SwipeableDrawer
                    anchor="right"
                    open={this.state.isMenuOpen}
                    disableSwipeToOpen={disableSwipeToOpen}
                    onClose={() => this.toggleMenu(false)}
                    onOpen={() => this.toggleMenu(true)}
                  >
                    <div
                      tabIndex={0}
                      role="button"
                      aria-hidden
                      onClick={() => this.toggleMenu(false)}
                    >
                      <NavBar />
                    </div>
                  </SwipeableDrawer>
                </ClickAwayListener>
                <SubCompontents />
                <Notify />
                  <AppBar position='fixed' className={classes.appBar}>
                    <Toolbar className={classes.topBar}>
                      {/* <Typography variant="title" color="inherit" className={classes.navTitle}>
                        {t('Title')}
                      </Typography> */}
                      <div className={classes.menuBtnWrapAvatar}>
                        {/*<IconButton*/}
                            {/*aria-owns={open ? 'menu-appbar' : null}*/}
                            {/*aria-haspopup="true"*/}
                            {/*onClick={this.handleMenu}*/}
                            {/*color="inherit"*/}
                            {/*disabled={ shouldProcessing }*/}
                          {/*>*/}
                            {/*/!*<Language />*!/*/}
                          {/*</IconButton>*/}
                        <img
                          alt={name}
                          src={avatarImg}
                          className={ classes.imgRaised +
                          " " +
                          classes.imgRounded +
                          " " +
                          classes.imgFluid
                          }
                        />

                      </div>
                      <div className={classes.menuUserName}>
                        <p >
                          {Username}
                        </p>
                      </div>
                      <div className={classes.menuBtnWrapLeft}>
                        <BottomNavigation
                          value={value}
                          onChange={this.handleChange}
                          showLabels
                          className={classes.BottomNavigation}
                        >
                          {
                            routeConfigs.slice(0,-1).map(route => (
                            <BottomNavigationAction label={t(route.title)} component="a" href={`#${route.url}`} icon={<route.icon />} className={classes.BottomNavigationIcon} disabled={ shouldProcessing }/>
                            )
                            )}
                        </BottomNavigation>
                      </div>
                      <div className={classes.menuBtnWrapRight}>
                        <Button onClick={this.handleStatus} className={`${statusClassName}`}>
                          {'连接信息'}
                        </Button>

                        <IconButton
                          aria-owns={open ? 'menu-appbar' : null}
                          aria-haspopup="true"
                          onClick={this.handleMenu}
                          color="inherit"
                          disabled={ shouldProcessing }
                        >
                          <Language />
                        </IconButton>
                        {/*<IconButton color="inherit" aria-label="Menu" className={classes.menuButton} onClick={(e) => this.toggleMenu(true,e)} disabled={ shouldProcessing }>*/}
                          {/*<MenuIcon />*/}
                        {/*</IconButton>*/}
                        <Menu
                          id="menu-appbar"
                          anchorEl={showStatus}
                          anchorOrigin={{
                            vertical: 'top',
                            horizontal: 'right',
                          }}
                          transformOrigin={{
                            vertical: 'top',
                            horizontal: 'right',
                          }}
                          open={openStatusMenu}
                          onClose={this.handleCloseStatus}
                          TransitionComponent={Fade}
                        >
                          <HealthCheck />
                        </Menu>
                        <Menu
                          id="menu-appbar"
                          anchorEl={anchorEl}
                          anchorOrigin={{
                            vertical: 'top',
                            horizontal: 'right',
                          }}
                          transformOrigin={{
                            vertical: 'top',
                            horizontal: 'right',
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
                            <ListItemText classes={{ primary: classes.primary }} inset primary={t('Language.en')}  />
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
                            <ListItemText classes={{ primary: classes.primary }} inset primary={t('Language.zh_CN')}  />
                          </MenuItem>
                        </Menu>
                      </div>
                    </Toolbar>
                  </AppBar>
              </div>
            )
          }
        </I18n>
      );
    }
  }
  /* eslint-disable react/forbid-prop-types */

  ConnectedLayout.propTypes = {
    classes: PropTypes.object.isRequired,
    orderStatus: PropTypes.string.isRequired,
    ControlMode: PropTypes.string,
    isAutoMode: PropTypes.bool,
    healthCheckResults: PropTypes.shape({}).isRequired,
  };

  const mapStateToProps = (state, ownProps) => ({
    debugInfo: state.debugInfo,
    enableDebugInfo: state.userConfigs.enableDebugInfo,
    avatarImg: state.userInfo.image_small,
    theme: PropTypes.object.isRequired,
    orderStatus: state.orderProgress.orderStatus,
    ControlMode: state.ControlMode,
    isAutoMode: state.isAutoMode,
    Username: state.userInfo.name,
    healthCheckResults: state.healthCheckResults,
    ...ownProps,
  });

  return withStyles(styles, { withTheme: true })(connect(mapStateToProps)(ConnectedLayout));
}
