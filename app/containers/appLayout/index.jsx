// @flow
import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
/* eslint-disable no-unused-vars */
import { I18n, Trans } from 'react-i18next';

import { push } from 'connected-react-router';
import Notify from '../../components/Notify';

import styles from './styles';

import i18n from '../../i18n';
import Avatar from '@material-ui/core/Avatar';

import { setNewNotification } from '../../modules/notification/action';
import { logoutRequest } from '../../modules/user/action';
import NavBar from '../../components/NavBar';

import LayoutDrawer from '../../components/LayoutDrawer';
import {Button} from '@material-ui/core'
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
              <LayoutDrawer
                contents={users.map((u) => {
                  return {
                    icon: u.avatar ? <Avatar src={u.avatar}/> : <Avatar>{u.name.slice(2)}</Avatar>,
                    label: (<div style={{display:'flex',flexDirection:'row',justifyContent:'space-between',alignItems:'center'}}>
                        {u.name}
                        <Button color="secondary" size="small" variant="contained">
                          Logout
                        </Button>
                    </div>)
                  };
                })}
              />
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
