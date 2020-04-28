import { model, watch } from 'saga-modux';
import { call, delay, put, select } from 'redux-saga/effects';
import React from 'react';
import { CommonLog } from '../../common/utils';
import dialogActions from '../dialog/action';

const viewingDialog = (endpoint) => ({
  maxWidth: 'xl',
  buttons: [{
    label: 'Common.Close',
    color: 'warning'
  }],
  content: <iframe src={endpoint} style={{ width: '80vw', height: '75vh' }}/>
});

export default model('viewOperationInfo', {
  * root() {
    try {
      yield call(watch([
        this.viewOperationResources,
        this.viewOperationTracing
      ]));
    } catch (e) {
      CommonLog.lError(e, { at: 'viewOperationInfo model root' });
    }
  },

  * viewOperationResources() {
    try {
      const { operationResources } = yield select(s => s.setting.system.endpoints);
      const envDialog = viewingDialog(operationResources);
      yield put(dialogActions.dialogShow(envDialog));
    } catch (e) {
      CommonLog.lError(e, { at: 'viewOperationInfo.viewOperationResources' });
    }
  },

  * viewOperationTracing() {
    try {
      yield delay(300);
      const { operationTracing } = yield select(s => s.setting.system.endpoints);
      const operationTracingDialog = viewingDialog(operationTracing);
      yield put(dialogActions.dialogShow(operationTracingDialog));
    } catch (e) {
      CommonLog.lError(e, { at: 'viewOperationInfo.viewOperationResources' });
    }
  }

}, {});
