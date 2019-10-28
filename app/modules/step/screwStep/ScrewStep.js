// @flow
import type { Saga } from 'redux-saga';
import { cloneDeep, isNil, isEmpty } from 'lodash-es';
import { call, put, take, all, actionChannel } from 'redux-saga/effects';
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
import { getDevice } from '../../external/device';
import dialogActions from '../../dialog/action';
import type { IWorkStep } from '../interface/IWorkStep';
import type { IScrewStep } from './interface/IScrewStep';
import { reduceResult2TimeLine } from './handleResult';
import type { ClsOperationPoint } from './classes/ClsOperationPoint';

export function* doPoint(
  points: Array<ClsOperationPoint>,
  isFirst: boolean,
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
      const success = yield call(
        [this, controllerModeTasks[data.controllerMode]],
        p
      );
      if (!success) {
        throw new Error(`${data.controllerMode} failed`);
      }
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'doPoint', points });
    yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
  }
}

export function getTools(points: Array<tPoint>) {
  const toolSnSet = new Set(points.map(p => p.toolSN));
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

    _activePoints = [];

    *_onLeave() {
      try {
        yield all(
          this._tools.map(t => (t.isEnable ? call(t.Disable) : call(() => {})))
        );
        this._tools = [];
        CommonLog.Info('tools cleared', {
          at: `screwStep(${(this: IWorkStep)._name},${
            (this: IWorkStep)._id
          })._onLeave`
        });
      } catch (e) {
        CommonLog.lError(e, {
          at: `screwStep(${(this: IWorkStep)._name},${
            (this: IWorkStep)._id
          })._onLeave`
        });
      }
    }

    // eslint-disable-next-line flowtype/no-weak-types
    constructor(...args: Array<any>) {
      super(...args);
      this.isValid = true; // 设置此工步是合法的
    }

    _statusTasks = {
      *[STEP_STATUS.ENTERING](ORDER, orderActions) {
        try {
          // init data
          const payload: tScrewStepPayload = this._payload;
          if (isNil(payload) || isEmpty(payload)) {
            throw new Error(
              `ScrewStepPayload Is Empty! id: ${
                this._id
              }, Payload: ${JSON.stringify(payload)}`
            );
          }

          this._pointsManager = new ClsOrderOperationPoints(payload.points);
          const points: Array<tPoint> = cloneDeep(payload?.points || []);

          this._tools = yield call(getTools, points);

          yield call(
            this.updateData,
            (data: tScrewStepData): tScrewStepData => ({
              ...data,
              points: this._pointsManager.points
            })
          );
          yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep ENTERING' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
        }
      },

      *[STEP_STATUS.DOING](ORDER, orderActions) {
        try {
          let isFirst = true;
          this._activePoints = this._pointsManager.start();

          const resultChannel = yield actionChannel([
            SCREW_STEP.RESULT,
            SCREW_STEP.REDO_POINT
          ]);
          while (true) {
            yield call(
              this.updateData,
              (data: tScrewStepData): tScrewStepData => ({
                ...data,
                points: this._pointsManager.points // results data.results
              })
            );

            if (this._pointsManager.isPass()) {
              yield put(orderActions.stepStatus(this, STEP_STATUS.FINISHED));
            }

            if (this._pointsManager.isFailed()) {
              yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL));
            }

            if (this._activePoints && this._activePoints.length > 0) {
              yield call(
                [this, doPoint],
                [...this._activePoints],
                isFirst,
                orderActions
              );
            }

            yield all(
              this._activePoints.map(p =>
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

            switch (action.type) {
              case SCREW_STEP.RESULT: {
                const {
                  results: { data: results }
                } = action;
                yield call(this.updateData, reduceResult2TimeLine(results));
                const { active, inactive } = this._pointsManager.newResult(
                  results
                );
                this._activePoints = active;
                yield all(
                  inactive.map(p =>
                    call(
                      getDevice(p.toolSN)?.Disable ||
                        (() => {
                          CommonLog.lError(
                            `tool ${p.toolSN}: no such tool or tool cannot be disabled.`
                          );
                        })
                    )
                  )
                );
                break;
              }
              case SCREW_STEP.REDO_POINT: {
                const { point } = action;
                if (!this._pointsManager.hasPoint(point)) {
                  break;
                }
                this._activePoints = [point.redo()];
                break;
              }
              default:
                break;
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

      *[STEP_STATUS.FINISHED](ORDER, orderActions) {
        try {
          yield put(orderActions.finishStep(this));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep FINISHED' });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, e));
        }
      },

      *[STEP_STATUS.FAIL](ORDER, orderActions, msg) {
        try {
          yield all(this._tools.map(t => call(t.Disable)));
          this._tools = [];
          yield put(
            dialogActions.dialogShow({
              buttons: [
                {
                  label: 'Common.Close',
                  color: 'danger'
                },
                {
                  label: 'Order.Next',
                  color: 'warning',
                  action: screwStepActions.confirmFail()
                }
              ],
              title: `工步失败：${this._name}`,
              content: `${msg || ''}`
            })
          );
          yield take(SCREW_STEP.CONFIRM_FAIL);
          yield put(orderActions.finishStep(this));
        } catch (e) {
          CommonLog.lError(e, { at: 'screwStep FAIL' });
          yield put(orderActions.finishStep(this));
        }
      }
    };
  };
export default ScrewStepMixin;
