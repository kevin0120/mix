import { take, put, fork, delay, cancel, call, takeEvery, select, race, takeLatest } from 'redux-saga/effects';
import React from 'react';
import Grid from '@material-ui/core/Grid';
import type { Saga } from 'redux-saga';
import { MANUAL, close, getresult, selectTool, selectPset, setData, inputOk1 } from './action';
import { CommonLog } from '../../common/utils';
import { tNS } from '../../i18n';
import { reworkDialogConstants as dia, reworkNS } from '../reworkPattern/constants';
import { getDevice, getDevicesByType } from '../deviceManager/devices';
import { deviceType } from '../deviceManager/constants';
import SelectCard from '../../components/SelectCard';
import { getPestListApi, psetApi } from '../../api/tools';
import { rushSendApi } from '../../api/rush';
import dialogActions from '../dialog/action';

import type { IDevice } from '../device/IDevice';
import notifierActions from '../Notifier/action';
import { addNewStory, STORY_TYPE } from './timeline';
import ResultInput from '../../components/ResultInput';
import type { tAction } from '../reworkPattern/interface/typeDef';
import reworkActions from '../reworkPattern/action';
import { DIALOG } from '../dialog/constants';
import { OPERATION_RESULT } from '../../containers/working/operations';


export default function* root() {
  try {
    yield takeEvery(MANUAL.TIGHTENING, startTightening);
    yield takeEvery(MANUAL.CLEAR, clearManual);
    // yield takeEvery(MANUAL.RESULTINPUT,recieveResult);
    yield takeLatest(MANUAL.START, manualWork);
    yield takeEvery(MANUAL.GET_RESULT, manualNewResult);
  } catch (e) {
    CommonLog.lError(e, {
      at: 'manual root'
    });
  }
}

const manualListeners = [];

// function* recieveResult(action) {
//   try{
//
//     if (action.resultIn?.sucess){
//       const r ={
//         tool_sn:"xx0011",
//         results:[{
//           tool_sn: "xx0011",
//           seq: 1,
//           group_seq: 1,
//           measure_time: 1,
//           measure_torque: action.resultIn?.result?.niu,
//           measure_angle: action.resultIn?.result?.jao,
//           measure_result: action.resultIn.result?.ok,
//           batch: '0',
//           count: 0}
//         ]
//       }
//
//       const tool = getDevice(r.tool_sn);
//       if (tool) {
//         yield call(tool.doDispatch, [r]);
//       } else {
//         CommonLog.lError('invalid tool', {
//           sn: r.tool_sn
//         });
//       }
//
//     }
//
//   }catch (e) {
//     console.error(e);
//   }
// }

function* clearManual() {
  try {
    while (manualListeners.length > 0) {
      const { tool, listener } = manualListeners.shift();
      tool.removeListener(listener);
      yield call(tool.Disable);
    }
  } catch (e) {
    notifierActions.enqueueSnackbar('Error', e.message, {
      at: 'clearManual'
    });
  }
}

// 手动模式收到结果
function* manualNewResult(action) {
  yield call(timeLineResult, action.result[0]);
}

