import { take } from 'redux-saga/effects';
import { MODEL_VIEWER } from './action';
import { CommonLog } from '../../common/utils';

const { exec } = require('child_process');

export default function* root() {
  try {
    while (true) {
      const {url}=yield take(MODEL_VIEWER.OPEN);
      if(url){
        yield exec(`firefox ${url}`);
      }else {
        CommonLog.lError('model viewing url not valid', {
          at: 'model viewer root',
          url
        });
      }
    }
  } catch (e) {
    CommonLog.lError(e, {
      at: 'model viewer root'
    });
  }

}
