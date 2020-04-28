// @flow
import { take, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { MODEL_VIEWER } from './action';
import { CommonLog } from '../../common/utils';

const { exec } = require('child_process');

export default function* root(): Saga<void> {
  try {
    while (true) {
      const { url } = yield take(MODEL_VIEWER.OPEN);
      console.log(url);
      try{
        if(!url){
          throw new Error('model viewing url not valid')
        }
        yield call(exec, `firefox ${url}`);
      }catch (e) {
        CommonLog.lError(`查看模型错误${e.message}`, {
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
