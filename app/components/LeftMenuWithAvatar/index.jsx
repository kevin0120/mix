import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Avatar from '@material-ui/core/Avatar';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import Clock from 'react-live-clock';

import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  name: state.userInfo.name,
  secondaryInfo: state.userInfo.login, // email or role etc
  avatarImg: state.userInfo.image_small,
  ...ownProps,
});

const mapDispatchToProps = {
};

function ConnectedLeftMenuWithAvatar(props) {
  const {
    classes,
    name,
    secondaryInfo,
    avatarImg,
  } = props;

  return (
    <Drawer
      variant="permanent"
      classes={{
        paper: classes.drawerPaper,
      }}
    >
      <div className={classes.toolbar} />

      {props.children}

      <div className={classes.baseInfo}>
        <List>
          <ListItem className={classes.userInfo}>
            <Avatar
              alt={name}
              src={avatarImg}
              className={classes.avatar}
            />
            <ListItemText className={classes.userText} primary={name} secondary={secondaryInfo} />
          </ListItem>
          <ListItem className={classes.timeWrap}>
            <Clock
              className={classes.timeContent}
              format="HH:mm:ss"
              ticking
              timezone="Asia/Shanghai"
            />
          </ListItem>
        </List>
      </div>
    </Drawer>
  );
}

ConnectedLeftMenuWithAvatar.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  children: PropTypes.shape({}).isRequired,
  name: PropTypes.string.isRequired,
  secondaryInfo: PropTypes.string.isRequired,
  avatarImg: PropTypes.string.isRequired,
};

const LeftMenuWithAvatar = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ConnectedLeftMenuWithAvatar);

export default withStyles(styles)(LeftMenuWithAvatar);
