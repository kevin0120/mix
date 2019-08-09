import React, { useState } from 'react';
import Menu from '@material-ui/core/Menu';
import Fade from '@material-ui/core/Fade';
import { createMuiTheme, makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import Button from '../CustomButtons/Button';
import styles from './styles';

function NavBarMenu({ statusOK, title, contents }) {
  const classes = makeStyles(styles.NavBarMenu)();
  const [showMenu, setShowMenu] = useState(null);
  const statusClassName = statusOK
    ? classes.menuStatusOK
    : classes.menuStatusFail;
  const open = Boolean(showMenu);

  return (
    <React.Fragment>
        <Button
          variant="contained"
          onClick={e => setShowMenu(e.currentTarget)}
          className={statusClassName}
        >
          {title}
        </Button>
        <Menu
          id="menu-sysInfo"
          anchorEl={showMenu}
          anchorOrigin={{
            vertical: 'top',
            horizontal: 'left'
          }}
          transformOrigin={{
            vertical: 'bottom',
            horizontal: 'left'
          }}
          open={open}
          onClose={() => setShowMenu(null)}
          TransitionComponent={Fade}
          classes={{
            paper: classes.popover
          }}
        >
          {contents}
        </Menu>
    </React.Fragment>
  );
}

export default NavBarMenu;
