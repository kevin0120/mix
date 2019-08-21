import { put, call } from 'redux-saga/effects';
import type { tPoint, tPointStatus, tResult, tResultStatus, tScrewStepData } from './model';
import { POINT_STATUS, RESULT_STATUS } from './model';
import STEP_STATUS from '../model';
import { CommonLog } from '../../../common/utils';

function formPointStatusFromResultStatus(point: tPoint, rStatus: tResultStatus, activeGroupSequence: number): tPointStatus {
  let pStatus = POINT_STATUS.WAITING;

  if (rStatus === RESULT_STATUS.nok) {
    pStatus = POINT_STATUS.ERROR;
  } else if (rStatus === RESULT_STATUS.ok) {
    pStatus = POINT_STATUS.SUCCESS;
  }

  const isActive = activeGroupSequence && activeGroupSequence === point.group_sequence;
  if (isActive) {
    switch (pStatus) {
      case POINT_STATUS.ERROR: {
        pStatus = POINT_STATUS.ERROR_ACTIVE;
        break;
      }
      case POINT_STATUS.WAITING:
        pStatus = POINT_STATUS.WAITING_ACTIVE;
        break;
      default:
        break;
    }
  }

  return pStatus;
}

const mergePointsAndResults = (points: Array<tPoint>, results: Array<tResult>, activeIndex: number, activeGroupSequence): Array<tPoint> => {
  const newPoints = [...points];
  if (activeIndex === -1) {
    return [
      ...newPoints.map((p) => {
        return ({
          ...p,
          status: formPointStatusFromResultStatus(p, null, activeGroupSequence)
        });
      })];
  }
  newPoints.splice(activeIndex, newPoints.length - activeIndex,
    ...newPoints.slice(activeIndex).map((p, idx) => {
      const r: tResult = results[idx];
      if (r) {
        return ({
          ...p,
          ti: r.ti,
          mi: r.mi,
          wi: r.wi,
          batch: r.batch,
          status: formPointStatusFromResultStatus(p, r.result, activeGroupSequence)
        });
      }
      return ({
        ...p,
        status: formPointStatusFromResultStatus(p, null, activeGroupSequence)
      });
    }));
  return newPoints;
};

const resultStatus = (results: Array<tResult>, data: tScrewStepData) => {
  const LSN = results.some((r: tResult): boolean => r.result === RESULT_STATUS.lsn) && 'LSN';
  const retry = results.some((r: tResult): boolean => r.result === RESULT_STATUS.nok) && 'retry';
  const fail = retry && (data.retryTimes >= data.points[data.activeIndex].maxRetryTimes) && 'fail';
  const finish = (!retry && (data.activeIndex + results.length >= data.points.length)) && 'finish';
  const next = !retry && !finish && 'next';
  CommonLog.Info([LSN, fail, retry, finish, next]);
  return [LSN, fail, retry, finish, next];
};

const resultStatusTasks = (ORDER, orderActions, results: Array<tResult>) => ({
  * retry() {
    try {
      yield call(this.updateData,(d: tScrewStepData): tScrewStepData => ({
        ...d,
        points: mergePointsAndResults(d.points, results, d.activeIndex, d.points[d.activeIndex]?.group_sequence),
        retryTimes: (d.retryTimes || 0) + 1
      }));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * fail() {
    try {
      yield call(this.updateData,(d: tScrewStepData): tScrewStepData => ({
        ...d,
        activeIndex: -1,
        points: mergePointsAndResults(d.points, results, d.activeIndex, d.points[d.activeIndex]?.group_sequence)
      }));
      yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * finish() {
    try {
      yield call(this.updateData,(d: tScrewStepData): tScrewStepData => ({
        ...d,
        activeIndex: -1,
        points: mergePointsAndResults(d.points, results, d.activeIndex, d.points[d.activeIndex]?.group_sequence)
      }));
      yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * next() {
    try {
      // update step data
      yield call(this.updateData,(d: tScrewStepData): tScrewStepData => ({
        ...d,
        activeIndex: d.activeIndex === -1 ? 0 : d.activeIndex + results.length,
        points: mergePointsAndResults(
          d.points,
          results,
          d.activeIndex,
          d.activeIndex === -1 ?
            d.points[0].group_sequence :
            d.points[d.activeIndex + results.length]?.group_sequence
        )
      }));
    } catch (e) {
      CommonLog.lError(e);
    }
  },
  * LSN() {
    try {
      // do nothing
    } catch (e) {
      CommonLog.lError(e);
    }
  }
});

export default function* handleResult(ORDER, orderActions, results, data) {
  try {
    yield call(this.updateData,(d: tScrewStepData): tScrewStepData => ({
      ...d,
      timeLine: [
        ...results.map(r => ({
          title: r.batch,
          color: r.result === 'ok' ? 'info' : 'danger',
          footerTitle: r.tool_sn,
          body: `${r.result}: wi=${r.wi},mi=${r.mi},ti=${r.ti}`
        })),
        ...(d.timeLine || [])
      ]
    }));
    const firstMatchResultStatus =
      resultStatusTasks(
        ORDER,
        orderActions,
        results
      )[resultStatus(results, data).find(v => !!v)];
    // 执行
    if (firstMatchResultStatus) {
      yield call([this, firstMatchResultStatus]);
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}
