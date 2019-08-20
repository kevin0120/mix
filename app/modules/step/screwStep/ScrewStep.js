import Step from '../Step';
import STEP_STATUS from '../model';
import type { tPoint, tScrewStepData, tScrewStepPayload } from './model';
import { call, put, take } from 'redux-saga/effects';
import { cloneDeep } from 'lodash-es';
import { CommonLog } from '../../../common/utils';
import handleResult from './handleResult';
import controllerModeTasks from './controllerModeTasks';
import screwStepActions, { SCREW_STEP } from './action';
import { getDevice } from '../../external/device';
import dialogActions from '../../dialog/action';
import i18n from '../../../i18n';

function* doPoint(point, isFirst, orderActions) {
  try {
    const data = this._data;
    if (data.controllerMode === 'pset' || (data.controllerMode === 'job' && isFirst)) {
      const success = yield call([this, controllerModeTasks[data.controllerMode]], point);
      if (!success) {
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL,`${data.controllerMode} failed`));
      }
      if (isFirst) {
        for (const t of this._tools) {
          yield call(t.Enable);
        }
      }
    }
    return yield take([SCREW_STEP.RESULT, SCREW_STEP.REDO_POINT]);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'doPoint',
      point
    });
  }
}

export default class ScrewStep extends Step {
  _tools = [];

  _statusTasks = {
    * [STEP_STATUS.ENTERING](ORDER, orderActions) {
      try {
        // init data
        const payload: tScrewStepPayload = this._payload;
        const points: Array<tPoint> = cloneDeep(payload?.points || [])
          .sort((a, b) => a.group_sequence - b.group_sequence);
        const toolSnSet = new Set(points.map(p => p.toolSN));
        const lostTool = [];
        toolSnSet.forEach(t => {
          const tool = getDevice(t);
          if (tool) {
            this._tools.push(tool);
          } else {
            lostTool.push(t);
          }
        });

        if (lostTool.length > 0) {
          lostTool.forEach(t => {
            CommonLog.lError(`tool not found: ${t}`);
          });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL,`tool not found: ${lostTool.map(t=>`${t}`)}`));
        }

        const unhealthyTools = this._tools.filter(t => !t.Healthz);
        if (unhealthyTools.length > 0) {
          unhealthyTools.forEach(t => {
            CommonLog.lError(`tool not found: ${t.sn}`);
          });
          yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL,`tool not connected: ${unhealthyTools.map(t=>`${t.sn}`)}`));
        }

        yield call(this.updateData, (data: tScrewStepData): tScrewStepData => {
          return {
            points, // results data.results
            activeIndex: -1, // <-activeResultIndex
            ...data,
            jobID: payload.jobID,
            controllerMode: payload.controllerMode,
            retryTimes: 0
          };
        });


        yield put(orderActions.stepStatus(this, STEP_STATUS.DOING));
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep ENTERING' });
      }
    },


    * [STEP_STATUS.DOING](ORDER, orderActions) {
      try {
        yield call([this, handleResult], ORDER, orderActions, [], this._data);
        let isFirst = true;
        const sData: tScrewStepData = this._data;
        const {
          activeIndex,
          points
        } = sData;
        let activePoint = points[activeIndex];
        while (true) {
          const data = this._data;
          const nextAction = yield call([this, doPoint], activePoint, isFirst, orderActions);
          switch (nextAction.type) {
            case SCREW_STEP.RESULT:
              const { results: { data: results } } = nextAction;
              yield call([this, handleResult], ORDER, orderActions, results, data);
              const {activeIndex:nextIndex,points:nextPoints}=this._data;
              activePoint = nextPoints[nextIndex];
              break;
            case SCREW_STEP.REDO_POINT:
              const { point } = nextAction;
              activePoint = point;
              break;
            default:
              break;
          }
          if (isFirst) {
            isFirst = false;
          }
        }
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep DOING' });
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL,e))
      } finally {
        for (const t of this._tools) {
          yield call(t.Disable);
        }
      }
    },


    * [STEP_STATUS.FINISHED](ORDER, orderActions) {
      try {
        yield put(orderActions.doNextStep());
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep FINISHED' });
      }
    },


    * [STEP_STATUS.FAIL](ORDER, orderActions,msg) {
      try {
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
            title: '',
            content: (
              `拧紧工步失败${
              JSON.stringify(msg)||''}`
            )
          })
        );
        yield take(SCREW_STEP.CONFIRM_FAIL);
        yield put(orderActions.doNextStep());
      } catch (e) {
        CommonLog.lError(e, { at: 'screwStep FAIL' });
      }
    }
  };
}
