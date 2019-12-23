import { all, call, take } from 'redux-saga/effects';
import { controllerModes } from './constants';
import controllerModeTasks from './controllerModeTasks';
import { CommonLog } from '../../../common/utils';
import { result2TimeLine } from './timeLine';

export function* getResult(activeConfigs, resultChannel, controllerMode, isFirst) {
  try {
    // if (!activeConfigs || activeConfigs.length === 0) {
    //   return;
    // }
    if (controllerMode === controllerModes.job && !isFirst) {
      return;
    }
    if (!controllerModeTasks[controllerMode]) {
      throw new Error(`未识别的控制器模式:${controllerMode}`);
    }
    const effects = activeConfigs.map(c => {
      const { point, tool, controllerModeId } = c;
      return call([this, controllerModeTasks[controllerMode]], point.point, tool, controllerModeId);
    });
    if (effects && effects.length > 0) {
      yield all(effects);
    }

    // 先设置pset再enable
    yield all(activeConfigs.map(c => {
      const { tool } = c;
      return call(tool?.Enable || (() => {
        CommonLog.lError(
          `tool ${tool?.Name}: no such tool or tool cannot be enabled.`
        );
      }));
    }));

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
