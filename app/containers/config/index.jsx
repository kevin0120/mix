import React from 'react';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import MenuList from '@material-ui/core/MenuList';
import MenuItem from '@material-ui/core/MenuItem';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import ViewModuleIcon from '@material-ui/icons/ViewModule';
import SettingsRemoteIcon from '@material-ui/icons/SettingsRemote';
import { I18n } from 'react-i18next';

import withLayout from '../../components/Layout/layout';
import LeftMenuWithAvatar from '../../components/LeftMenuWithAvatar';
import AppBarBack from '../../components/AppBarBack';
import Net from '../../components/PreferencesContent/Net';
import Io from '../../components/PreferencesContent/Io';
import Connect from '../../components/PreferencesContent/Connect';

import styles from './styles';

const menuContents = [
  {
    text: 'Configuration.network.name',
    icon: <RssFeedIcon style={{ fill: '#009688' }} />,
    component: <Net />
  },
  {
    text: 'Configuration.IO.name',
    icon: <ViewModuleIcon style={{ fill: '#ff9800' }} />,
    component: <Io />
  },
  {
    text: 'Configuration.connections.name',
    icon: <SettingsRemoteIcon style={{ fill: '#3492ff' }} />,
    component: <Connect />
  }
];

class ConnectedPreferences extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      activeMenu: 'Configuration.network.name'
    };

    this.handleChangeMenu = this.handleChangeMenu.bind(this);
  }

  componentDidMount() {
    // const { odooUrl, hmiSn } = this.props;
    // this.props.initHmiConnInfo(odooUrl, hmiSn, false);
  }

  handleChangeMenu(text) {
    this.setState({
      activeMenu: text
    });
  }

  render() {
    const { classes } = this.props;
    const { activeMenu } = this.state;
    const currentComponentIdx = menuContents.findIndex(
      ele => activeMenu === ele.text
    );

    const menuList = t => {
      const menus = menuContents.map(item => (
        <MenuItem
          selected={activeMenu === item.text}
          key={item.text}
          className={classes.menuItem}
          onClick={() => this.handleChangeMenu(item.text)}
        >
          <ListItemIcon>{item.icon}</ListItemIcon>
          <span className={classes.itemText}>{t(item.text)}</span>
        </MenuItem>
      ));
      return <div >{menus}</div>;
    };

    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.root}>
            <AppBarBack />
            <LeftMenuWithAvatar>
              <MenuList>{menuList(t)}</MenuList>
            </LeftMenuWithAvatar>
            <div className={classes.content}>
              {/*<div className={classes.toolbar} />*/}
              {menuContents[currentComponentIdx].component}
            </div>
          </div>
        )}
      </I18n>
    );
  }
}

ConnectedPreferences.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

export default withLayout(withStyles(styles)(ConnectedPreferences), false);
