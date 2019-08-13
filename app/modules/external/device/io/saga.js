// @flow
// redux-saga
import {
  call
} from 'redux-saga/effects';
import type { Saga } from 'redux-saga';
// reducers
import ClsIOModule from './model';
import { DefaultInput, DefaultOutput } from './type';
import { CommonLog } from '../../../../common/utils';

// config

// const io = {
//   channel: null,
//   client: null,
//   recon: null,
//   currentKeyStatus: null,
//   modbusClient: null,
//   senderReceiver: null,
//   runningTask: null,
//   health: false,
//   i: {
//     resetKey: 0,
//     byPass: 1,
//     modeSelect: 2
//   },
//   o: {
//     white: 0,
//     yellow: 1,
//     green: 2,
//     red: 3,
//     beep: 4
//   }
// };

export const defaultIO = new ClsIOModule('IO Module', '1', { input: DefaultInput, output: DefaultOutput });

type tIORushData = {
  type: string,
  data: Object
};

export default function* ioNewData(data: tIORushData): Saga<void> {
  try {
    const msgType = data.type;
    switch (msgType) {
      case 'WS_IO_CONTACT':
        yield call(defaultIO.doDispatch, {
          sn: data.data.sn,
          direction: data.data.type,
          contact: data.data.contact
        });
        break;
      case 'WS_IO_STATUS':
        console.log(data);
        break;
      default:
        CommonLog.lError('IO Message Type Is Not Defined', { msgType });
    }
  } catch (e) {
    CommonLog.lError(e, { at: 'rush event io' });
  }
}

//
// function* ioInputHandler(action: tCommonActionType & tDeviceNewData): Saga<void> {
//   try {
//     const { data } = action;
//     if (defaultIO.doValidate(data)) {
//       const respActions = defaultIO.doDispatch(data);
//       if (respActions instanceof Array) {
//         // eslint-disable-next-line
//         for (const a of respActions) {
//           yield put(a);
//         }
//       }
//     } else {
//       // do nothing
//     }
//   } catch (e) {
//     CommonLog.lError(e);
//   }
// }

// export default function* watchIOEvent(): Saga<void> {
//   try {
//     while (true) {
//       const action: AnyAction = yield take([IO_ACTION.RESET, IO_ACTION.DATA_ONCHANGE]);
//       switch (action.type) {
//         case IO_ACTION.DATA_ONCHANGE: {
//           yield fork(ioInputHandler, action);
//           break;
//         }
//         // case IO_ACTION.RESET: {
//         //   yield fork(resetIO, action.modbusConfig);
//         //   break;
//         // }
//         default:
//           break;
//       }
//     }
//   } catch (e) {
//     CommonLog.lError(e);
//   }
// }

// export function* handleIOFunction(data: ?tIOContact): Saga<void> {
//   try {
//     const state = yield select();
//     if (state.router.location.pathname !== '/working') {
//       return;
//     }
//     switch (data) {
//       case IO_FUNCTION.IN.RESET: {
//         // 复位
//         if (state.operations.operationStatus === OPERATION_STATUS.FAIL) {
//           // 只有在fail阶段可以执行复位功能
//           yield call(continueOperation);
//         }
//         break;
//       }
//       case IO_FUNCTION.IN.BYPASS: {
//         // 强制放行
//         // yield put({ type: OPERATION.FINISHED, data: [] });
//         if (
//           [OPERATION_STATUS.DOING, OPERATION_STATUS.FAIL].includes(
//             state.operations.operationStatus
//           )
//         ) {
//           // yield put(openShutdown('bypass'));
//           if (state.setting.operationSettings.byPass) {
//             yield put(operationBypassIO());
//
//           }
//         }
//         break;
//       }
//       case IO_FUNCTION.IN.MODE_SELECT: {
//         // 模式选择
//
//         break;
//       }
//
//       default:
//         break;
//     }
//   } catch (e) {
//     console.error(e);
//   }
// }


// export function setModBusIO(modbusConfig) {
//   const modbusOutConfig = modbusConfig.out;
//   const modbusInConfig = modbusConfig.in;
//   lodash.forEach(modbusOutConfig, setOutBit);
//   lodash.forEach(modbusInConfig, setInBit);
// }
//
// function resetIO(modbusConfig) {
//   const preIOStatus = lodash.cloneDeep(ioStatus);
//   const preO = lodash.cloneDeep(io.o);
//   const { o } = io;
//
//   setModBusIO(modbusConfig);
//   const keys = Object.keys(o);
//   for (const key of keys) {
//     ioStatus[o[key]] = sOff; // 首先复位为关闭
//     ioStatus[o[key]] = preIOStatus[preO[key]];
//   }
// }
//
// export function setLedStatusReady() {
//   const { o } = io;
//   ioStatus[o.red] = sOff;
//   ioStatus[o.white] = sOff;
//   ioStatus[o.yellow] = sOff;
//   ioStatus[o.green] = sBlinkOn;
// }
//
// export function setLedStatusDoing() {
//   const { o } = io;
//   ioStatus[o.red] = sOff;
//   ioStatus[o.yellow] = sOff;
//   ioStatus[o.white] = sOn;
//   ioStatus[o.green] = sOn;
// }
//
// export function setLedInfo(flag) {
//   ioStatus[io.o.white] = flag;
// }
//
// export function setLedWarning(flag) {
//   ioStatus[io.o.green] = flag;
// }
//
// export function setLedError(flag) {
//   ioStatus[io.o.red] = flag;
// }
//
// export function testIO(ioConfig, idx) {
//   switch (ioConfig) {
//     case 'out': {
//       io.modbusClient
//         .readCoils(idx, 1)
//         .then(resp => {
//           const got = resp.response.body.valuesAsArray[0];
//           io.modbusClient.writeSingleCoil(idx, got === 0);
//           return true;
//         })
//         .catch(() => 'fail');
//       break;
//     }
//     case 'in': {
//       return io.modbusClient.readDiscreteInputs(idx, 1);
//     }
//     default:
//       break;
//   }
// }
//
// export function getIBypass() {
//   return io.i.byPass;
// }
//
// export function getIModeSelect() {
//   return io.i.modeSelect;
// }
