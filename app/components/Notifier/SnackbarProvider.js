import React from 'react';
import { makeStyles } from '@material-ui/core';
import { SnackbarProvider } from 'notistack';

const styles = (t) => ({
  success: {
    backgroundColor: t.palette.success.main,
    color: t.palette.common.white,
    ...t.typography.body1
  },
  error: {
    backgroundColor: t.palette.danger.main,
    color: t.palette.common.white,
    ...t.typography.body1

  },
  warning: {
    backgroundColor: t.palette.warning.main,
    color: t.palette.common.white,
    ...t.typography.body1

  },
  info: {
    backgroundColor: t.palette.info.main,
    color: t.palette.common.white,
    ...t.typography.body1

  }
});

const StyledSnackbarProvider = (props) => {
  const { children, ...restProps } = props;
  const classes = makeStyles(styles)();

  return (
    <SnackbarProvider
      classes={{
        variantSuccess: classes.success,
        variantError: classes.error,
        variantWarning: classes.warning,
        variantInfo: classes.info
      }}
      {...restProps}
    >
      {children}
    </SnackbarProvider>
  );
};

export default StyledSnackbarProvider;
