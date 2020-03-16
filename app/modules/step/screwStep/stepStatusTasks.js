import { actionChannel, all, call, delay, put, select, take } from 'redux-saga/effects';
import { isEmpty, isNil } from 'lodash-es';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import { STEP_STATUS } from '../constants';
import { orderActions } from '../../order/action';
import { CommonLog } from '../../../common/utils';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './interface/typeDef';
import { ClsOrderOperationPoints } from './classes/ClsOrderOperationPoints';
import { workModes } from '../../workCenterMode/constants';
import { tNS } from '../../../i18n';
import { reworkDialogConstants as dia, reworkNS } from '../../reworkPattern/constants';
import { getDevice, getDevicesByType } from '../../deviceManager/devices';
import { deviceType } from '../../deviceManager/constants';
import SelectCard from '../../../components/SelectCard';
import screwStepActions from './action';
import dialogActions from '../../dialog/action';
import { controllerModes, SCREW_STEP } from './constants';
import { getPestListApi } from '../../../api/tools';
import type { ClsOperationPoint } from './classes/ClsOperationPoint';
import { getTools } from './getTools';
import { byPassPoint } from './byPassPoint';
import { stepDataApi } from '../../../api/order';
import type { IWorkable } from '../../workable/IWorkable';
import { getResult } from './getResult';
import { disableTools, setTools } from './setTools';

function* enteringState(config) {
  try {
    // init data
    // check payload
    const payload: tScrewStepPayload = this._payload;
    if (isNil(payload) || isEmpty(payload)) {
      throw new Error(
        `ScrewStepPayload Is Empty! code: ${
          this._id
        }, Payload: ${JSON.stringify(payload)}`
      );
    }
    const points = payload.tightening_points;
    this._pointsManager = new ClsOrderOperationPoints(
      payload.tightening_points
    );

    const { workcenterType } = yield select(s => s.setting.system.workcenter);
    const { workCenterMode } = yield select();

    switch (workCenterMode) {
      case workModes.reworkWorkCenterMode: {
        // is rework workcenter
        // TODO judge by checking if the order belongs to current workcenter
        if (workcenterType === 'rework') {
          // do not check pset/job, set pest manually
          const btnCancel = {
            label: tNS(dia.cancel, reworkNS),
            color: 'warning',
            action: orderActions.stepStatus(this, STEP_STATUS.FAIL)
          };
          const tools = getDevicesByType(deviceType.tool);
          const ToolSelectCard = SelectCard(screwStepActions.selectTool);

          yield put(dialogActions.dialogShow({
            buttons: [btnCancel],
            title: '选择工具',
            content: (<Grid spacing={1} container>
              {tools.map(t => <Grid item xs={6} key={`${t.serialNumber}`}>
                <ToolSelectCard
                  name={t.Name}
                  status={t.Healthz ? '已连接' : '已断开'}
                  infoArr={[t.serialNumber]}
                  height={130}
                  item={t}
                />
              </Grid>)}
            </Grid>)
          }));
          const { tool } = yield take(SCREW_STEP.SELECT_TOOL);
          this._tools = [tool];
          this._forceTool = tool;
          yield put(dialogActions.dialogClose());
          const PsetSelect = SelectCard(screwStepActions.selectPset);

          // TODO get tool psets
          const psets = (yield call(getPestListApi, tool.serialNumber))?.data?.pset_list || [];
          yield delay(300);
          yield put(dialogActions.dialogShow({
            buttons: [btnCancel],
            title: '选择PSET',
            content: (<Grid spacing={1} container>
              {psets.map(p => <Grid item xs={6} key={`${p}`}>
                <PsetSelect
                  name={p}
                  height={130}
                  item={p}
                />
              </Grid>)}
            </Grid>)
          }));
          const { pset } = yield take(SCREW_STEP.SELECT_PSET);
          this._forcePset = pset;
          yield put(dialogActions.dialogClose());
        } else {
          const { reworkConfig } = config || {};
          const { point } = reworkConfig;
          this._tools = yield call(getTools, [point] || []);
          this._forcePset = null;
          this._forceTool = null;
          yield call(
            this.updateData,
            (data: tScrewStepData): tScrewStepData => ({
              ...data,
              controllerMode: controllerModes.pset
            })
          );
        }
        break;
      }
      case workModes.normWorkCenterMode: {
        if (!isNil(payload.jobID)) {
          yield call(
            this.updateData,
            (data: tScrewStepData): tScrewStepData => ({
              ...data,
              controllerMode: controllerModes.job
            })
          );
          points.reduce((tSN: string, p: tPoint): string => {
            if (tSN && p.tightening_tool !== tSN) {
              throw new Error('工具序列号不匹配');
            }
            return p.tightening_tool || tSN || '';
          }, null);
        } else if (points.every(p => !isNil(p.pset))) {
          yield call(
            this.updateData,
            (data: tScrewStepData): tScrewStepData => ({
              ...data,
              controllerMode: controllerModes.pset
            })
          );
        } else {
          throw new Error('缺少JOB号或Pset号');
        }
        this._tools = yield call(getTools, payload?.tightening_points || []);
        this._forcePset = null;
        this._forceTool = null;

        break;
      }
      default:
        break;
    }
    // eslint-disable-next-line camelcase
    this._tools.forEach(t => {
      this._listeners.push(
        t.addListener(
          () => true,
          input => screwStepActions.result(input.data)
        )
      );
    });
    if (this._data && this._data.results) {
      if (!(window.debugSettings && window.debugSettings.disableOrderTriggerLimit)) {
        this._pointsManager.newResult(this._data.results);
      }
      this._pointsManager.stop();
    }

    yield call(
      this.updateData,
      (data: tScrewStepData): tScrewStepData => ({
        ...data,
        tightening_points: this._pointsManager.points.map(p => p.data)
      })
    );
    // goto doing
    yield put(orderActions.stepStatus(this, STEP_STATUS.DOING, config));
  } catch (e) {
    CommonLog.lError(e, { at: 'screwStep ENTERING' });
    yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: e }));
  }
}