// 手动输入拧紧结果
export function* manualResult(action: tAction = {}): Saga<void> {
  console.log('请手动输入拧紧结果');
  const { point } = action;
  try {
    if (point.isActive) {
      const buttons = [
        {
          label: '取消',
          color: 'info',
          action: reworkActions.cancelRework()
        },
        {
          label: '完成',
          color: 'success',
          action: inputOk1()
        }
      ];

      yield put(
        dialogActions.dialogShow({
          buttons,
          title: `请输入拧紧结果`,
          content: (<ResultInput/>)
        })
      );

      yield take(MANUAL.INPUT_OK);

      yield delay(300);
      const state = yield select();
      const { manual } = state;
      if (!manual.resultIn?.sucess) {
        yield put(
          notifierActions.enqueueSnackbar('Warn', '输入的拧紧结果不符合(扭矩值必填,ok值必填,或者数据类型)的约束', {
            at: '输入结果'
          })
        );
        return;
      }

      yield put(
        dialogActions.dialogShow({
          buttons,
          title: `请确认输入完成`
        })
      );

      yield take(MANUAL.INPUT_OK);

      if (manual.resultIn?.sucess) {

        // const r = {
        //     tool_sn: point.tightening_tool,
        //     seq: point.sequence,
        //     group_seq: point.group_sequence,
        //     measure_time: 0,
        //     measure_torque: manual.resultIn?.result?.niu,
        //     measure_angle: manual.resultIn?.result?.jao,
        //     measure_result: manual.resultIn?.result?.ok,
        //     batch: '1/24',
        //     count: point.max_redo_times,
        //     scanner:'',
        //   }

        const tool = getDevice(point.tightening_tool);
        if (tool) {
          // yield call(tool.doDispatch, [r]);
          const ControllerSN = ((tool.parent: any): IDevice)?.serialNumber;
          if (!ControllerSN) {
            throw new Error(`工具(${tool?.serialNumber})缺少控制器`);
          }
          console.log(ControllerSN);

          yield call(
            rushSendApi,
            'WS_TOOL_RESULT_SET',
            {
              tool_sn: point.tightening_tool,
              controller_sn: ControllerSN,
              measure_result: manual.resultIn?.result?.ok.toUpperCase(),
              measure_torque: parseFloat(manual.resultIn?.result?.niu),
              measure_angle: parseFloat(manual.resultIn?.result?.jao),
              count: point.max_redo_times + 1
            }
          );

        }
      }
    }

  } catch (e) {
    CommonLog.lError(e, {
      at: 'manualResult'
    });
  }
}

// let result;

function* timeLineResult(result) {
  const { measure_result } = result;
  const storyTypes = {
    [OPERATION_RESULT.OK]: STORY_TYPE.PASS,
    [OPERATION_RESULT.NOK]: STORY_TYPE.FAIL
  };
  const storyText = {
    [OPERATION_RESULT.OK]: '成功',
    [OPERATION_RESULT.NOK]: '失败'
  };
  yield call(
    addNewStory,
    storyTypes[measure_result] || STORY_TYPE.INFO,
    `结果 ${storyText[measure_result] || ''}`,
    `T=${result?.measure_torque}Nm A=${result?.measure_angle}° Tool=${result?.tool_sn} Scanner=${result?.scanner_code}`
  );
}

function isToolWorking(toolSN) {
  return manualListeners.find(l => {
    return l.toolSN === toolSN;
  });
}

function isToolWorkingWithPSET(toolSN, pset) {
  return manualListeners.find(l => {
    return l.toolSN === toolSN && l.pset === pset;
  });
}


function* startTightening() {
  try {
    const { tool: toolSN, scanner, controllerSN, pset } = yield select(s => s.manual);

    const tool = getDevice(toolSN);

    if (!tool) {
      throw new Error('没有找到工具');
    }

    if (scanner === '') {
      yield put(
        notifierActions.enqueueSnackbar('Warn', '没有条码进行拧紧结果追溯,不能进行手动作业,请输入条码', {
          at: 'controllerModes.pset'
        })
      );
      return;
    }

    if (isToolWorkingWithPSET(toolSN, pset)) {
      yield put(notifierActions.enqueueSnackbar('Warn', '当前工具和PSET正在作业中', {
        at: 'controllerModes.pset'
      }));
      return;
    }

    const retries = 1;
    for (let retry = 1; retry <= retries; retry += 1) {
      try {
        const set = yield call(
          psetApi,
          toolSN || '',
          controllerSN || '',
          '手动工步',
          '',
          pset,
          0,
          0,
          '手动工单',
          scanner,
          1
        );
        if (set.data.result === 0) {
          const msg1 = `pset设置成功, 工具：${toolSN}`;
          yield put(notifierActions.enqueueSnackbar('Info', msg1, {
            at: 'manual oK'
          }));
          break;
        }
      } catch (e) {
        const msg = `pset失败，${e.message}, 工具：${toolSN}`;
        if (retry === retries) {
          throw new Error(msg);
        }
      }
    }

    yield call(tool?.Enable || (() => {
      CommonLog.lError(
        `tool ${tool?.Name}: no such tool or tool cannot be enabled.`
      );
    }));


    const toolWorking = isToolWorking(toolSN, pset);
    if (toolWorking) {
      toolWorking.pset = pset;
      return;
    }

    const listener = tool.addListener(
      () => true,
      input => getresult(input.data)
    );
    manualListeners.push({
      tool,
      toolSN,
      pset,
      listener
    });

  } catch (e) {
    yield put(
      notifierActions.enqueueSnackbar('Error', e.message, {
        at: 'manual startTightening'
      })
    );
  }
}

