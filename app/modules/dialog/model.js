// @flow

import React from 'react';

export type dialogConfig = {
  hasOk: false,
  hasCancel: true,
  cancelAction: () => {},
  okAction: () => {},
  content: React.Component,
  title: React.Component
};

export type dialogState = {
  open: false,
  config: dialogConfig
};
