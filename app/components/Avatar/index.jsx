import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import { Grid } from '@material-ui/core';

const styles = {
  menuBtnWrapAvatar: {
    display: 'flex',
    borderStyle: 'solid',
    borderRadius: '50%',
    borderWidth: '1px',
    overflow: 'hidden',
    flex: '1',
    maxHeight: '100%',
    maxWidth: '100%',
    margin:'2px -10px 2px 0',
  },
  container: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    marginRight:'10px'
  }
};

export default withStyles(styles)((props) => {
  const { className, avatars, classes } = props;
  return <div className={`${classes.container} ${className || ''}`}>
    {avatars.map((a) => <div className={classes.menuBtnWrapAvatar}>
      <div style={{
        height: '0', minWidth: '100%', paddingBottom: '100%', overFlow: 'visible',
        backgroundImage: `url(${a})`,
        backgroundSize: 'cover',
        backgroundPosition: 'center'
      }}/>
    </div>)}
  </div>;
});