function* doingState(config) {
  try {
    const { reworkConfig } = config || {};
    const { workCenterMode } = yield select();
    let redoPointClsObj = null;

    switch (workCenterMode) {
      case workModes.reworkWorkCenterMode: {
        const { point } = reworkConfig;
        if (point) {
          redoPointClsObj = this.points.find(p => p.sequence === point.sequence);
        } else {
          redoPointClsObj = this.points.find(p => !p.noRedo);
        }
        if (redoPointClsObj) {
          this._pointsToActive = [redoPointClsObj.start(true)];
        }
        break;
      }
      case workModes.normWorkCenterMode: {
        if (!this._pointsManager) {
          throw new Error('拧紧点异常');
        }
        break;
      }
      default:
        break;
    }

    const resultChannel = yield actionChannel([SCREW_STEP.RESULT]);
    const { controllerMode, jobID } = this._data;
    let isFirst = true; // job只设置一次，记录状态

    // let test = true; // job只设置一次，记录状态


    while (true) {
      switch (workCenterMode) {
        case workModes.reworkWorkCenterMode: {
          if (this._pointsToActive && this._pointsToActive.length > 0) {
            break;
          }
          const canRedoPoint = this.points.find(p => !p.noRedo);
          const success = redoPointClsObj
            && redoPointClsObj.isSuccess
            && !canRedoPoint;
          if (success) {
            yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED)); // 成功退出
          } else {
            yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { silent: true })); // 失败退出
          }
          break;
        }
        case workModes.normWorkCenterMode: {
          const finalFailPoints = (this._newInactivePoints || []).filter(
            (p: ClsOperationPoint) => p.isFinalFail
          );
          yield call([this, byPassPoint], finalFailPoints);

          // if (test){
          // yield  put(screwStepActions.byPassSpecPoint());
          //   test=false;
          // }

          const newActivePoints = this._pointsManager.start();
          this._pointsToActive = newActivePoints.filter(p =>
            this._pointsToActive.every(pp => pp.sequence !== p.sequence)
          );
          if (
            this._pointsManager.isFailed &&
            this._pointsManager.points.filter(p => p.isActive).length === 0
          ) {
            yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: '拧紧失败' })); // 失败退出
          } else if (this._pointsManager.isPass) {
            yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED)); // 成功退出
          }
          break;
        }
        default:
          break;
      }

      yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
        ...data,
        tightening_points: this._pointsManager.points.map(p => p.data)
      }));

      const activeConfigs = this._pointsToActive.map(p => {
        const cModeId = controllerMode === controllerModes.job ? jobID : this._forcePset || p.pset;
        const tool = this._forceTool || getDevice(p.toolSN);
        return {
          point: p,
          tool,
          controllerModeId: cModeId
        };
      });
      yield call([this, setTools],
        activeConfigs,
        controllerMode,
        isFirst
      );
      const results = yield call(
        [this, getResult],
        resultChannel
      );
      this._pointsToActive = [];

      this._newInactivePoints = this._pointsManager.newResult(results);
      // disable tools before bypass point
      yield call([this, disableTools], this._newInactivePoints);
      yield call(this.updateData, (data: tScrewStepData): tScrewStepData => ({
        ...data,
        tightening_points: this._pointsManager.points.map(p => p.data)
      }));

      if (isFirst) {
        isFirst = false;
      }
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'screwStep DOING' });
    yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: e }));
  }
}

