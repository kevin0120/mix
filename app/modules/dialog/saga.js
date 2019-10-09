import { take, call, put, fork, takeEvery } from 'redux-saga/effects';
import actions, { DIALOG } from './action';
import type { tDialogConfig } from './model';
import { CommonLog } from '../../common/utils';

export default function* root() {
  try {
    while (true) {
      const action = yield take(DIALOG.SHOW);
      const { config } = action;
      yield fork(showDialog, action);
      yield take(DIALOG.CLOSE);
      yield fork(handleClose, config);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

const dialogActions = {
  *[DIALOG.BUTTON](config: tDialogConfig, { idx }) {
    try {
      yield put(actions.dialogClose());
      if (config?.buttons?.[idx]?.action) {
        yield put(config.buttons[idx].action);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};
function* handleClose(config) {
  try {
    if (config?.closeAction) {
      yield put(config.closeAction);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* showDialog(action) {
  try {
    const { config } = action;
    const dialogAction = yield take(Object.keys(dialogActions));
    if (dialogActions[dialogAction.type]) {
      yield fork(dialogActions[dialogAction.type], config, dialogAction);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}
