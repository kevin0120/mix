import { CommonLog } from '../../common/utils';

export default function*(payload){
  try{
    console.log('healthz',payload);
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
