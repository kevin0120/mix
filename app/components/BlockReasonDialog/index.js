import React from 'react';
import i18n, { withI18n } from '../../i18n';
import { orderActions } from '../../modules/order/action';

function BlockReasonDialogContent() {

  return withI18n(t =>
    <div />
  );
}

export function BlockReasonDialog() {

  return {
    maxWidth: 'md',
    buttons: [
      {
        label: i18n.t('Common.Cancel'),
        color: 'danger'
      },
      {
        label: i18n.t('Common.Confirm'),
        color: 'info',
        action: orderActions.pendingOrder()
      }
    ],
    title: i18n.t('OEE.title'),
    content: BlockReasonDialogContent()
  };
}
