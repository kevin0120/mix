import { all, call } from 'redux-saga/effects';
import { controllerModes } from './constants';
import controllerModeTasks from './controllerModeTasks';
import { CommonLog } from '../../../common/utils';
import type { tScrewStepData } from './interface/typeDef';
import { getDevice } from '../../deviceManager/devices';

export function* setTools(activeControls, controllerMode, isFirst) {
  if (controllerMode === controllerModes.job && !isFirst) {
    return;
  }
  if (!controllerModeTasks[controllerMode]) {
    throw new Error(`未识别的控制器模式:${controllerMode}`);
  }
  try {
    let successCount = 0;
    for (const c of activeControls) {
      successCount += yield call([this, setSingleTool], controllerMode, c);
    }
    return successCount;
  } catch (e) {
    CommonLog.lError(e);
    // throw e;
  }
}

function* setSingleTool(controllerMode, singleControl) {
  try {
    const { sequence, tool, controllerModeId } = singleControl;
    yield call([this, controllerModeTasks[controllerMode]], sequence, tool, controllerModeId);
    yield call(tool?.Enable || (() => {
      throw new Error(
        `tool ${tool?.Name}: no such tool or tool cannot be enabled.`
      );
    }));
    return 1;
  } catch (e) {
    yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
      ...data,
      tightening_points: this._pointsManager.points.map(p => p.data)
    }));
    return 0;
  }
}

export function* disableTools(controls) {
  yield all(controls.map(c => call([this, disableSingleTool], c)));
}

function* disableSingleTool(control) {
  try {
    yield call(getDevice(control.toolSN)?.Disable || (() => {
      throw new Error(
        `tool ${control.toolSN}: no such tool or tool cannot be disabled.`
      );
    }));
  } catch (e) {
    CommonLog.lError(e);
    yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
      ...data,
      tightening_points: this._pointsManager.points.map(p => p.data)
    }));
  }
}
