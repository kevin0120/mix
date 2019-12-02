// @flow
import type { Saga } from 'redux-saga';
import { isNil, isEmpty } from 'lodash-es';
import { call, put, take, all, actionChannel, race, select } from 'redux-saga/effects';
import { STEP_STATUS } from '../constants';
import type {
  tPoint,
  tScrewStepData,
  tScrewStepPayload
} from './interface/typeDef';
import { ClsOrderOperationPoints } from './classes/ClsOrderOperationPoints';
import { CommonLog } from '../../../common/utils';
import controllerModeTasks from './controllerModeTasks';
import screwStepActions from './action';
import { SCREW_STEP, controllerModes } from './constants';
import { getDevice } from '../../deviceManager/devices';
import dialogActions from '../../dialog/action';
import type { IWorkStep } from '../interface/IWorkStep';
import type { IScrewStep } from './interface/IScrewStep';
import { result2TimeLine } from './timeLine';
import type { ClsOperationPoint } from './classes/ClsOperationPoint';
import { stepDataApi } from '../../../api/order';
import type { IWorkable } from '../../workable/IWorkable';
import { workModes } from '../../workCenterMode/constants';

export function* doPoint(
  points: Array<ClsOperationPoint>,
  isFirst: boolean,
  // eslint-disable-next-line flowtype/no-weak-types
  orderActions: Object
): Saga<void> {
  try {
    const data = this._data;

    if (
      !(
        data.controllerMode === controllerModes.pset ||
        (data.controllerMode === controllerModes.job && isFirst)
      )
    ) {
      return;
    }
    // eslint-disable-next-line no-restricted-syntax
    for (const p of points) {
      yield call([this, controllerModeTasks[data.controllerMode]], p.point);
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'doPoint', points });
    yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
  }
}

function* getResult(pointsToActive, resultChannel, isFirst, orderActions) {
  try {
    if (pointsToActive && pointsToActive.length > 0) {
      yield call(
        [this, doPoint],
        [...pointsToActive],
        isFirst,
        orderActions
      );
    }

    // 先设置pset再enable
    yield all(
      pointsToActive.map(p =>
        call(
          getDevice(p.toolSN)?.Enable ||
          (() => {
            CommonLog.lError(
              `tool ${p.toolSN}: no such tool or tool cannot be enabled.`
            );
          })
        )
      )
    );

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
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL)); // 失败退出
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

    * _onLeave() {
      try {
        yield call(
          this.updateData,
          (data: tScrewStepData): tScrewStepData => ({
            ...data,
            tightening_points: this._pointsManager.points.map(p => p.data)
          })
        );
        yield call(stepDataApi, this.id, this._data);
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

    _statusTasks = {
      * [STEP_STATUS.READY](ORDER, orderActions, config) {
        try {
          console.warn(config);
          yield put(orderActions.stepStatus(this, STEP_STATUS.ENTERING, config));
        } catch (e) {
          CommonLog.lError(e);
        }
      },
      * [STEP_STATUS.ENTERING](ORDER, orderActions, config) {
        try {
          // init data
          const payload: tScrewStepPayload = this._payload;
          if (isNil(payload) || isEmpty(payload)) {
            throw new Error(
              `ScrewStepPayload Is Empty! code: ${
                this._id
              }, Payload: ${JSON.stringify(payload)}`
            );
          }

          const points = payload.tightening_points;

          if (!isNil(payload.jobID)) {
            yield call(
              this.updateData,
              (data: tScrewStepData): tScrewStepData => ({
                ...data,
                controllerMode: controllerModes.job
              })
            );
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

          this._pointsManager = new ClsOrderOperationPoints(
            payload.tightening_points
          );

          // eslint-disable-next-line camelcase
          this._tools = yield call(getTools, payload?.tightening_points || []);
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
          yield put(orderActions.stepStatus(this, STEP_STATUS.DOING, config));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep ENTERING' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
        }
      },
      * [STEP_STATUS.DOING](ORDER, orderActions, config) {
        try {
          const { reworkConfig } = config || {};
          const { workCenterMode } = yield select();
          let redoPointClsObj = null;
          if (workCenterMode === workModes.reworkWorkCenterMode) {
            const { point } = reworkConfig;
            if (point) {
              redoPointClsObj = this.points.find(p => p.sequence === point.sequence);
            } else {
              redoPointClsObj = this.points.find(p => p.canRedo);
            }
            console.warn('rework step with point', redoPointClsObj);
            if (redoPointClsObj) {
              this._pointsToActive = [redoPointClsObj.start()];
            }
          } else {
            if (!this._pointsManager) {
              throw new Error('拧紧点异常');
            }
            this._pointsToActive = this._pointsManager.start();
          }


          let isFirst = true; // job只设置一次，记录状态
          const resultChannel = yield actionChannel([SCREW_STEP.RESULT]);
          while (true) {
            yield call(
              this.updateData,
              (data: tScrewStepData): tScrewStepData => ({
                ...data,
                tightening_points: this._pointsManager.points.map(p => p.data)
              })
            );
            console.warn(this._pointsToActive);

            const results = yield call(
              [this, getResult],
              this._pointsToActive,
              resultChannel,
              isFirst,
              orderActions
            );

            const { inactive } = this._pointsManager.newResult(results);
            // disable tools before bypass point
            yield all(inactive.map(p => call(
              getDevice(p.toolSN)?.Disable || (() => {
                CommonLog.lError(
                  `tool ${p.toolSN}: no such tool or tool cannot be disabled.`
                );
              })
            )));

            if (workCenterMode === workModes.reworkWorkCenterMode) {
              const canRedoPoint = this.points.find(p => p.canRedo);
              const success = redoPointClsObj
                && redoPointClsObj.isSuccess
                && !canRedoPoint;
              if (success) {
                yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED)); // 成功退出
              } else {
                yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL)); // 失败退出
              }
            } else {
              const finalFailPoints = inactive.filter(
                (p: ClsOperationPoint) => p.isFinalFail
              );
              yield call([this, byPassPoint], finalFailPoints, orderActions);

              this._pointsToActive = this._pointsManager.start();
              if (
                this._pointsManager.isFailed &&
                this._pointsManager.points.filter(p => p.isActive).length === 0
              ) {
                yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL)); // 失败退出
              } else if (this._pointsManager.isPass) {
                yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED)); // 成功退出
              }
            }


            if (isFirst) {
              isFirst = false;
            }
          }
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep DOING' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
        }
      },

      * [STEP_STATUS.FINISHED](ORDER, orderActions) {
        try {
          yield put(orderActions.finishStep(this));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep FINISHED' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
        }
      },

      * [STEP_STATUS.FAIL](ORDER, orderActions, msg) {
        try {
          yield all(this._tools.map(t => call(t.Disable)));
          const { workCenterMode } = yield select();
          const isNormal = workCenterMode === workModes.normWorkCenterMode;
          if (isNormal) {
            yield put(
              dialogActions.dialogShow({
                buttons: [
                  {
                    label: 'Common.Close',
                    color: 'danger',
                    action: screwStepActions.confirmFail()
                  },
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
                ],
                title: `工步失败：${this._code}`,
                content: `${msg || this.failureMsg}`
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
  };
export default ScrewStepMixin;
