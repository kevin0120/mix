import React from 'react';
import MenuList from '@material-ui/core/MenuList';
import MenuItem from '@material-ui/core/MenuItem';
import Typography from '@material-ui/core/Typography';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import makeStyles from '@material-ui/core/styles/makeStyles';
import { push } from 'connected-react-router';
import { connect } from 'react-redux';
import { I18n } from 'react-i18next';
import LeftMenuWithAvatar from '../../components/LeftMenuWithAvatar';

import styles from './styles';

function Preference({ children, childRoutes, doPush, path }) {
  const classes = makeStyles(styles.config)();

  const menuItems = () => childRoutes.map(r => r && (
    <I18n ns="translations" key={r.title}>
      {t => (
        <MenuItem
          selected={path === r.url}

          className={classes.menuItem}
          component={React.forwardRef((p, ref) => (
            <Card
              onClick={() => doPush(r.url)}
              className={
                path === r.url
                  ? classes.menuItemSelected
                  : classes.menuItem
              }
              ref={ref}
            >
              <CardActionArea
                classes={{
                  root:
                    path === r.url
                      ? classes.cardActionAreaSelected
                      : classes.cardActionArea
                }}
              >
                {<r.icon style={{ fontSize: 45, margin: 5, fill: path === r.url ? '#FAFAFA' : r.color }}/>}
                <span className={classes.itemText}>
                  <Typography variant="h6">
                  {t(r.title)}
                  </Typography>
                </span>
              </CardActionArea>
            </Card>
          ))}
        />)}
    </I18n>
  ));

  return (
    <div className={classes.root}>
      {/* <AppBarBack /> */}
      <LeftMenuWithAvatar>
        <MenuList>{menuItems()}</MenuList>
      </LeftMenuWithAvatar>
      <div className={classes.content}>
        {children}
      </div>
    </div>
  );
}


const mapState = (state, ownProps) => ({
  path: state.router.location.pathname,
  ...ownProps
});

const mapDispatch = {
  doPush: push
};

export default connect(
  mapState,
  mapDispatch
)(Preference);
