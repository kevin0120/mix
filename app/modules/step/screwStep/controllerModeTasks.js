// @flow
import { put, select, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { controllerModes } from './model';
import { jobManual, pset } from '../../../api/operation';
import { setNewNotification } from '../../notification/action';
import { OPERATION } from '../../operation/action';
import { toolEnable } from '../../tools/action';
import { stepData, workingStep, workingOrder } from '../../order/selector';
import type { tPoint, tResult, tScrewStepData } from './model';
import { CommonLog } from '../../../common/utils';

export default {
  * [controllerModes.pset](): Saga<void> {
    try {
      console.log('in pset progress');
      const rush = yield select(s => s.connections);
      const sData = yield select(stepData(workingStep(workingOrder(s.order))));
      const {
        activeResultIndex,
        failCount,
        results,
        workorderID
      } = sData;
      const userID = 1;
      yield call(
        pset,
        rush,
        results[activeResultIndex].controller_sn,
        results[activeResultIndex].gun_sn,
        0,
        failCount + 1,
        userID,
        results[activeResultIndex].pset,
        workorderID,
        results[activeResultIndex].group_sequence
      );
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
      const { operationID, carType, carID, jobID, source, points } = sData;

      const toolSN = points.reduce((tSN: string, p: tPoint): string => {
        if (tSN && p.toolSN !== tSN) {
          CommonLog.lError('结果中的toolSN不匹配');
        }
        return p.toolSN || tSN || '';
      }, '');

      const userID = 1;
      const skip = false;
      const hasSet = false;
      try {
        const resp = yield call(
          jobManual,
          rushUrl,
          controller_sn,
          toolSN,
          carType,
          carID,
          userID,
          jobID,
          hmiSn.value,
          operationID,
          skip,
          hasSet,
          source
        );

        if (resp.status === 200) {
          // 程序号设置成功

          // 设置workorder_id
          yield put({
            type: OPERATION.JOB_MANUAL.OK,
            workorderID: resp.data.workorder_id
          });
          // 启动用具
          yield put(toolEnable('开始作业'));
          // yield put({ type: OPERATION.STARTED });
        } else {
          CommonLog.lError('程序号设置失败');
        }
      } catch (e) {
        // 程序号设置失败
        CommonLog.lError(e);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  }
};

