import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { I18n } from 'react-i18next';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import Divider from '@material-ui/core/Divider';

import { routeConfigs } from '../../routes';

import styles from './styles';

/* eslint-disable react/prefer-stateless-function */
class ConnectedNavBar extends React.Component {

  render() {
    const { classes, styleOptions } = this.props;

    return (
      <I18n ns="translations">
        {
          t => (
            <div>
              <List className={classes.sideNav}>
                {routeConfigs.slice(0,-1).map(route => (
                  <li className={classes.itemWrap} key={route.name}>
                    <ListItem button component="a" href={`#${route.url}`} className={classes.menuItem}>
                      <ListItemIcon>
                        <route.icon style={{ color: route.color }} />
                      </ListItemIcon>
                      <ListItemText
                        primary={t(route.title)}
                        disableTypography={styleOptions.disableTypography}
                      />
                    </ListItem>
                    <Divider />
                  </li>
                ))}
                </List>
              </div>
            )
          }
      </I18n>
    );
  }
}

ConnectedNavBar.propTypes = {
  classes: PropTypes.shape({
  }).isRequired,
  styleOptions: PropTypes.shape({
    disableTypography: PropTypes.bool,
  }),
};

ConnectedNavBar.defaultProps = {
  styleOptions: {
    disableTypography: true,
  },
};

const NavBar = withStyles(styles)(ConnectedNavBar);

export default NavBar;
