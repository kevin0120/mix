// @flow
import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import { I18n } from 'react-i18next';
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Divider from '@material-ui/core/Divider';
import lodash from 'lodash';

import Clock from 'react-live-clock';
import IconButton from '@material-ui/core/IconButton';
import Language from '@material-ui/icons/LanguageRounded';
import Menu from '@material-ui/core/Menu';
import Fade from '@material-ui/core/Fade';
import MenuItem from '@material-ui/core/MenuItem';
import Flag from 'react-flags';
import { push } from 'connected-react-router';
import { connect } from 'react-redux';
import { setNewNotification } from '../../modules/notification/action';
import { logoutRequest } from '../../modules/user/action';
import HealthCheck from '../HealthCheck';
import Button from '../CustomButtons/Button';
import PageEntrance from '../pageEntrance';
import styles from './styles';
import Avatar from '../Avatar';
import SysInfo from '../sysInfo';

import { getContentByUrl } from '../../containers/pages';

type Props = {
  classes: {},
  disabled: boolean,
  contents: Array<string>,
  childRoutes: Array
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedNavBar extends React.Component<Props> {
  constructor(props) {
    super(props);
    this.state = {
      anchorEl: null,
      showStatus: null,
      showSysInfo: null
    };
  }

  renderSysInfoMenu = (key) => {
    const {
      classes
    } = this.props;
    const statusClassName = this.HealthCheckOk()
      ? classes.menuStatusOK
      : classes.menuStatusFail;
    const { showSysInfo } = this.state;
    const openSysInfo = Boolean(showSysInfo);
    return <React.Fragment key={key}>
      <Button

        onClick={this.handleSysInfo}
        className={`${statusClassName}`}
      >
        {'系统'}
      </Button>
      <Menu
        id="menu-sysInfo"
        anchorEl={showSysInfo}
        anchorOrigin={{
          vertical: 'top',
          horizontal: 'left'
        }}
        transformOrigin={{
          vertical: 'bottom',
          horizontal: 'left'
        }}
        open={openSysInfo}
        onClose={this.handleCloseSysInfo}
        TransitionComponent={Fade}
        classes={{
          paper: classes.popover
        }}
      >
        <SysInfo/>
      </Menu>
    </React.Fragment>;
  };

  renderHealthCheckMenu = (key) => {
    const {
      classes,
      healthCheckResults
    } = this.props;
    const { showStatus } = this.state;
    const openStatusMenu = Boolean(showStatus);
    return <React.Fragment key={key}>
      <Button

        onClick={this.handleStatus}
        className={`${statusClassName}`}
      >
        {'连接'}
      </Button>
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
        <HealthCheck healthCheckResults={healthCheckResults}/>
      </Menu>
    </React.Fragment>;
  };

  renderLanguageMenu = (key) => {
    const { classes, disabled } = this.props;
    const { anchorEl } = this.state;
    const open = Boolean(anchorEl);
    return <I18n ns="translations" key={key}>
      {t => (
        <React.Fragment key={key}>
          <IconButton
            aria-owns={open ? 'menu-appbar' : null}
            aria-haspopup="true"
            onClick={this.handleMenu}
            color="inherit"
            disabled={disabled}
          >
            <Language/>
          </IconButton>
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
            <Divider/>
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
        </React.Fragment>
      )}
    </I18n>;
  };

  renderPageEntrance = (key) => {
    const {
      classes,
      users,
      doPush,
      notification,
      path,
      childRoutes,
      self,
      disabled
    } = this.props;
    return <PageEntrance
      key={key}
      type="navigation"
      value={path}
      routes={[self, ...childRoutes, getContentByUrl('/pages/login')]}
      onItemClick={(route) => {
        if (disabled) {
          return;
        }
        if (!route.role || route.role.length === 0 || users.some((u) => lodash.includes(route.role, u.role))) {
          doPush(route.url);
        } else {
          notification('error', '没有访问权限');
        }
      }}
      navigationClassName={classes.BottomNavigation}
      ActionClassName={classes.BottomNavigationIcon}
    />;
  };

  renderClock = (key) => {
    const { classes } = this.props;
    return <div key={key} className={classes.menuClock}>
      <Clock
        className={classes.timeContent}
        format="HH:mm:ss"
        ticking
        timezone="Asia/Shanghai"
      />
    </div>;
  };

  renderAvatar = (key) => {
    const { classes, users, logout } = this.props;
    return <Avatar
      key={key}
      className={classes.menuBtnWrapAvatar}
      users={users}
      onClickAvatar={logout}
    />;
  };


  renderContentsMapping = {
    sysInfo: this.renderSysInfoMenu,
    healthCheck: this.renderHealthCheckMenu,
    language: this.renderLanguageMenu,
    pages: this.renderPageEntrance,
    clock: this.renderClock,
    avatar: this.renderAvatar
  };

  render() {
    const {
      classes,
      contents
    } = this.props;

    return (
      <div className={classes.appBar}>
        {contents.map((c) => this.renderContentsMapping[c](c))}
      </div>
    );
  }
}


const mapState = (state, ownProps) => ({
  users: state.users,
  orderStatus: state.operations.operationStatus,
  workMode: state.workMode.workMode,
  healthCheckResults: state.healthCheckResults,
  path: state.router.location.pathname,
  ...ownProps
});

const mapDispatch = {
  doPush: push,
  notification: setNewNotification,
  logout: logoutRequest
};

export default withStyles(styles, { withTheme: true })(
  connect(
    mapState,
    mapDispatch
  )(ConnectedNavBar)
);
