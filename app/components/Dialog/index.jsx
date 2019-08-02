// @flow

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
import type { tDialogConfig } from '../../modules/dialog/model';

const mapState = (state, props) => ({
  ...props,
  config: state?.dialog?.config || {},
  open: state?.dialog?.open || false
});

const mapDispatch = {
  closeAction: dialogActions.dialogClose,
  buttonAction: dialogActions.dialogButton
};

const Transition = React.forwardRef((props, ref) => (
  <Slide direction="up" {...props} ref={ref}/>
));

type Props = {
  config: tDialogConfig,
  open: boolean,
  buttonAction: Dispatch,
  closeAction: Dispatch
};

function customDialog(props: Props) {
  const { config, open, buttonAction, closeAction } = props;
  const classes = makeStyles(styles)();
  const { buttons, title, content } = config;

  const onButton = (idx) => {
    buttonAction(idx);
    closeAction();
  };
  const onClose = () => {
    closeAction();
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
          open={open}
          onClose={onClose}
          aria-labelledby="form-dialog-title"
          scroll="paper"
        >
          <DialogTitle id="form-dialog-title" className={classes.modalHeader}>
            {(typeof title === 'string' ? t(title) : title) || ''}
          </DialogTitle>
          <DialogContent className={classes.modalBody}>
            {content || ''}
          </DialogContent>
          <DialogActions>
            {
              buttons?.map((b, idx) => b ?
                <Button key={b.label} onClick={() => onButton(idx)} color={b.color || 'info'} regular>
                  {t(b.label || '')}
                </Button> : null) || null
            }
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
