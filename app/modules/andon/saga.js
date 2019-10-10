import { select, put, call, race, delay, take } from 'redux-saga/effects';
import { OPERATION_SOURCE, OPERATION_STATUS } from '../operation/model';
import { operationTrigger, OPERATION } from '../operation/action';
import { andonVehicleApi } from '../../api/andon';
import { jobManual } from '../../api/operation';

import { ANDON } from './action';
import { watchWorkers } from '../util';
import notifierActions from '../Notifier/action';
import { addNewStory } from '../timeline/saga';
import { STORY_TYPE } from '../timeline/model';

const workers = {
  [ANDON.NEW_DATA]: [call, handleAndonData],
  [ANDON.SCANNER]: [call, handleAndonScanner]
};

export default watchWorkers(workers);

function* handleAndonData(action) {
  try {
    const { data } = action;
    const state = yield select();
    if (state.operations.operationStatus !== OPERATION_STATUS.DOING) {
      if (data.cartype_code.length) {
        // 车辆拧紧作业
        yield put(
          operationTrigger(
            data.vin_code,
            data.cartype_code,
            null,
            OPERATION_SOURCE.ANDON
          )
        );
      } else {
        // 空车信息

        const { carType, carID } = state.operations;

        const { emptyCarJob } = state.setting.operationSettings;

        const controllerSN = state.connections.controllers[0].serial_no;
        const rushUrl = state.connections.masterpc;
        const { hmiSn } = state.setting.page.odooConnection;
        const userID = 1;
        const skip = true;
        const hasSet = false;
        const resp = yield call(
          jobManual,
          rushUrl,
          controllerSN,
          '',
          carType,
          carID,
          userID,
          emptyCarJob,
          hmiSn.value,
          0,
          skip,
          hasSet,
          ''
        );

        if (resp.statusCode !== 200) {
          yield put({ type: OPERATION.PROGRAMME.SET_FAIL });
        }
      }
    }
  } catch (e) {
    console.error('handleAndonData:', e);
  }
}

function* handleAndonScanner(action) {
  try {
    const state = yield select();
    const { operationStatus } = state.operations;
    if (operationStatus !== 'Ready') {
      return;
    }
    const { vin } = action;
    yield call(addNewStory, STORY_TYPE.INFO, 'scanner', vin);
    yield put({
      type: OPERATION.TRIGGER.NEW_DATA,
      carID: vin,
      carType: null
    });
    const { aiis, workcenterCode } = state.setting.system.connections;
    const resp = yield call(andonVehicleApi, aiis, vin, workcenterCode);
    if (resp) {
      const { overtime, successAction } = yield race({
        overtime: delay(8000),
        successAction: take(ANDON.NEW_DATA)
      });
      if (overtime) {
        yield put(notifierActions.enqueueSnackbar('Warn', '获取工单信息超时'));
      }
      if (successAction) {
        yield call(handleAndonData, successAction);
      }
    }
  } catch (e) {
    console.error(e);
    if (e.response) {
      console.error(e.response);
    }
  }
}
