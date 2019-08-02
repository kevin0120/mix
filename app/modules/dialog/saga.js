import { take, call, put } from 'redux-saga/effects';
import { DIALOG } from './action';
import type { tDialogConfig } from './model';

export default function* root() {
  try {
    while (true) {
      const action = yield take(DIALOG.SHOW);
      yield call(showDialog, action);
    }
  } catch (e) {
    console.error(e);
  }
}

const dialogActions = {
  * [DIALOG.BUTTON](config: tDialogConfig,{idx}) {
    try {
      if (config?.buttons?.[idx]?.action) {
        yield put(config.buttons[idx].action);
      }
    } catch (e) {
      console.error(e);
    }
  },
  * [DIALOG.CLOSE](config) {
    try {
      if (config?.closeAction) {
        yield put(config.closeAction);
      }
    } catch (e) {
      console.error(e);
    }
  }
};

function* showDialog(action) {
  try {
    const { config } = action;
      const dialogAction = yield take(Object.keys(dialogActions));
      if (dialogActions[dialogAction.type]) {
        yield call(dialogActions[dialogAction.type], config,dialogAction);
      }
  } catch (e) {
    console.error(e);
  }
}
