// @flow
import React, { useEffect, useState } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { AppBar } from '@material-ui/core';
import * as lodash from 'lodash-es';
import Clock from 'react-live-clock';
import { push } from 'connected-react-router';
import { connect } from 'react-redux';
import { logoutRequest } from '../../modules/user/action';
import HealthCheck from '../HealthCheck';
import PageEntrance from '../pageEntrance';
import styles from './styles';
import Avatar from '../Avatar';
import SysInfo from '../sysInfo';
import healthzActions from '../../modules/healthz/action';
import NavBarMenu from './NavBarMenu';
import LanguageMenu from './LanguageMenu';
import type { tRouteObj, tUrl } from '../../containers/typeDef';
import type { Dispatch, tAction } from '../../modules/typeDef';
import type { tUser } from '../../modules/user/interface/typeDef';
import notifyActions from '../../modules/Notifier/action';
import workCenterModeActions from '../../modules/workCenterMode/action';
import type { CommonLogLvl } from '../../common/utils';
import Button from '../CustomButtons/Button';
import { translation as trans, navBarNs } from './local';
import Alert from '../Alert';
import { withI18n } from '../../i18n';
import { CommonLog } from '../../common/utils';


type OP = {|
  disabled: boolean,
  contents: Array<string>,
  childRoutes: Array<tRouteObj>,
  self: tRouteObj,
  getContentByUrl: (tUrl)=>tRouteObj
|};

type SP = {|
  ...OP,
  users: Array<tUser>,
  workMode: string,
  workCenterMode: string,
  healthCheckResults: {},
  path: string
|};

type DP = {|
  doPush: Dispatch,
  notification: (variant: CommonLogLvl, message: string, meta: Object)=>tAction<any, any>,
  logout: Dispatch,
  switchWorkCenterMode: Dispatch,
  updateHealthz: Dispatch
|};

type Props = {|
  ...SP,
  ...DP
|};

const mapState = (state, ownProps: OP): SP => ({
  ...ownProps,
  users: state.users,
  // orderStatus: state.operations.operationStatus,
  workMode: state.workMode.workMode,
  workCenterMode: state.workCenterMode,
  healthCheckResults: state.healthz.status || {},
  path: state.router.location.pathname
});

const mapDispatch: DP = {
  doPush: push,
  notification: notifyActions.enqueueSnackbar,
  switchWorkCenterMode: workCenterModeActions.aSwitchWorkCenterMode,
  logout: logoutRequest,
  updateHealthz: healthzActions.update
};

function ConnectedNavBar(
  {
    healthCheckResults,
    disabled,
    users,
    doPush,
    notification,
    path,
    workCenterMode,
    childRoutes,
    self,
    logout,
    contents,
    getContentByUrl,
    updateHealthz,
    switchWorkCenterMode
  }: Props) {
  
  
  const [healthOK, setHealthOK] = useState(false);
  
  useEffect(() => {
    const HealthCheckOk = (): boolean => (!Object.keys(healthCheckResults).some(r => !healthCheckResults[r]));
    setHealthOK(HealthCheckOk());
  }, [healthCheckResults]);
  
  const [showSwitchWorkCenterModeDiag, setShowSwitchWorkCenterModeDiag] = useState(false);
  
  // const [workCenterMode, setWorkCenterMode] = useState(trans.normWorkCenterMode); // 将其翻译直接作为工作模式
  
  const [showButtonColor, setShowButtonColor] = useState('primary');
  useEffect(() => {
    const btnColor = (): string => {
      switch (workCenterMode) {
        case trans.reworkWorkCenterMode:
          return 'danger';
        case trans.normWorkCenterMode:
          return 'primary';
        default:
          return 'primary';
      }
    };
    setShowButtonColor(btnColor());
  }, [workCenterMode]);
  
  
  const renderSysInfoMenu = (key) =>
    <NavBarMenu
      key={key}
      statusOK={healthOK}
      title="系统"
      contents={<SysInfo/>}
    />;
  
  const renderHealthCheckMenu = (key) =>
    <NavBarMenu
      key={key}
      statusOK={healthOK}
      title="连接"
      onClick={() => updateHealthz()}
    >
      <HealthCheck status={healthCheckResults}/>
    </NavBarMenu>;
  
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
          notification('Error', '没有访问权限');
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
  
  const switchWorkCenterModeClasses = makeStyles(styles.switchWorkCenterButton)();
  
  const renderWorkCenterModeToggleButton = (key) =>
    withI18n(
      t => (
        <div>
          <Button
            key={key}
            type="button"
            onClick={() => {
              setShowSwitchWorkCenterModeDiag(!showSwitchWorkCenterModeDiag);
            }}
            variant="contained"
            color={showButtonColor}
            className={switchWorkCenterModeClasses.bigButton}
          >
            {t(workCenterMode)}
          </Button>
          <Alert
            show={showSwitchWorkCenterModeDiag}
            title={t(trans.switchWorkCenterModeTitle)}
            onConfirm={() => {
              switch (workCenterMode) {
                case trans.reworkWorkCenterMode:
                  switchWorkCenterMode(trans.normWorkCenterMode);
                  break;
                case trans.normWorkCenterMode:
                  switchWorkCenterMode(trans.reworkWorkCenterMode);
                  break;
                default:
                  CommonLog.lError(`{workCenterMode} Is Not Support`);
                  break;
              }
              setShowSwitchWorkCenterModeDiag(false);
            }}
            onCancel={() => {
              setShowSwitchWorkCenterModeDiag(false);
            }}
            confirmBtnCssClass={`${switchWorkCenterModeClasses.button} ${
              switchWorkCenterModeClasses.successWarn
              }`}
            cancelBtnCssClass={`${switchWorkCenterModeClasses.button} ${
              switchWorkCenterModeClasses.danger
              }`}
            confirmBtnText={t(trans.confirm)}
            cancelBtnText={t(trans.cancel)}
            content={t(trans.switchWorkCenterModeContent)}
            showCancel
          />
        </div>
      ),
      navBarNs);
  
  
  const renderContentsMapping = {
    sysInfo: renderSysInfoMenu,
    healthCheck: renderHealthCheckMenu,
    language: renderLanguageMenu,
    pages: renderPageEntrance,
    clock: renderClock,
    avatar: renderAvatar,
    switchWorkCenterButton: renderWorkCenterModeToggleButton
  };
  
  const classes = makeStyles(styles.root)();
  return (
    <AppBar className={classes.appBar}>
      {contents.map((c) => renderContentsMapping[c](c))}
    </AppBar>
  );
}

export default connect<Props, OP, SP, DP, _, _>(
  mapState,
  mapDispatch
)(ConnectedNavBar);
