// @flow
import { put, select, take, call } from 'redux-saga/effects';
import { cloneDeep } from 'lodash-es';
import STEP_STATUS from '../model';
import { stepPayload, workingStep, stepData, workingOrder } from '../../order/selector';
import { SCREW_STEP } from './action';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './model';
import controllerModeTasks from './controllerModeTasks';
import handleResult from './handleResult';
import { CommonLog } from '../../../common/utils';
import { staticScrewTool } from '../../tools/saga';
// import { toolEnableApi } from '../../../api/order';


export default {


  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      // init data
      const payload: tScrewStepPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      // const { result, msg } = yield call(toolEnableApi, staticScrewTool.SerialNumber, false);
      // console.log(result, msg);
      yield put(orderActions.stepData((data: tScrewStepData): tScrewStepData => {
        const points: Array<tPoint> = cloneDeep(payload?.points || []).sort((a, b) => a.group_sequence - b.group_sequence);
        return {
          points, // results data.results
          activeIndex: -1, // <-activeResultIndex
          ...data,
          jobID: payload.jobID,
          controllerMode: payload.controllerMode,
          retryTimes: 0
        };
      }));
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      CommonLog.lError(e, { at: 'screwStep ENTERING' });
    }
  },


  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      const initData = yield select(s => stepData(workingStep(workingOrder(s.order))));
      yield call(handleResult, ORDER, orderActions, [], initData);
      let isFirst = true;
      while (true) {
        const data = yield select(s => stepData(workingStep(workingOrder(s.order))));
        if (data.controllerMode === 'pset' || (data.controllerMode === 'job' && isFirst)) {
          const succ = yield call(controllerModeTasks[data.controllerMode], orderActions);
          if (!succ) {
            // TODO: on set pset fail
            yield put(orderActions.stepStatus(STEP_STATUS.FAIL));
          }
          if (isFirst) {
            yield call([staticScrewTool, staticScrewTool.Enable]);
          }
          isFirst = false;
        }
        const { results: { data: results } } = yield take(SCREW_STEP.RESULT);
        yield call(handleResult, ORDER, orderActions, results, data);
      }
    } catch (e) {
      CommonLog.lError(e, { at: 'screwStep DOING' });
    } finally {
      yield call([staticScrewTool, staticScrewTool.Disable]);
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
