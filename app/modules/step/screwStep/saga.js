// @flow
import { put, select, take, call } from 'redux-saga/effects';
import { cloneDeep } from 'lodash-es';
import STEP_STATUS from '../model';
import { stepPayload, workingStep, stepData, workingOrder } from '../../order/selector';
import { SCREW_STEP } from './action';
import type { tResultAction } from './action';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './model';
import controllerModeTasks from './controllerModeTasks';
import handleResult from './handleResult';
import { CommonLog } from '../../../common/utils';


export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      // init data
      const payload: tScrewStepPayload = yield select(s => stepPayload(workingStep(workingOrder(s.order))));
      yield put(orderActions.stepData((data: tScrewStepData): tScrewStepData => {
        const points: Array<tPoint> = cloneDeep(payload?.points || []).sort((a, b) => a.group_sequence - b.group_sequence);
        return {
          points, // results data.results
          activeIndex: 0, // <-activeResultIndex
          ...data,
          jobID: payload.job_id,
          carType: payload.model,
          workSheet: payload.work_sheet,
          lnr: payload.lnr,
          workorderID: payload.workorder_id,
          controllerMode: payload.controllerMode
        };
      }));

      // enable tools

      //
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      while (true) {
        const data = yield select(s => stepData(workingStep(workingOrder(s.order))));

        console.log(data.controllerMode);
        // call controllerModeTasks(pset/job)
        yield call(controllerModeTasks[data.controllerMode], orderActions);

        // take result
        const { results }: tResultAction = yield take(SCREW_STEP.RESULT);
        console.log('result taken');

        // handle result
        yield call(handleResult, ORDER, orderActions, results, data);
      }
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * [STEP_STATUS.FAIL](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      CommonLog.lError(e);
    }
  }

};
