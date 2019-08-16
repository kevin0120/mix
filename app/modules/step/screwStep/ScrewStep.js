import Step from '../Step';
import STEP_STATUS from '../model';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './model';
import { call, put, select, take } from 'redux-saga/effects';
import { cloneDeep } from 'lodash-es';
import { CommonLog } from '../../../common/utils';
import handleResult from './handleResult';
import controllerModeTasks from './controllerModeTasks';
import { staticScrewTool } from '../../external/device/tools/saga';
import { SCREW_STEP } from './action';


export default class ScrewStep extends Step {
  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        // init data
        const payload: tScrewStepPayload = this._payload;
        this.updateData((data: tScrewStepData): tScrewStepData => {
          const points: Array<tPoint> = cloneDeep(payload?.points || []).sort((a, b) => a.group_sequence - b.group_sequence);
          return {
            points, // results data.results
            activeIndex: -1, // <-activeResultIndex
            ...data,
            jobID: payload.jobID,
            controllerMode: payload.controllerMode,
            retryTimes: 0
          };
        });
        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep ENTERING' });
      }
    },


    * [STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        yield call([this,handleResult], ORDER, orderActions, [], this._data);
        let isFirst = true;
        while (true) {
          const data = this._data;
          if (data.controllerMode === 'pset' || (data.controllerMode === 'job' && isFirst)) {
            const succ = yield call([this,controllerModeTasks[data.controllerMode]], orderActions);
            if (!succ) {
              // TODO: on set pset fail
              yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
            }
            if (isFirst) {
              yield call(staticScrewTool.Enable);
            }
            isFirst = false;
          }
          const { results: { data: results } } = yield take(SCREW_STEP.RESULT);
          yield call([this,handleResult], ORDER, orderActions, results, data);
        }
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep DOING' });
      } finally {
        yield call(staticScrewTool.Disable);
      }
    },


    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.doNextStep());
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep FINISHED' });
      }
    },


    * [STEP_STATUS.FAIL](ORDER, orderActions) {
      try {
        yield put(orderActions.doNextStep());
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep FAIL' });
      }
    }
  };
}