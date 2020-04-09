import React, { useState } from 'react';
import DialogContent from '@material-ui/core/DialogContent';
import Dialog from '@material-ui/core/Dialog';
import { makeStyles } from '@material-ui/core/styles';
import MenuItem from '@material-ui/core/MenuItem';
import { Typography } from '@material-ui/core';
import Select from '@material-ui/core/Select';
import Button from '../CustomButtons/Button';
import { withI18n } from '../../i18n';
import styles from './styles';
import { blockReasons } from '../../modules/order/constants';

export function BlockReasonDialog({ AnchorButton, onConfirm }) {
  const classes = makeStyles(styles)();
  const [dialogOpen, setDialogOpen] = useState(false);
  const [selectValue, setSelectValue] = useState(null);
  return withI18n(t => <React.Fragment>
    <AnchorButton onClick={() => setDialogOpen(true)}/>
    <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
      <div style={{ backgroundColor: 'white' }}>
        <DialogContent className={classes.dialogContainer}>
          <div>
            <Typography
              variant="h4"
              color="textSecondary"
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
        </DialogContent>
      </div>
    </Dialog>
  </React.Fragment>);
}
