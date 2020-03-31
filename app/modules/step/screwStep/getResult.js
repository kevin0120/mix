import { call, take } from 'redux-saga/effects';
import { CommonLog } from '../../../common/utils';
import { result2TimeLine } from './timeLine';

export function* getResult(resultChannel) {
  try {
    const action = yield take(resultChannel);
    const { results } = action;
    yield call(this.updateData, data => ({
      ...data,
      results: [...(data.results || []), ...results],
      timeLine: [
        ...result2TimeLine(results),
        ...(data.timeLine || [])
      ]
    }));
    return results;
  } catch (e) {
    CommonLog.lError(e);
    throw e;
  }
}
