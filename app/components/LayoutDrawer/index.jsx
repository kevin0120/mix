// @flow
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import clsx from 'clsx';
import type { Node } from 'react';

import { ListItemAvatar } from '@material-ui/core';
import styles from './styles';

type tContent = {
  key: string,
  label: Node,
  icon: Node
};

type tProps = {
  contents: Array<tContent>
};

export default function LayoutDrawer({ contents }: tProps) {
  const classes = makeStyles(styles)();

  const [open, setOpen] = React.useState(false);

  function toggleDrawer() {
    setOpen(!open);
  }

  return (
    <Drawer
      variant="permanent"
      className={clsx(classes.drawer, {
        [classes.drawerOpen]: open,
        [classes.drawerClose]: !open
      })}
      classes={{
        paper: clsx({
          [classes.drawerOpen]: open,
          [classes.drawerClose]: !open
        })
      }}
      open={open}
    >
      <div className={classes.toolbar}>
        <IconButton
          color="inherit"
          aria-label="Open drawer"
          onClick={toggleDrawer}
          className={classes.menuButton}
        >
          <MenuIcon/>
        </IconButton>
      </div>
      <Divider/>
      <List>
        {contents.map((c) => (
          <ListItem key={c.key}>
            <ListItemAvatar>
              {c.icon}
            </ListItemAvatar>
            <ListItemText>
              {c.label}
            </ListItemText>
          </ListItem>
        ))}
      </List>
      <Divider/>
    </Drawer>
  );
}
