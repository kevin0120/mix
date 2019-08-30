import React from 'react';
import SweetAlert from 'react-bootstrap-sweetalert';
import { makeStyles } from '@material-ui/core';

const styles = {
  root:{
    '& .lead':{
      maxHeight:'40vh',
      overflow:'auto'
    }
  }
};

export default function (props) {
  const { show, style, children, ...restProps } = props;
  const classes=makeStyles(styles)();
  return (
    show ? (<div
        style={{
          position: 'fixed',
          height: '100vh',
          width: '100vw',
          top: 0,
          left: 0,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          alignContents: 'center',
          zIndex: 99999
        }}
      >
        <SweetAlert
          show={show}
          customClass={classes.root}
          style={{
            display: 'block',
            position: 'relative',
            margin: 'auto',
            top: 'auto',
            left: 'auto',
            maxHeight: '100%',
            overflow: 'auto',
            ...style
          }}
          {...restProps}
        >
          {children}
        </SweetAlert>
      </div>)
      : null
  );
}
