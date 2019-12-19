// @flow
import { isEmpty, isNil } from 'lodash-es';
import { actionChannel, all, call, delay, put, race, select, take } from 'redux-saga/effects';
import React from 'react';
import Grid from '@material-ui/core/Grid';
import { STEP_STATUS } from '../constants';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './interface/typeDef';
import { ClsOrderOperationPoints } from './classes/ClsOrderOperationPoints';
import { CommonLog } from '../../../common/utils';
import controllerModeTasks from './controllerModeTasks';
import screwStepActions from './action';
import { controllerModes, SCREW_STEP } from './constants';
import { getDevice, getDevicesByType } from '../../deviceManager/devices';
import dialogActions from '../../dialog/action';
import type { IWorkStep } from '../interface/IWorkStep';
import type { IScrewStep } from './interface/IScrewStep';
import { result2TimeLine } from './timeLine';
import type { ClsOperationPoint } from './classes/ClsOperationPoint';
import { stepDataApi } from '../../../api/order';
import type { IWorkable } from '../../workable/IWorkable';
import { workModes } from '../../workCenterMode/constants';
import { reworkDialogConstants as dia, reworkNS } from '../../reworkPattern/constants';
import { tNS } from '../../../i18n';
import { deviceType } from '../../deviceManager/constants';
import SelectCard from '../../../components/SelectCard';
import { getPestListApi } from '../../../api/tools';

function* getResult(activeConfigs, resultChannel, controllerMode, isFirst) {
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

function* byPassPoint(finalFailPoints, orderActions) {
  try {
    const n: string = finalFailPoints
      .map((p: ClsOperationPoint) => p.point.nut_no)
      .join(',');
    if (finalFailPoints.length > 0) {
      CommonLog.Debug('Show Next Point By Pass Diag');
      yield put(
        dialogActions.dialogShow({
          buttons: [
            {
              label: 'Order.Next',
              color: 'danger',
              action: screwStepActions.confirmFailSpecPoint()

            },
            {
              label: 'Screw.Next',
              color: 'warning',
              action: screwStepActions.byPassSpecPoint()
            }
          ],
          // eslint-disable-next-line camelcase
          title: `拧紧点失败：${n}`,
          content: `${this.failureMsg}`
        })
      );
      const { bypass, fail } = yield race({
        bypass: take(SCREW_STEP.BYPASS_SPEC_POINT),
        fail: take(SCREW_STEP.CONFIRM_FAIL_SPEC_POINT)
      });
      if (fail) {
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: '拧紧失败' })); // 失败退出
      }
    }
  } catch (e) {
    CommonLog.lError(e);
    throw e;
  }
}

export function getTools(points: Array<tPoint>) {
  const toolSnSet = new Set(points.map(p => p.tightening_tool));
  const lostTool = [];
  const tools = [];
  toolSnSet.forEach(t => {
    const tool = getDevice(t);
    if (tool) {
      tools.push(tool);
    } else {
      lostTool.push(t);
    }
  });

  if (lostTool.length > 0) {
    throw new Error(`tools not found: ${String(lostTool.map(t => `${t}`))}`);
  }

  const unhealthyTools = tools.filter(t => !t.Healthz);
  if (unhealthyTools.length > 0) {
    throw new Error(
      `tool not connected: ${JSON.stringify(
        unhealthyTools.map(t => `${String(t.serialNumber)}`)
      )}`
    );
  }
  return tools;
}

const ScrewStepMixin = (ClsBaseStep: Class<IWorkStep>) =>
  class ClsScrewStep extends ClsBaseStep implements IScrewStep {
    _tools = [];

    isValid: boolean = false;

    _orderOperationPoints: ClsOrderOperationPoints;

    _pointsToActive = [];

    _listeners = [];

    _statusTasks = {
      * [STEP_STATUS.READY](ORDER, orderActions, config) {
        try {
          yield put(orderActions.stepStatus(this, STEP_STATUS.ENTERING, config));
        } catch (e) {
          CommonLog.lError(e);
        }
      },
      * [STEP_STATUS.ENTERING](ORDER, orderActions, config) {
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
                console.log('tool selected', tool);
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
                this._tools = yield call(getTools, payload?.tightening_points || []);
                this._forcePset = null;
                this._forceTool = null;
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
            this._pointsManager.newResult(this._data.results);
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
      },
      * [STEP_STATUS.DOING](ORDER, orderActions, config) {
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
                yield call([this, byPassPoint], finalFailPoints, orderActions);
                this._pointsToActive = this._pointsManager.start();
                console.warn(this._pointsToActive);
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
            const results = yield call(
              [this, getResult],
              activeConfigs,
              resultChannel,
              controllerMode,
              isFirst
            );
            this._pointsToActive = [];

            this._newInactivePoints = this._pointsManager.newResult(results);
            // disable tools before bypass point
            yield all(this._newInactivePoints.map(p => call(
              getDevice(p.toolSN)?.Disable || (() => {
                CommonLog.lError(
                  `tool ${p.toolSN}: no such tool or tool cannot be disabled.`
                );
              })
            )));
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
          console.log(e);
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: e }));
        }
      },

      * [STEP_STATUS.FINISHED](ORDER, orderActions) {
        try {
          yield put(orderActions.finishStep(this));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep FINISHED' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: e }));
        }
      },

      * [STEP_STATUS.FAIL](ORDER, orderActions, config) {
        try {
          const { error, silent } = config;
          yield all(this._tools.map(t => call(t.Disable)));
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
    };

    // eslint-disable-next-line flowtype/no-weak-types
    constructor(...args: Array<any>) {
      super(...args);
      this.isValid = true; // 设置此工步是合法的
      this._onLeave = this._onLeave.bind(this);
    }

    get points() {
      if (this._pointsManager) {
        return this._pointsManager.points;
      }
      return [];
    }

    * _onLeave() {
      try {
        yield call(stepDataApi, this.id, this._data);
        this._pointsManager.stop();
        yield all(
          this._tools.map(t => (t.isEnable ? call(t.Disable) : call(() => {
          })))
        );
        this._tools.forEach(t => {
          this._listeners.forEach(l => {
            t.removeListener(l);
          });
        });

        this._tools = [];
        this._listeners = [];
        CommonLog.Info('tools cleared', {
          at: `screwStep(${String((this: IWorkable)._code)})._onLeave`
        });
      } catch (e) {
        CommonLog.lError(e, {
          at: `screwStep(${String((this: IWorkable)._code)})._onLeave`
        });
      }
    }
  };
export default ScrewStepMixin;
