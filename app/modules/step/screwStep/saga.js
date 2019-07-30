import { put, select, take, call } from 'redux-saga/effects';
import { STEP_STATUS } from '../model';
import { stepPayload, workingStep, stepData } from '../../order/selector';
import { SCREW_STEP } from './action';


const mergePointsAndResults = (points, result, activeIndex) => {
  const newPoints = [...points];
  newPoints.splice(activeIndex, result.length,
    ...newPoints.slice(activeIndex, activeIndex + result.length).map((p, idx) => ({
      ...p,
      ...result[idx]
    })));
  return newPoints;
};

const resultStatus = (result, data) => {
  const retry = result.some(r => r.status === 'error') && 'retry';
  const fail = retry && data.retryTimes >= data.maxRetryTimes && 'fail';
  const finish = !retry && (data.activeIndex + result.length >= data.points.length) && 'finish';
  const next = !retry && !finish && 'next';
  return [fail, retry, finish, next];
};

const resultStatusTasks = (ORDER, orderActions, result) => ({
  * retry() {
    try {
      yield put(orderActions.stepData((d) => ({
        ...d,
        points: mergePointsAndResults(d.points, result, d.activeIndex),
        retryTimes: (d.retryTimes || 0) + 1
      })));
    } catch (e) {
      console.error(e);
    }
  },
  * fail() {
    try {
      yield put(orderActions.stepData((d) => ({
        ...d,
        activeIndex: -1,
        points: mergePointsAndResults(d.points, result, d.activeIndex)
      })));
      yield put(orderActions.stepStatus(STEP_STATUS.FAIL));
    } catch (e) {
      console.error(e);
    }
  },
  * finish() {
    try {
      yield put(orderActions.stepData((d) => ({
        ...d,
        activeIndex: -1,
        points: mergePointsAndResults(d.points, result, d.activeIndex)
      })));
      yield put(orderActions.stepStatus(STEP_STATUS.FINISHED));
    } catch (e) {
      console.error(e);
    }
  },
  * next() {
    try {
      yield put(orderActions.stepData((d) => ({
        ...d,
        activeIndex: d.activeIndex + result.length,
        points: mergePointsAndResults(d.points, result, d.activeIndex)
      })));
    } catch (e) {
      console.error(e);
    }
  }
});

export default {
  * [STEP_STATUS.ENTERING](ORDER, orderActions) {
    try {
      const payload = yield select(s => stepPayload(workingStep(s.order)));
      yield put(orderActions.stepData((data) => {
        const points = JSON.parse(JSON.stringify(payload?.points || []));
        return {
          ...data,
          points,
          maxRetryTimes: payload.maxRetryTimes,
          activeIndex: -1
        };
      }));
      yield put(orderActions.stepStatus(STEP_STATUS.DOING));
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.DOING](ORDER, orderActions) {
    try {
      yield take(SCREW_STEP.IMAGE_READY);
      yield put(orderActions.stepData((data) => ({
        ...data,
        activeIndex: 0
      })));
      while (true) {
        const { result } = yield take(SCREW_STEP.RESULT);
        console.log('result taken');
        const data = yield select(s => stepData(workingStep(s.order)));

        const firstMatchResultStatus =
          resultStatusTasks(
            ORDER,
            orderActions,
            result
          )[resultStatus(result, data).find(v => !!v)];

        if (firstMatchResultStatus) {
          yield call(firstMatchResultStatus);
        }
      }
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.FINISHED](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      console.error(e);
    }
  },
  * [STEP_STATUS.FAIL](ORDER, orderActions) {
    try {
      yield put(orderActions.doNextStep());
    } catch (e) {
      console.error(e);
    }
  }

};
