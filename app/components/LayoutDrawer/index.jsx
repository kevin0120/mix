import React from 'react';
import { makeStyles, createMuiTheme, MuiThemeProvider } from '@material-ui/core/styles';

import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import clsx from 'clsx';

import { ListItemAvatar } from '@material-ui/core';
import styles from './styles';

export default function LayoutDrawer(props) {
  const { contents } = props;
  const classes = makeStyles(styles)();
  const theme = createMuiTheme({
    overrides: {
      MuiDrawer: {
        root: {
          zIndex: '0'
        },
        paper: {
          zIndex: '0',
          backgroundImage: 'url("../resources/imgs/texture.png")',
          backgroundRepeat: 'repeat',
          backgroundColor: '#444',
          color:'#ddd'
        }
      }
    }
  });

  const [open, setOpen] = React.useState(false);

  function toggleDrawer() {
    setOpen(!open);
  }

  return (
    <MuiThemeProvider theme={theme}>
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
          {contents.map((C) => (
            <ListItem key={`${C.label}`}>
              <ListItemAvatar>
                {C.icon}
              </ListItemAvatar>
              <ListItemText>
                {C.label}
              </ListItemText>
            </ListItem>
          ))}
        </List>
        <Divider/>
      </Drawer>
    </MuiThemeProvider>
  );
}
