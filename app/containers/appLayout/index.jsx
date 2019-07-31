// @flow
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Button from '@material-ui/core/Button';
import Avatar from '@material-ui/core/Avatar';
import Notify from '../../components/Notify';
import { setNewNotification } from '../../modules/notification/action';
import NavBar from '../../components/NavBar';
import LayoutDrawer from '../../components/LayoutDrawer';
import type { tUser } from '../../modules/user/model';
import { logoutRequest } from '../../modules/user/action';
import type { Dispatch } from '../../modules/indexReducer';

type Props = {
  users: Array<tUser>,
  path: string,
  children: PropTypes.element,
  childRoutes: [],
  self: {},
  logout: Dispatch
};

function AppLayout(
  {
    users,
    path,
    children,
    childRoutes,
    self,
    logout
  }: Props) {
  const { DefaultContent, navBarContents } = self;
  return (
    <React.Fragment>
      <Notify/>
      <div style={{ height: 'calc(100% - 64px)', display: 'flex' }}>
        <LayoutDrawer
          contents={users.map((u) => ({
            key: u.name,
            icon: u.avatar ? <Avatar src={u.avatar}/> : <Avatar>{u.name.slice(0, 2)}</Avatar>,
            label: (<div style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              {u.name}
              <Button
                color="secondary"
                size="small"
                variant="contained"
                onClick={() => logout(u.uuid)}
              >
                Logout
              </Button>
            </div>)
          }))}
        />
        {path === '/app' ? <DefaultContent childRoutes={childRoutes}/> : children}
      </div>
      <NavBar
        contents={navBarContents}
        self={self}
        childRoutes={childRoutes}
      />
    </React.Fragment>
  );
}

const mapStateToProps = (state, ownProps) => ({
  users: state.users,
  healthCheckResults: state.healthCheckResults,
  path: state.router.location.pathname,
  ...ownProps
});

const mapDispatchToProps = {
  logout: logoutRequest,
  doPush: push,
  notification: setNewNotification
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AppLayout);
