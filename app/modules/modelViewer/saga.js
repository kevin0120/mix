import { take } from 'redux-saga/effects';
import { MODEL_VIEWER } from './action';
import { CommonLog } from '../../common/utils';

const { exec } = require('child_process');

export default function* root() {
  try {
    while (true) {
      yield take(MODEL_VIEWER.OPEN);
      // yield exec('');
    }
  } catch (e) {
    CommonLog.lError(e, {
      at: 'model viewer root'
    });
  }

}
