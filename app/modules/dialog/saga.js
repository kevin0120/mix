import { take, call, put } from 'redux-saga/effects';
import { DIALOG } from './action';

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
  *[DIALOG.OK](config) {
    try {
      if (config?.okAction) {
        yield put(config.okAction);
      }
    } catch (e) {
      console.error(e);
    }
  },
  *[DIALOG.CANCEL](config) {
    try {
      if (config?.cancelAction) {
        yield put(config.cancelAction);
      }
    } catch (e) {
      console.error(e);
    }
  },
  *[DIALOG.CLOSE](config) {
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
    while (true) {
      const { config } = action;
      const dialogAction = yield take(Object.keys(dialogActions));
      if (dialogActions[dialogAction.type]) {
        yield call(dialogActions[dialogAction.type], config);
      }
    }
  } catch (e) {
    console.error(e);
  }
}
