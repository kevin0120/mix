import React, { useState } from 'react';
import PropTypes from 'prop-types';
import MenuList from '@material-ui/core/MenuList';
import MenuItem from '@material-ui/core/MenuItem';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import ViewModuleIcon from '@material-ui/icons/ViewModule';
import SettingsRemoteIcon from '@material-ui/icons/SettingsRemote';
import makeStyles from '@material-ui/core/styles/makeStyles';
import i18n from '../../i18n';
import LeftMenuWithAvatar from '../../components/LeftMenuWithAvatar';
import Net from '../../components/PreferencesContent/Net';
import Io from '../../components/PreferencesContent/Io';
import Connect from '../../components/PreferencesContent/Connect';

import styles from './styles';

const menuContents = [
  {
    text: 'Configuration.network.name',
    icon: selected => (
      <RssFeedIcon
        style={{ fontSize: 45, margin: 5, fill: selected ? '#FAFAFA' : '#009688' }}
      />
    ),
    component: <Net/>
  },
  {
    text: 'Configuration.IO.name',
    icon: selected => (
      <ViewModuleIcon
        style={{ fontSize: 45, margin: 5, fill: selected ? '#FAFAFA' : '#ff9800' }}
      />
    ),
    component: <Io/>
  },
  {
    text: 'Configuration.connections.name',
    icon: selected => (
      <SettingsRemoteIcon
        style={{ fontSize: 45, margin: 5, fill: selected ? '#FAFAFA' : '#3492ff' }}
      />
    ),
    component: <Connect/>
  }
];

function ConnectedPreferences() {
  const [activeMenu, setMenu] = useState('Configuration.network.name');
  const classes = makeStyles(styles)();

  const currentComponentIdx = menuContents.findIndex(
    ele => activeMenu === ele.text
  );

  const menuItems = () => menuContents.map(item => (
    <MenuItem
      selected={activeMenu === item.text}
      key={item.text}
      className={classes.menuItem}
      component={React.forwardRef((p, r) => (
        <Card
          onClick={() => setMenu(item.text)}
          className={
            activeMenu === item.text
              ? classes.menuItemSelected
              : classes.menuItem
          }
          ref={r}
        >
          <CardActionArea
            classes={{
              root:
                activeMenu === item.text
                  ? classes.cardActionAreaSelected
                  : classes.cardActionArea
            }}
          >
            {item.icon(activeMenu === item.text)}
            <span className={classes.itemText}>{i18n.t(item.text)}</span>
          </CardActionArea>
        </Card>
      ))}
    />
  ));

  return (
    <div className={classes.root}>
      {/*<AppBarBack />*/}
      <LeftMenuWithAvatar>
        <MenuList>{menuItems()}</MenuList>
      </LeftMenuWithAvatar>
      <div className={classes.content}>
        {menuContents[currentComponentIdx].component}
      </div>
    </div>
  );
}

ConnectedPreferences.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

export default ConnectedPreferences;
