import React, { useState } from 'react';
import DialogContent from '@material-ui/core/DialogContent';
import Dialog from '@material-ui/core/Dialog';
import { makeStyles } from '@material-ui/core/styles';
import MenuItem from '@material-ui/core/MenuItem';
import { DialogActions, Typography } from '@material-ui/core';
import Select from '@material-ui/core/Select';
import Button from '../CustomButtons/Button';
import { withI18n } from '../../i18n';
import styles from './styles';

export function BlockReasonDialog({ AnchorButton, onConfirm, blockReasons }) {
  const classes = makeStyles(styles)();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectValue, setSelectValue] = useState(null);
  return withI18n(t => <React.Fragment>
    <AnchorButton onClick={() => setDialogOpen(true)}/>
    <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
      <DialogContent className={classes.dialogContainer}>
        <div>
          <Typography
            variant="h5"
            color="textPrimary"
            align="left"
            className={classes.orderInfoText}
          >
            {t('OEE.SelectBlockReason')}
          </Typography>
          <Select
            labelId="block-reason-select-label"
            id="block-reason-select"
            value={selectValue}
            onChange={(event) => {
              setSelectValue(event.target.value);
            }}
          >
            {
              Object.values(blockReasons).map(r => <MenuItem value={r}>{t(`blockReasons.${r.name}`)}</MenuItem>)
            }
          </Select>
        </div>

      </DialogContent>
      <DialogActions>
        <Button
          type="button"
          disabled={!selectValue}
          onClick={() => {
            onConfirm(selectValue);
            setDialogOpen(false);
          }}
          variant="contained"
          color="warning"
        >
          {t('Common.Confirm')}
        </Button>
      </DialogActions>
    </Dialog>
  </React.Fragment>);
}
