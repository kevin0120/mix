import { all, call } from 'redux-saga/effects';
import { controllerModes } from './constants';
import controllerModeTasks from './controllerModeTasks';
import { CommonLog } from '../../../common/utils';
import type { tScrewStepData } from './interface/typeDef';
import { getDevice } from '../../deviceManager/devices';
import { byPassPoint } from './byPassPoint';


export function* setTools(activeConfigs, controllerMode, isFirst) {
  if (controllerMode === controllerModes.job && !isFirst) {
    return;
  }
  if (!controllerModeTasks[controllerMode]) {
    throw new Error(`未识别的控制器模式:${controllerMode}`);
  }
  try {
    const effects = activeConfigs.map(c => {
      const { point, tool, controllerModeId } = c;
      return call([this, setSingleTool], controllerMode, point, tool, controllerModeId);
    });
    yield all(effects);
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* setSingleTool(controllerMode, point, tool, controllerModeId) {
  try {
    yield call([this, controllerModeTasks[controllerMode]], point.point, tool, controllerModeId);
    yield call(tool?.Enable || (() => {
      CommonLog.lError(
        `tool ${tool?.Name}: no such tool or tool cannot be enabled.`
      );
    }));
  } catch (e) {
    CommonLog.lError(e);
    yield call([this, byPassPoint], [point]);
    yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
      ...data,
      tightening_points: this._pointsManager.points.map(p => p.data)
    }));
  }
}

export function* disableTools(pointsToDisable) {
  yield all(pointsToDisable.map(p => call([this, disableSingleTool], p)));
}

function* disableSingleTool(point) {
  try {
    yield call(getDevice(point.toolSN)?.Disable || (() => {
      CommonLog.lError(
        `tool ${point.toolSN}: no such tool or tool cannot be disabled.`
      );
    }));
  } catch (e) {
    CommonLog.lError(e);
    yield call([this, byPassPoint], [point]);
    yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
      ...data,
      tightening_points: this._pointsManager.points.map(p => p.data)
    }));
  }
}
