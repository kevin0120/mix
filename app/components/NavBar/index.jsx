// @flow
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { AppBar } from '@material-ui/core';
import lodash from 'lodash';
import Clock from 'react-live-clock';
import { push } from 'connected-react-router';
import { connect } from 'react-redux';
import { setNewNotification } from '../../modules/notification/action';
import { logoutRequest } from '../../modules/user/action';
import HealthCheck from '../HealthCheck';
import PageEntrance from '../pageEntrance';
import styles from './styles';
import Avatar from '../Avatar';
import SysInfo from '../sysInfo';

import NavBarMenu from './NavBarMenu';
import LanguageMenu from './LanguageMenu';
import type { tRouteObj,tUrl } from '../../containers/model';
import type { Dispatch } from '../../modules/indexReducer';
import type { tUser } from '../../modules/user/model';

type Props = {
  classes: {},
  disabled: boolean,
  contents: Array<string>,
  childRoutes: Array<tRouteObj>,
  healthCheckResults: {},
  users: Array<tUser>,
  doPush: Dispatch,
  notification: Dispatch,
  path: string,
  self: tRouteObj,
  logout: Dispatch,
  getContentByUrl: (tUrl)=>tRouteObj
};

function ConnectedNavBar(
  {
    healthCheckResults,
    disabled,
    users,
    doPush,
    notification,
    path,
    childRoutes,
    self,
    logout,
    contents,
    getContentByUrl
  }: Props) {

  function HealthCheckOk() {
    const results = lodash.filter(healthCheckResults, 'enable');
    return lodash.every(results, ['isHealth', true]);
  }

  const renderSysInfoMenu = (key) =>
    <NavBarMenu
      key={key}
      statusOK={HealthCheckOk()}
      title="系统"
      contents={<SysInfo/>}
    />;

  const renderHealthCheckMenu = (key) =>
    <NavBarMenu
      key={key}
      statusOK={HealthCheckOk()}
      title="连接"
      contents={<HealthCheck healthCheckResults={healthCheckResults}/>
      }
    />;

  const renderLanguageMenu = (key) =>
    <LanguageMenu
      key={key}
      disabled={disabled}
    />;


  const pagesClasses = makeStyles(styles.pages)();
  const renderPageEntrance = (key) =>
    <PageEntrance
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
      navigationClassName={pagesClasses.BottomNavigation}
      ActionClassName={pagesClasses.BottomNavigationIcon}
    />;


  const clockClasses = makeStyles(styles.clock)();
  const renderClock = (key) =>
    <div key={key} className={clockClasses.menuClock}>
      <Clock
        className={clockClasses.timeContent}
        format="HH:mm:ss"
        ticking
        timezone="Asia/Shanghai"
      />
    </div>;


  const avatarClasses = makeStyles(styles.avatar)();
  const renderAvatar = (key) => <Avatar
    key={key}
    className={avatarClasses.menuBtnWrapAvatar}
    users={users}
    onClickAvatar={logout}
  />;


  const renderContentsMapping = {
    sysInfo: renderSysInfoMenu,
    healthCheck: renderHealthCheckMenu,
    language: renderLanguageMenu,
    pages: renderPageEntrance,
    clock: renderClock,
    avatar: renderAvatar
  };

  const classes = makeStyles(styles.root)();
  return (
    <AppBar className={classes.appBar}>
      {contents.map((c) => renderContentsMapping[c](c))}
    </AppBar>
  );
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

export default connect(
  mapState,
  mapDispatch
)(ConnectedNavBar);
