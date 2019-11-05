// @flow
import { put, select, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { controllerModes } from './constants';
import notifierActions from '../../Notifier/action';
import type { tPoint, tScrewStepData } from './interface/typeDef';
import { CommonLog } from '../../../common/utils';
import { jobApi, psetApi } from '../../../api/tools';
import { workingOrder } from '../../order/selector';

// pset/job模式
export default {
  *[controllerModes.pset](point: tPoint): Saga<boolean> {
    try {
      const sData: tScrewStepData = this._data;
      const stepId = this._id;
      const { points, retryTimes } = sData;
      const userIDs = yield select(s => s.users.map(u => u.uid));
      // const userID = 1;

      // const { toolSN, pset, sequence } = points[activeIndex];
      // eslint-disable-next-line camelcase
      const { tightening_tool, pset, sequence } = point;
      const total = points.length || 0;
      const workorderID = yield select(s => workingOrder(s.order)?.code);
      yield call(
        psetApi,
        // eslint-disable-next-line camelcase
        tightening_tool || '',
        stepId,
        // userID,
        userIDs,
        pset,
        sequence,
        retryTimes,
        total,
        workorderID
      );
    } catch (e) {
      // 程序号设置失败
      yield put(
        notifierActions.enqueueSnackbar('Error', `pset失败，${e.message}`, {
          // meta message,
          at: 'controllerModes.pset'
        })
      );
      return false;
    }
    return true;
  },

  *[controllerModes.job](): Saga<boolean> {
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
      yield put(
        notifierActions.enqueueSnackbar('Error', `程序号设置失败，${e.message}`, {
          // meta message,
          at: 'controllerModes.job'
        })
      );
      return false;
    }
    return true;
  }
};
