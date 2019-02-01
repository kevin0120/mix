/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

import { OPERATION } from './actionTypes';
import { OPERATION_SOURCE } from '../reducers/operations';

export function switch2Ready() {
  return {
    type: OPERATION.FINISHED
  };
}

export function switch2Doing() {
  return {
    type: OPERATION.STARTED
  };
}

// export function operationVerified(data) {
//   return {
//     type: OPERATION.VERIFIED,
//     data
//   };
// }

export function switch2Timeout() {
  return {
    type: OPERATION.TIMEOUT
  };
}

export function switch2PreDoing() {
  return {
    type: OPERATION.PREDOING
  };
}

export function operationTrigger(carID, carType, job, source) {
  return {
    type: OPERATION.TRIGGER.TRIGGER,
    carID,
    carType,
    job,
    source
  };
}

export function operationBypassIO(){
  return{
    type:OPERATION.BYPASS.IO
  }
}

export function operationBypassConfirm(){
  return{
    type:OPERATION.BYPASS.CONFIRM
  }
}

export function operationBypassCancel(){
  return{
    type:OPERATION.BYPASS.CANCEL
  }
}

export function operationConflictDetected(data){
  return{
    type:OPERATION.CONFLICT.DETECTED,
    data
  }
}

export function operationConflictConfirm(data){
  return{
    type:OPERATION.CONFLICT.CONFIRM,
    data
  }
}

export function operationConflictCancel(){
  return{
    type:OPERATION.CONFLICT.CANCEL
  }
}

export function operationTriggerBlock(block){
  return{
    type:OPERATION.TRIGGER.BLOCK,
    block
  }
}
