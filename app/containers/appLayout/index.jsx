// @flow
import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Button from '../../components/CustomButtons/Button';
import Avatar from '@material-ui/core/Avatar';
import Typography from '@material-ui/core/Typography';
import Notify from '../../components/Notify';
import NavBar from '../../components/NavBar';
import LayoutDrawer from '../../components/LayoutDrawer';
import type { tUser } from '../../modules/user/model';
import { logoutRequest } from '../../modules/user/action';
import type { Dispatch } from '../../modules/indexReducer';
import type { tRouteComponent, tRouteObj } from '../model';
import Notifier from '../../components/Notifier';

type Props = {
  users: Array<tUser>,
  path: string,
  children: Array<tRouteComponent>,
  childRoutes: Array<tRouteObj>,
  self: tRouteObj & {
    DefaultContent: PropTypes.element,
    navBarContents: Array<string>
  },
  logout: Dispatch,
  getContentByUrl: (string)=>tRouteObj
};

function AppLayout(
  {
    users,
    path,
    children,
    childRoutes,
    self,
    logout,
    getContentByUrl
  }: Props) {
  const { DefaultContent, navBarContents } = self;
  return (
    <React.Fragment>
      <Notify/>
      <Notifier/>
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
                color="warning"
                size="md"
                regular
                variant="contained"
                onClick={() => logout(u.uuid)}
              >
                <Typography variant="body1">
                  Logout
                </Typography>
              </Button>
            </div>)
          }))}
        />
        {path === '/app' ? <DefaultContent childRoutes={childRoutes}/> : children}
      </div>
      <NavBar
        getContentByUrl={getContentByUrl}
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
  doPush: push
};

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(AppLayout);
