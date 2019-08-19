// @flow
import { put, select, call } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { controllerModes } from './model';
import { setNewNotification } from '../../notification/action';
import { stepData, workingStep, workingOrder, workingIndex } from '../../order/selector';
import type { tPoint, tScrewStepData } from './model';
import { CommonLog } from '../../../common/utils';
import { jobApi, psetApi } from '../../../api/tools';

export default {


  * [controllerModes.pset](point): Saga<void> {
    try {
      const sData: tScrewStepData = this._data;
      const stepId = this._id;
      const {
        activeIndex,
        points,
        retryTimes
      } = sData;
      // TODO: pass correct userID
      const userID = 1;

      // const { toolSN, pset, sequence } = points[activeIndex];
      const { toolSN, pset, sequence } = point;

      const data = yield call(psetApi, toolSN || '', stepId, userID, pset, sequence, retryTimes);
      if (data && data.result !== 0) {
        CommonLog.lError(`pset失败${data.msg}`, {
          at: 'pset',
          toolSN, stepId, userID, pset, sequence, retryTimes
        });
        return false;
      }
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
      const { jobID, points }: tScrewStepData = this._data;
      const stepId = this._id;
      const toolSN = points.reduce((tSN: string, p: tPoint): string => {
        if (tSN && p.toolSN !== tSN) {
          CommonLog.lError('结果中的toolSN不匹配');
        }
        return p.toolSN || tSN || '';
      }, '');

      const userID = 1;

      const data = yield call(jobApi, toolSN, stepId, userID, jobID);
      if (data && data.result !== 0) {
        CommonLog.lError(`程序号设置失败:${data.msg}`, {
          at: 'job',
          toolSN, stepId, userID, jobID
        });
        return false;
      }
    } catch (e) {
      CommonLog.lError(e);
      return false;
    }
    return true;
  }
};

