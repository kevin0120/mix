import React from 'react';
import { connect } from 'react-redux';
import Dialog from '@material-ui/core/Dialog/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import Slide from '@material-ui/core/Slide';
import { makeStyles } from '@material-ui/core/styles';
import { I18n } from 'react-i18next';
import { DialogActions } from '@material-ui/core';
import styles from './style';
import dialogActions from '../../modules/dialog/action';
import Button from '../CustomButtons/Button';
import { Dispatch } from '../../modules/indexReducer';

const mapState = (state, props) => ({
  ...props,
  dialog: state.dialog
});

const mapDispatch = {
  cancel: dialogActions.cancelDialog,
  ok: dialogActions.okDialog,
  close: dialogActions.closeDialog
};

const Transition = React.forwardRef((props, ref) => (
  <Slide direction="up" {...props} ref={ref} />
));

type Props = {
  dialog: {},
  cancel: Dispatch,
  ok: Dispatch,
  close: Dispatch
};

function customDialog(props: Props) {
  const { dialog, cancel, ok, close } = props;
  const classes = makeStyles(styles)();
  const { hasCancel, hasOk } = dialog?.config || {};

  const onCancel = () => {
    if (hasCancel && cancel) {
      cancel();
    }
    close();
  };
  const onOk = () => {
    if (hasOk && ok) {
      ok();
    }
    close();
  };

  return (
    <I18n>
      {t => (
        <Dialog
          classes={{
            root: classes.modalRoot,
            paper: `${classes.modal} ${classes.modalLarge}`
          }}
          TransitionComponent={Transition}
          keepMounted
          open={dialog.open}
          onClose={onCancel}
          aria-labelledby="form-dialog-title"
          scroll="paper"
        >
          <DialogTitle id="form-dialog-title" className={classes.modalHeader}>
            {t('Common.Result')}
          </DialogTitle>
          <DialogContent className={classes.modalBody}>{}</DialogContent>
          <DialogActions>
            {hasCancel ? (
              <Button onClick={onCancel} color="info" autoFocus round>
                {t('Common.Close')}
              </Button>
            ) : null}
            {hasOk ? (
              <Button onClick={onOk} color="info" autoFocus round>
                {t('Common.Yes')}
              </Button>
            ) : null}
          </DialogActions>
        </Dialog>
      )}
    </I18n>
  );
}

export default connect(
  mapState,
  mapDispatch
)(customDialog);