function* manualWork() {
  try {
    // const SelectCard1 = SelectCard(noScanner);

    // yield put(dialogActions.dialogShow({
    //   buttons: [btnCancel],
    //   title: '请扫码用以追溯',
    //   content: (<Grid spacing={1} container>
    //     <Grid item xs={6} key={`1`}>
    //       <SelectCard1
    //         name={'不扫码'}
    //         height={130}
    //         item={this}
    //       />
    //     </Grid>
    //    </Grid>)
    // }));
    //
    // const {data} =yield take([manual.NOSCANNER,manual.NEWCANNERDATA])


    const { toolAction, canceled } = yield call(showToolSelectDialog);
    if (!toolAction || canceled) {
      return;
    }
    const { tool } = toolAction;
    const ControllerSN = ((tool.parent: any): IDevice)?.serialNumber;

    yield delay(300); // 动画时间
    const { psetAction, canceled: psetCanceled } = yield call(showPsetSelectDialog, tool, ControllerSN);
    if (!psetAction || psetCanceled) {
      return;
    }
    const { pset } = psetAction;
    yield put(setData(ControllerSN, tool.serialNumber, pset));
  } catch (e) {
    CommonLog.lError(e, {
      at: 'manualWork'
    });
  }
}

function* showToolSelectDialog() {
  const btnCancel = {
    label: tNS(dia.manualCancel, reworkNS),
    color: 'warning',
    action: close()
  };
  const tools = getDevicesByType(deviceType.tool);
  const ToolSelectCard = SelectCard(selectTool);

  yield put(dialogActions.dialogShow({
    buttons: [btnCancel],
    title: '选择工具1',
    content: (<Grid spacing={1} container>
      {tools.map(t => <Grid item xs={6} key={`${t.serialNumber}`}>
        <ToolSelectCard
          name={t.Name}
          status={t.Healthz ? '已连接' : '已断开'}
          infoArr={[t.serialNumber]}
          height={130}
          item={t}
        />
      </Grid>)}
    </Grid>)
  }));
  const resp = yield race({
    toolAction: take(MANUAL.SELECT_TOOL),
    canceled: take([MANUAL.CLOSE, DIALOG.CLOSE])
  });
  yield put(dialogActions.dialogClose());
  return resp;
}

function* showPsetSelectDialog(tool, ControllerSN) {
  const btnCancel = {
    label: tNS(dia.manualCancel, reworkNS),
    color: 'warning',
    action: close()
  };
  const PsetSelect = SelectCard(selectPset);
  const psets = (yield call(getPestListApi, tool.serialNumber, ControllerSN))?.data || [];
  yield put(dialogActions.dialogShow({
    buttons: [btnCancel],
    title: '选择PSET',
    content: (<Grid spacing={1} container>
      {psets.map(p => <Grid item xs={6} key={`${p}`}>
        <PsetSelect
          name={p}
          height={130}
          item={p}
        />
      </Grid>)}
    </Grid>)
  }));
  const resp = yield race({
    psetAction: take(MANUAL.SELECT_PSET),
    canceled: take([MANUAL.CLOSE, DIALOG.CLOSE])
  });
  yield put(dialogActions.dialogClose());
  return resp;
}
