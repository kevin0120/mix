import OWebSocket from 'ws';
import userConfigs from '../shared/config';
// import { job } from './jobs';
// import {
//   fetchRoutingWorkcenterbyCarType,
//   fetchRoutingWorkcenterbyJobId,
//   invalidRoutingWorkcenter,
//   receiveRoutingWorkcenter
// } from './ongoingRoutingWorkcenter';

import { AIIS } from './actionTypes';
// import { canAcceptNewCar } from './index';

// import { addNewStory } from './storiesTimeLine';
// let ws = null;

// const WebSocket = require('@oznu/ws-connect');

// import {
//   GENERAL_NEW_TASK,
//   WEBSOCKET_CLOSED,
//   WEBSOCKET_STARTED,
//   WEBSOCKET_STARTINT,
//   GENERAL_NEW_JOB
// } from './index';
// import { WorkMode, AutoMode } from './commonActions'
// import { RECEIVE_NEW_CAR, RECEIVE_EMPTY_CAR } from './actionTypes';
// import { toolEnable } from './tools';
// import { setHealthzCheck } from './healthCheck';

export function initAiis(aiisUrl, hmiSN) {
  return {
    type: AIIS.INIT,
    aiisUrl,
    hmiSN
  };
}


// export function stopAiisWebsocket() {
//   if (
//     ws.ws.readyState === OWebSocket.OPEN ||
//     ws.ws.readyState === OWebSocket.CONNECTING
//   ) {
//     ws.close();
//   }
//   ws = null;
// }


// export function initAiis(dispatch, aiisUrl, hmiSN) {
//   if (ws) {
//     stopAiisWebsocket();
//   }
//
//   ws = new WebSocket(aiisUrl);
//   ws.on('open', () => {
//     // reg msg
//     ws.sendJson({ hmi_sn: hmiSN }, err => {
//       if (err) {
//         ws.close();
//       }
//     });
//
//     dispatch(setHealthzCheck('Andon', true));
//   });
//
//   ws.on('close', (code, reason) => {
//     dispatch(setHealthzCheck('Andon', false));
//     console.log(
//       `websocket disconnected. retry in 1s code: ${code}, reason: ${reason}`
//     );
//   });
//
//   ws.on('error', () => {
//     dispatch(setHealthzCheck('Andon', false));
//     console.log('websocket error. reconnect after 1s');
//   });
//   ws.on('ping', () => {
//     console.log(' receive ping Msg');
//   });
//   ws.on('pong', () => {
//     console.log(' receive pong Msg');
//   });
//
//   ws.on('message', dataRaw => {
//     const dataArray = dataRaw.split(';');
//
//     const data = dataArray.slice(-1);
//     const json = JSON.parse(data);
//
//     dispatch({ type: ANDON.NEW_DATA, json });
// /////////////////
//     // const currentState = getState();
//     //
//     // if (!currentState.isAutoMode || currentState.realMode === 'scanner'){
//     //   //不接受andon发来的信息
//     //   return
//     // }
//     //
//     // dispatch(addNewStory('info', 'Andon', json.vin_code));
//     //
//     // if (json.vin_code.length && canAcceptNewCar(currentState.orderProgress.orderStatus)) {
//     //
//     //
//     //   fetchRoutingWorkcenterbyCarType(currentState.connInfo.masterpc.connection, json.workcenter_code, json.cartype_code)
//     //     .then ( response => {
//     //       const statusCode = response.status;
//     //       if (statusCode === 200) {
//     //         dispatch(receiveRoutingWorkcenter(json.cartype_code, json.vin_code, response.data));
//     //         let mode = 'manual';
//     //         if (currentState.isAutoMode) {
//     //           mode = currentState.realMode;
//     //         }
//     //         job(currentState.connInfo.masterpc.connection,
//     //           currentState.connInfo.controllers[0],
//     //           json.cartype_code,
//     //           json.vin_code,
//     //           currentState.userInfo.id,
//     //           response.data.job,
//     //           response.data.points,
//     //           currentState.userConfigs.odooConnection.hmiSn.value,
//     //           response.data.product_id,
//     //           currentState.workingPage.infoUser.workcenter.id,
//     //           false,
//     //           !currentState.isAutoMode, mode)
//     //           .then ( resp => {
//     //             if (resp.status === 200) {
//     //               dispatch({
//     //                 type: GENERAL_JOB.SUCCESS,
//     //                 jobID: response.data.job
//     //               });
//     //
//     //             } else {
//     //             }
//     //           })
//     //           .catch(() => {
//     //           });
//     //
//     //       } else {
//     //         dispatch(invalidRoutingWorkcenter(null));
//     //       }
//     //     })
//     //     .catch(() => {
//     //       dispatch(invalidRoutingWorkcenter(null));
//     //     });
//     // }else {
//     //   // 空车信息
//     //
//     //   let mode = 'manual';
//     //   if (currentState.isAutoMode) {
//     //     mode = currentState.realMode;
//     //   }
//     //
//     //   if (canAcceptNewCar(currentState.orderProgress.orderStatus)) {
//     //     job(currentState.connInfo.masterpc.connection,
//     //       currentState.connInfo.controllers[0],
//     //       json.cartype_code,
//     //       json.vin_code,
//     //       currentState.userInfo.id,
//     //       currentState.userConfigs.byPassJob,
//     //       null,
//     //       currentState.userConfigs.odooConnection.hmiSn.value,
//     //       0,
//     //       currentState.workingPage.infoUser.workcenter.id,
//     //       true,
//     //       !currentState.isAutoMode, mode)
//     //       .then ( resp => {
//     //         if (resp.status === 200) {
//     //           dispatch({
//     //             type: GENERAL_JOB.SUCCESS,
//     //             jobID: currentState.userConfigs.byPassJob
//     //           });
//     //         } else {
//     //         }
//     //       })
//     //       .catch(() => {
//     //       });
//     //   }
//     //   // job();
//     //   // dispatch({
//     //   //   type: RECEIVE_EMPTY_CAR,
//     //   //   isEmptyCar: true
//     //   // });
//     // }
//   });
// }
