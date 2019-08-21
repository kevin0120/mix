import React from 'react';
import { GridLoader } from 'react-spinners';
import { connect } from 'react-redux';
import Dialog from '@material-ui/core/Dialog';
import { makeStyles } from '@material-ui/styles';
import styles from './style';
import Fade from '@material-ui/core/Fade';

function Loading({ loading }) {
  const classes=makeStyles(styles)();
  return <Dialog
    fullScreen
    classes={{paper:classes.root}}
    open={loading}
    style={{opacity: 0.7}}
    TransitionComponent={Fade}
  >
    <GridLoader
      sizeUnit="px"
      size={50}
      color="#36D7B7"
      loading={loading}
      className={classes.spinner}
    />
  </Dialog>;
}

const mapDispatch = {};

const mapState = (state, props) => ({
  ...props,
  loading: state.loading
});

export default connect(mapState, mapDispatch)(Loading);
