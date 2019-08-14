import { put } from 'redux-saga/effects';
import { CommonLog } from '../../common/utils';

const onConnectActions=[];

export function* bindOnConnectAction(actionToBind){
  try {
    onConnectActions.push(actionToBind);
    if(rushHealthz){
      yield put(actionToBind())
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

let rushHealthz=false;

export default function*(payload){
  try{
    rushHealthz=payload;
    if(payload){
      for(const action of onConnectActions){
        yield put(action());
      }
    }
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
  }catch(e){
    CommonLog.lError(e,{at:'rush handleHealthz'})
  }
}
