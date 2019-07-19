import React, { useState } from 'react';
import Menu from '@material-ui/core/Menu';
import Fade from '@material-ui/core/Fade';
import { makeStyles } from '@material-ui/core/styles';
import IconButton from '@material-ui/core/IconButton';
import Language from '@material-ui/icons/LanguageRounded';
import MenuItem from '@material-ui/core/MenuItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Flag from 'react-flags';
import ListItemText from '@material-ui/core/ListItemText';
import Divider from '@material-ui/core/Divider';
import i18n from '../../i18n';
import styles from './styles';
import { I18n } from 'react-i18next';

type Props = {
  disabled: boolean
};

function LanguageMenu({ disabled }: Props) {
  const [anchorEl, setAnchorEl] = useState(null);
  const classes = makeStyles(styles.languageMenu)();
  const open = Boolean(anchorEl);

  function handleChangeLng(lng) {
    i18n.changeLanguage(lng);
    setAnchorEl(null);
  }

  return <I18n ns="translations">
    {t => (
      <React.Fragment>
        <IconButton
          aria-owns={open ? 'menu-appbar' : null}
          aria-haspopup="true"
          onClick={e => {
            setAnchorEl(e.currentTarget);
          }}
          color="inherit"
          disabled={disabled}
        >
          <Language/>
        </IconButton>
        <Menu
          id="menu-i18n"
          anchorEl={anchorEl}
          anchorOrigin={{
            vertical: 'top',
            horizontal: 'left'
          }}
          transformOrigin={{
            vertical: 'bottom',
            horizontal: 'left'
          }}
          open={open}
          onClose={() => setAnchorEl(null)}
          TransitionComponent={Fade}
        >
          <MenuItem
            button
            className={classes.menuItem}
            onClick={() => handleChangeLng('en')}
          >
            <ListItemIcon className={classes.icon}>
              <Flag name="GB" format="png" pngSize={24} basePath="./flags"/>
            </ListItemIcon>
            <ListItemText
              inset
              primary={t('Language.en')}
            />
          </MenuItem>
          <Divider/>
          <MenuItem
            button
            className={classes.menuItem}
            onClick={() => handleChangeLng('zh_CN')}
          >
            <ListItemIcon className={classes.icon}>
              <Flag
                name="CN"
                format="png"
                pngSize={24}
                basePath="./flags"
              />
            </ListItemIcon>
            <ListItemText
              inset
              primary={t('Language.zh_CN')}
            />
          </MenuItem>
        </Menu>
      </React.Fragment>
    )}
  </I18n>;
}

export default LanguageMenu;
