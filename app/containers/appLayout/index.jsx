// @flow
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
/* eslint-disable no-unused-vars */
import { I18n, Trans } from 'react-i18next';

import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import Language from '@material-ui/icons/LanguageRounded';
import Fade from '@material-ui/core/Fade';
import Divider from '@material-ui/core/Divider';

import Flag from 'react-flags';
import Clock from 'react-live-clock';

import { push } from 'connected-react-router';
import Notify from '../../components/Notify';
import SysInfo from '../../components/sysInfo';
import HomePage from '../home';

import styles from './styles';

import i18n from '../../i18n';
import HealthCheck from '../../components/HealthCheck';
import Button from '../../components/CustomButtons/Button';

import { setNewNotification } from '../../modules/notification/action';
import HomeOperationList from '../HomeOperationList';
import Avatar from '../../components/Avatar';
import { logoutRequest } from '../../modules/user/action';
import PageEntrance from '../../components/pageEntrance';
import NavBar from '../../components/NavBar';

import LayoutDrawer from '../../components/LayoutDrawer';

const lodash = require('lodash');

/* eslint-disable react/prefer-stateless-function */

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
      users,
      doPush,
      notification,
      path,
      children,
      childRoutes,
      self,
      logout
    } = this.props;
    const { DefaultContent } = self;
    const isAutoMode = workMode === 'auto';
    // const { name, avatar, role } = users[0];
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
            {/* aria-hidden *}
            {/* onClick={() => this.toggleMenu(false)} */}
            {/* > */}
            {/*<NavBar/>*/}
            {/* </div> */}
            {/* </SwipeableDrawer> */}
            {/* </ClickAwayListener> */}
            {/* <SubCompontents /> */}

            <Notify/>
            <div style={{ height: 'calc(100% - 64px)', display: 'flex' }}>
              <LayoutDrawer/>
              {path === '/app' ? <DefaultContent childRoutes={childRoutes}/> : children}
            </div>
            <NavBar
              contents={self.navBarContents}
              self={self}
              childRoutes={childRoutes}
            />
          </React.Fragment>
        )}
      </I18n>
    );
  }
}

/* eslint-disable react/forbid-prop-types */

ConnectedLayout.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  users: PropTypes.array.isRequired,
  orderStatus: PropTypes.string.isRequired,
  workMode: PropTypes.string.isRequired,
  healthCheckResults: PropTypes.shape({}).isRequired
};

const mapStateToProps = (state, ownProps) => ({
  users: state.users,
  orderStatus: state.operations.operationStatus,
  workMode: state.workMode.workMode,
  healthCheckResults: state.healthCheckResults,
  path: state.router.location.pathname,
  ...ownProps
});

const mapDispatchToProps = {
  doPush: push,
  notification: setNewNotification,
  logout: logoutRequest
};

export default withStyles(styles, { withTheme: true })(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(ConnectedLayout)
);
