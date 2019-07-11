import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import { Grid } from '@material-ui/core';

const styles = {
  avatarBorder: {
    display: 'flex',
    borderStyle: 'solid',
    borderRadius: '50%',
    borderWidth: '1px',
    overflow: 'hidden',
    flex: '1',
    maxHeight: '100%',
    maxWidth: '100%',
    // margin: '2px -10px 2px 0'
  },
  container: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginRight: '10px'
  }
};

export default withStyles(styles)((props) => {
  const { className, users, classes, onClickAvatar } = props;
  return <div className={`${classes.container} ${className || ''}`}>
    {users.map(user => <div key={`${user.name} ${user.uuid}`} className={classes.avatarBorder}>
      <div
        style={{
          height: '0', minWidth: '100%', paddingBottom: '100%', overFlow: 'visible',
          backgroundImage: `url(${user.avatar})`,
          backgroundSize: 'cover',
          backgroundPosition: 'center'
        }}
        onClick={() => onClickAvatar(user.uuid)}
      />
    </div>)}
  </div>;
});
