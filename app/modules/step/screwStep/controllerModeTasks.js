// @flow
import { put, select, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { controllerModes } from './model';
import { pset } from '../../../api/operation';
import { setNewNotification } from '../../notification/action';
import { stepData, workingStep, workingOrder } from '../../order/selector';
import type { tPoint, tScrewStepData } from './model';
import { CommonLog } from '../../../common/utils';
import rushActions from '../../rush/action';


export default {
  * [controllerModes.pset](): Saga<void> {
    try {
      console.log('in pset progress');
      const rush = yield select(s => s.connections);
      const sData = yield select(s => stepData(workingStep(workingOrder(s.order))));
      const {
        activeResultIndex,
        failCount,
        results,
        workorderID
      } = sData;
      const userID = 1;
      yield put(rushActions.sendJson({

      }))
      // yield call(
      //   pset,
      //   rush,
      //   results[activeResultIndex].controller_sn,
      //   results[activeResultIndex].gun_sn,
      //   0,
      //   failCount + 1,
      //   userID,
      //   results[activeResultIndex].pset,
      //   workorderID,
      //   results[activeResultIndex].group_sequence
      // );
    } catch (e) {
      // 程序号设置失败
      yield put(setNewNotification('Error', 'pset failed', {
        // meta message
      }));
      return false;
    }
    return true;
  },
  * [controllerModes.job](): Saga<void> {
    try {
      const sData: tScrewStepData = yield select(s => stepData(workingStep(workingOrder(s.order))));
      const rushUrl = yield select(s => s.setting.system.connections.rush);

      const { hmiSn } = yield select(s => s.setting.page.odooConnection);
      const { jobID, points } = sData;

      const toolSN = points.reduce((tSN: string, p: tPoint): string => {
        if (tSN && p.toolSN !== tSN) {
          CommonLog.lError('结果中的toolSN不匹配');
        }
        return p.toolSN || tSN || '';
      }, '');

      const userID = 1;

      // TODO: set job here
      // yield put(rushActions.sendJson({
      //
      // }))
      // if (resp.status === 200) {
      //   // 程序号设置成功
      //
      //   // 设置workorder_id
      //   yield put({
      //     type: OPERATION.JOB_MANUAL.OK,
      //     workorderID: resp.data.workorder_id
      //   });
      //   // 启动用具
      //   yield put(toolEnable('开始作业'));
      //   return true;
      // }
      CommonLog.lError('程序号设置失败');

    } catch (e) {
      CommonLog.lError(e);
      return false;
    }
    return false;
  }
};