function* failState(config) {
  try {
    const { error, silent } = config;
    yield all(this._tools.map(t => call(t.Disable)));
    yield call([this, clearStepData]);
    if (!silent) {
      const { workCenterMode } = yield select();
      const isNormal = workCenterMode === workModes.normWorkCenterMode;
      let buttons = [
        {
          label: 'Common.Close',
          color: 'danger',
          action: screwStepActions.confirmFail()
        }
      ];
      if (isNormal) {
        buttons = buttons.concat([
          {
            label: 'Order.Retry',
            color: 'info',
            action: orderActions.stepStatus(this, STEP_STATUS.READY)
          },
          {
            label: 'Order.Next',
            color: 'warning',
            action: screwStepActions.confirmFail()
          }
        ]);
      }
      yield put(
        dialogActions.dialogShow({
          buttons,
          title: `工步失败：${this._code}`,
          content: `${error || this.failureMsg}`
        })
      );
      yield take(SCREW_STEP.CONFIRM_FAIL);
    }
    yield put(orderActions.finishStep(this));
  } catch (e) {
    CommonLog.lError(e, { at: 'screwStep FAIL' });
    yield put(orderActions.finishStep(this));
  }
}

export const screwStepStatusTasksMixin = (superTasks) => ({
  ...superTasks,
  [STEP_STATUS.ENTERING]: enteringState,
  [STEP_STATUS.DOING]: doingState,
  [STEP_STATUS.FAIL]: failState
});

function* clearStepData() {
  try {
    if (this._pointsManager) {
      this._pointsManager.stop();
      this._pointsManager.clearByPass();
    }
    if (this._tools) {
      yield all(
        this._tools.map(t => (t.isEnable ? call(t.Disable) : call(() => {
        })))
      );
    }
    if (this._tools && this._listeners) {
      this._tools.forEach(t => {
        this._listeners.forEach(l => {
          t.removeListener(l);
        });
      });
      this._tools = [];
      this._listeners = [];
    }
    this._pointsToActive = [];

    CommonLog.Info('tools cleared', {
      at: `screwStep(${String((this: IWorkable)._code)})._onLeave`
    });
  } catch (e) {
    CommonLog.lError(e, {
      at: `screwStep(${String((this: IWorkable)._code)}).clearStepData`
    });
  }

}

export function* onLeave() {
  try {
    yield call(stepDataApi, this.code, this._data);
    yield call([this, clearStepData]);
  } catch (e) {
    CommonLog.lError(e, {
      at: `screwStep(${String((this: IWorkable)._code)})._onLeave`
    });
  }
}
