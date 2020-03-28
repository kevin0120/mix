// @flow
import { call, put, select } from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
import { isNil } from 'lodash-es';
import { controllerModes } from './constants';
import notifierActions from '../../Notifier/action';
import type { tPoint, tScrewStepData } from './interface/typeDef';
import { jobApi, psetApi } from '../../../api/tools';
import { workingOrder } from '../../order/selector';
import type { IDevice } from '../../device/IDevice';

// pset/job模式
export default {
  * [controllerModes.pset](sequence, tool, pset): Saga<void> {
    try {
      const sData: tScrewStepData = this.data;
      const stepCode = this.code;
      // const { retryTimes } = sData;
      const { points } = this._pointsManager;
      const userIDs: Array<number> = yield select(s => s.users.map(u => u.uid));
      if (!tool) {
        throw new Error(`未指定工具`);
      }
      const ControllerSN = ((tool.parent: any): IDevice)?.serialNumber;
      if (!ControllerSN) {
        throw new Error(`工具(${tool?.serialNumber})缺少控制器`);
      }
      const total = points.length;
      const workorderCode = yield select(s => workingOrder(s.order)?.code);
      if (isNil(pset)) {
        throw new Error('pset号为空');
      }
      if (isNil(workorderCode)) {
        throw new Error('工单ID为空');
      }
      const retries = 1;
      for (let retry = 0; retry <= retries; retry += 1) {
        try {
          yield call(
            psetApi,
            tool?.serialNumber || '',
            ControllerSN || '',
            stepCode,
            userIDs,
            pset,
            sequence,
            total,
            workorderCode,
            ''
          );
          break;
        } catch (e) {
          if (retry === retries) {
            throw e;
          }
        }
      }

    } catch (e) {
      // 程序号设置失败
      const msg = `pset失败，${e.message}, 工具：${tool?.serialNumber}`;
      yield put(
        notifierActions.enqueueSnackbar('Error', msg, {
          at: 'controllerModes.pset'
        })
      );
      throw new Error(msg);
    }
  },

  * [controllerModes.job](sequence, tool, jobID): Saga<void> {
    try {
      const stepId = this._id;
      const userIDs: Array<number> = yield select(s => s.users.map(u => u.uid));
      if (!tool) {
        throw new Error(`未指定工具`);
      }
      const ControllerSN = ((tool.parent: any): IDevice)?.serialNumber;
      if (!ControllerSN) {
        throw new Error(`工具(${tool.Name})缺少控制器`);
      }
      yield call(jobApi, tool.serialNumber, ControllerSN, stepId, userIDs, jobID);
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
