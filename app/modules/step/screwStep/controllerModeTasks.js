// @flow
import { put, select, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { controllerModes } from './constants';
import notifierActions from '../../Notifier/action';
import type { tPoint, tScrewStepData } from './interface/typeDef';
import { CommonLog } from '../../../common/utils';
import { jobApi, psetApi } from '../../../api/tools';
import { workingOrder } from '../../order/selector';

// pset/job模式
export default {
  * [controllerModes.pset](point: tPoint): Saga<void> {
    try {
      console.log(point);
      const sData: tScrewStepData = this.data;
      const stepId = this.id;
      const { retryTimes } = sData;
      const { points } = this._pointsManager;
      const userIDs: Array<number> = yield select(s => s.users.map(u => u.uid));

      // eslint-disable-next-line camelcase
      const { tightening_tool: toolSN, pset, sequence } = point;
      const total = points.length || 0;
      const workorderID = yield select(s => workingOrder(s.order)?.id);
      if (isNil(pset)) {
        throw new Error('pset号为空');
      }
      if (isNil(workorderID)) {
        throw new Error('工单ID为空');
      }
      yield call(
        psetApi,
        toolSN || '',
        stepId,
        userIDs,
        pset,
        sequence,
        retryTimes,
        total,
        workorderID
      );
    } catch (e) {
      // 程序号设置失败
      const msg = `pset失败，${e.message}`;
      yield put(
        notifierActions.enqueueSnackbar('Error', msg, {
          at: 'controllerModes.pset'
        })
      );
      throw new Error(msg);
    }
  },

  * [controllerModes.job](): Saga<void> {
    try {
      const { jobID, points }: tScrewStepData = this._data;
      const stepId = this._id;
      const toolSN = points.reduce((tSN: string, p: tPoint): string => {
        if (tSN && p.tightening_tool !== tSN) {
          CommonLog.lError('结果中的toolSN不匹配');
        }
        return p.tightening_tool || tSN || '';
      }, '');

      const userID = 1;

      yield call(jobApi, toolSN, stepId, userID, jobID);
    } catch (e) {
      const msg = `程序号设置失败，${e.message}`;

      yield put(
        notifierActions.enqueueSnackbar('Error', msg, {
          at: 'controllerModes.job'
        })
      );
      throw new Error(msg);
    }
  }
};
