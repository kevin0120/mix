import { put } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';
import NotifierActions from '../Notifier/action';

const onConnectActions = [
  ()=>NotifierActions.enqueueSnackbar('Info', 'rush 已连接'),
];
const onDisconnectActions = [
  ()=>NotifierActions.enqueueSnackbar('Error', 'rush 已断开'),

];
const onChangeActions = [

];

export const bindRushAction = {
  onConnect: bindOnConnectAction,
  onDisconnect: bindOnDisconnectAction,
  onChange: bindOnChangeAction
};
let rushHealthz = false;


function* bindOnConnectAction(actionToBind) {
  try {
    onConnectActions.push(actionToBind);
    if (rushHealthz) {
      yield put(actionToBind());
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindOnDisconnectAction(actionToBind) {
  try {
    onConnectActions.push(actionToBind);
    if (!rushHealthz) {
      yield put(actionToBind());
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* bindOnChangeAction(actionToBind) {
  try {
    onConnectActions.push(actionToBind);
    yield put(actionToBind());
  } catch (e) {
    CommonLog.lError(e);
  }
}


export default function* (payload) {
  try {
    if (rushHealthz !== payload) {
      for (const action of onChangeActions) {
        yield put(action(payload));
      }
    }
    if (payload) {
      for (const action of onConnectActions) {
        yield put(action(payload));
      }
      CommonLog.lError('rush 已连接');
    } else {
      for (const action of onDisconnectActions) {
        yield put(action(payload));
      }
      CommonLog.lError('rush 已断开');
    }
    rushHealthz = payload;

    // const healthzStatus = state.healthCheckResults; // 获取整个healthz
    // if (!lodash.isEqual(healthzStatus.masterpc.isHealth, payload)) {
    //   // 如果不相等 更新
    //   yield put(setHealthzCheck('masterpc', payload));
    //   yield put(
    //     setNewNotification('info', `masterPC连接状态更新: ${payload}`)
    //   );
    // }
    // if (!payload) {
    //   yield put(setHealthzCheck('controller', false));
    //   yield put(
    //     setNewNotification('info', `controller连接状态更新: ${false}`)
    //   );
    // }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush handleHealthz' });
  }
}
