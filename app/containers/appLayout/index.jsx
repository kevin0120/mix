// @flow
import React from 'react';
import type { Node } from 'react';
import { connect } from 'react-redux';
import { push } from 'connected-react-router';
import Avatar from '@material-ui/core/Avatar';
import Typography from '@material-ui/core/Typography';
import { I18n } from 'react-i18next';
import Button from '../../components/CustomButtons/Button';
import NavBar from '../../components/NavBar';
import LayoutDrawer from '../../components/LayoutDrawer';
import type { tUser } from '../../modules/user/interface/typeDef';
import { logoutRequest } from '../../modules/user/action';
import type { Dispatch } from '../../modules/typeDef';
import type { tRouteComponent, tRouteObj } from '../typeDef';
import Notifier from '../../components/Notifier';
import Dialog from '../../components/Dialog';
import Loading from '../../components/Loading';

type tOP = {|
  children: Array<tRouteComponent>,
  childRoutes: Array<tRouteObj>,
  self: tRouteObj & {
    DefaultContent: Node,
    navBarContents: Array<string>
  },
  getContentByUrl: (string)=>tRouteObj
|};
type tSP = {|
  ...tOP,
  users: Array<tUser>,
  path: string,
|};
type tDP = {|
  logout: Dispatch,
  doPush: Dispatch
|};

type Props = {|
  ...tOP,
  ...tSP,
  ...tDP,

|};

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
    <I18n ns="translations">
      {t => (
        <React.Fragment>
          <Notifier/>
          <Dialog/>
          <Loading/>
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
                      {t('Common.Logout')}
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
      )}
    </I18n>
  );
}

const mapStateToProps = (state, ownProps: tOP): tSP => ({
  ...ownProps,
  users: state.users,
  path: state.router.location.pathname
});

const mapDispatchToProps: tDP = {
  logout: logoutRequest,
  doPush: push
};

export default connect<Props, tOP, tSP, tDP, _, _>(
  mapStateToProps,
  mapDispatchToProps
)(AppLayout);
