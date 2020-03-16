import { take, put, fork,delay,cancel,call,takeEvery,select,takeLatest} from 'redux-saga/effects';
import React from 'react';
import Grid from '@material-ui/core/Grid';
import { MANUAL, start, close, getresult, selectTool, selectPset, setData } from './action';
import { CommonLog } from '../../common/utils';
import { tNS } from '../../i18n';
import { reworkDialogConstants as dia, reworkNS } from '../reworkPattern/constants';
import { getDevice, getDevicesByType } from '../deviceManager/devices';
import { deviceType } from '../deviceManager/constants';
import SelectCard from '../../components/SelectCard';
import { getPestListApi, psetApi } from '../../api/tools';
import dialogActions from '../dialog/action';

import type { IDevice } from '../device/IDevice';
import notifierActions from '../Notifier/action';
import { addNewStory, STORY_TYPE } from './timeline';
import ResultInput from '../../components/ResultInput';
import type { tDeviceSN } from '../device/typeDef';
import type { tResult, tResultStatus } from '../step/screwStep/interface/typeDef';


export default function* root() {
  try {
    yield takeEvery(MANUAL.CANCEL,initManual);
    yield takeLatest(MANUAL.TIGHTENING,oK);
    yield takeEvery(MANUAL.CLICKPOINT,manualResult);
    yield takeEvery(MANUAL.RESULTINPUT,recieveResult);

    while (true) {
      yield take(MANUAL.START);
      const work =yield fork(manualWork);
      yield take(MANUAL.CLOSE);
      if (work) {
        yield cancel(work);
      }
    }
  } catch (e) {
    CommonLog.lError(e);
  }
}

function* recieveResult(action) {
  try{

    if (action.resultIn?.sucess){
      const r ={
        tool_sn:"xx0011",
        results:[{
          tool_sn: "xx0011",
          seq: 1,
          group_seq: 1,
          measure_time: 1,
          measure_torque: action.resultIn?.result?.niu,
          measure_angle: action.resultIn?.result?.jao,
          measure_result: action.resultIn.result?.ok,
          batch: '0',
          count: 0}
        ]
      }

      const tool = getDevice(r.tool_sn);
      if (tool) {
        yield call(tool.doDispatch, [r]);
      } else {
        CommonLog.lError('invalid tool', {
          sn: r.tool_sn
        });
      }


    }


  }catch (e) {
    console.error(e);
  }
}

function* manualResult() {
  try{
    const buttons = [
      {
        label: "取消",
        color: 'info',
        action: dialogActions.dialogClose()
      },
      {
        label: "完成",
        color: 'success',
        action: dialogActions.dialogClose()
      }
    ];

    yield put(
      dialogActions.dialogShow({
        buttons,
        title: `请输入拧紧结果`,
        content: (<ResultInput/>)
      })
    );
  }catch (e) {
    console.error(e);
  }
}

let result;

function* oK() {
  try {
    const state = yield select();

    const  {manual} = state;
    const tool = getDevice(manual?.tool);

    if (result !==null&& typeof result !== 'undefined') {
      tool.removeListener(result);
    }


    const retries = 1;
    for (let retry = 1; retry <= retries; retry += 1) {
      if (manual.scanner===''){
        yield put(
          notifierActions.enqueueSnackbar('Warn', '没有条码进行拧紧结果追溯,不能进行手动作业请输入条码', {
            at: 'controllerModes.pset'
          })
        );
        break;
      }
      try {
        const set = yield call(
          psetApi,
          manual?.tool || '',
          manual?.controllerSN || '',
          '手动工步',
          '',
          manual?.pset,
          0,
          // retryTimes,
          0,
          '手动工单',
          manual?.scanner
        );
        if (set.data.result === 0) {
          const msg1 = `pset设置成功, 工具：${manual?.tool}`;
          yield put(
            notifierActions.enqueueSnackbar('Info', msg1, {
              at: 'controllerModes.pset'
            })
          );
        }

        break;
      } catch (e) {
        const msg = `pset失败，${e.message}, 工具：${manual?.tool}`;
        yield put(
          notifierActions.enqueueSnackbar('Error', msg, {
            at: 'controllerModes.pset'
          })
        );
        if (retry === retries) {
          throw e;
        }
      }
    }

    if (manual.scanner!=='') {
      yield call(tool?.Enable || (() => {
        CommonLog.lError(
          `tool ${tool?.Name}: no such tool or tool cannot be enabled.`
        );
      }));


      result = tool.addListener(
        () => true,
        input => getresult(input.data)

      )

      while (true){
      const action= yield take(MANUAL.GETRESULT);

      if (action.result[0]?.measure_result==="NOK"){
        yield call(
          addNewStory,
          STORY_TYPE.FAIL,
          `结果 失败`,
          `T=${action.result[0]?.measure_torque}Nm A=${action.result[0]?.measure_angle}° Tool=${action.result[0]?.tool_sn} Scanner=${action.result[0]?.scanner_code}`
        );

      } else {
        yield call(
          addNewStory,
          STORY_TYPE.PASS,
          `结果 成功`,
          `T=${action.result[0]?.measure_torque}Nm A=${action.result[0]?.measure_angle}° Tool=${action.result[0]?.tool_sn} Scanner=${action.result[0]?.scanner_code}`
        );
      }
      }
    }

  }catch (e) {
    console.error(e);
  }
}

function* initManual() {
  try{
  yield put(close());
  yield delay(300);
  yield put(start());
  }catch (e) {
    console.error(e);
  }
}

function* manualWork() {
  try {
    // const SelectCard1 = SelectCard(noScanner);
    const btnCancel = {
      label: tNS(dia.manualCancel, reworkNS),
      color: 'warning',
      action: close()
    };
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

    const tools = getDevicesByType(deviceType.tool);
    const ToolSelectCard = SelectCard(selectTool);
    yield delay(300);

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
    const { tool } = yield take(MANUAL.SELECT_TOOL,);
    // this._tools = [tool];
    // this._forceTool = tool;
    yield put(dialogActions.dialogClose());
    const PsetSelect = SelectCard(selectPset);
    const ControllerSN = ((tool.parent: any): IDevice)?.serialNumber;
    // TODO get tool psets
    const psets= (yield call(getPestListApi, tool.serialNumber, ControllerSN))?.data || [];
    yield delay(300);
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
    const { pset } = yield take(MANUAL.SELECT_PSET);
    // this._forcePset = pset;
    yield put(dialogActions.dialogClose());

    yield put(setData(ControllerSN,tool.serialNumber,pset));

  }
   catch (e) {
    CommonLog.lError(e);
  }
}
