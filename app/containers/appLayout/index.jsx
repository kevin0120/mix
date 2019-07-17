// @flow
import React from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Button from '@material-ui/core/Button';
import Avatar from '@material-ui/core/Avatar';
import Notify from '../../components/Notify';
import { setNewNotification } from '../../modules/notification/action';
import NavBar from '../../components/NavBar';
import LayoutDrawer from '../../components/LayoutDrawer';

type Props = {
  users: [],
  path: string,
  children: {},
  childRoutes: [],
  self: {}
};

function AppLayout(
  {
    users,
    path,
    children,
    childRoutes,
    self
  }: Props) {
  const { DefaultContent } = self;
  return (
    <React.Fragment>
      <Notify/>
      <div style={{ height: 'calc(100% - 64px)', display: 'flex' }}>
        <LayoutDrawer
          contents={users.map((u) => ({
            icon: u.avatar ? <Avatar src={u.avatar}/> : <Avatar>{u.name.slice(2)}</Avatar>,
            label: (<div style={{
              display: 'flex',
              flexDirection: 'row',
              justifyContent: 'space-between',
              alignItems: 'center'
            }}>
              {u.name}
              <Button color="secondary" size="small" variant="contained">
                Logout
              </Button>
            </div>)
          }))}
        />
        {path === '/app' ? <DefaultContent childRoutes={childRoutes}/> : children}
      </div>
      <NavBar
        contents={self.navBarContents}
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
  doPush: push,
  notification: setNewNotification
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AppLayout);
