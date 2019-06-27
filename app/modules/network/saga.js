/* eslint-disable no-bitwise,no-plusplus */
import { call, put, select } from 'redux-saga/effects';
import { cloneDeep } from 'lodash';
import { NETWORK } from './action';
import saveConfigs from '../setting/action';
import { watch } from '../indexSaga';
import { setNewNotification } from '../notification/action';


const lodash = require('lodash');
const util = require('util');
const exec = util.promisify(require('child_process').exec);

const section = 'network';


const netmask2CIDR = netmask =>
  netmask
    .split('.')
    .map(Number)
    .map(part => (part >>> 0).toString(2))
    .join('')
    .split('1').length - 1;

const CDIR2netmask = (bitCount) => {
  const mask = [];
  let count = bitCount;
  for (let i = 0; i < 4; i++) {
    const n = Math.min(count, 8);
    mask.push(256 - 2 ** (8 - n));
    count -= n;
  }
  return mask.join('.');
};

const workers = {
  [NETWORK.CHECK]: doCheckCurrentConnection,
  [NETWORK.SCAN]: doScan,
  [NETWORK.SET]: doSet,
  [NETWORK.SIGNAL]: doCheckActiveSignal,
  [NETWORK.SET_FAIL]: networkFail
};

export default watch(workers);

function* doCheckCurrentConnection() {
  try {
    const { error, stdout, stderr } = yield call(execNmcli, 'nmcli -t -s con show default');
    if (error) {
      console.log(`exec error: ${error}`);
      return;
    }
    if (stderr) {
      console.log(`stderr: ${stderr}`);
      return;
    }
    if (stdout) {
      const result = stdout.toString();
      const ssid = /802-11-wireless.ssid:(.+)\n/.exec(result)[1] || '';
      const password = /802-11-wireless-security.psk:(.+)\n/.exec(result)[1] || '';
      const addresses = /ipv4.addresses:(.+)\/(.+)\n/.exec(result) || '';
      const address = addresses[1] || '';
      const gateway = /ipv4.gateway:(.+)\n/.exec(result)[1] || '';
      const CDIR = parseInt(addresses[2], 10) || 0;
      const mask = CDIR2netmask(CDIR);
      const data = yield select(state => state.network.config);
      const tempData = cloneDeep(data);

      tempData.gateway.value = gateway;
      tempData.ipAddress.value = address;
      tempData.netmask.value = `${mask}`;
      tempData.password.value = password;
      tempData.ssid.value = ssid;

      yield put({
        type: NETWORK.CHECK_OK,
        data: tempData
      });
    }
  } catch (e) {
    console.error(e);
  }
}

function* doCheckActiveSignal() {
  try {
    const { error, stdout, stderr } = yield call(execNmcli, 'LANG=eng nmcli -f active,ssid,signal -t -e no dev wifi');
    if (error) {
      console.log(`exec error: ${error}`);
      return;
    }
    if (stderr) {
      console.log(`stderr: ${stderr}`);
      return;
    }
    if (stdout) {
      const result = /(yes):(.*):(\d*)/.exec(stdout);
      if (result) {
        const ssid = result[2];
        const signal = result[3];
        yield put({
          type: NETWORK.SIGNAL_OK,
          ssid,
          signal
        });
        return;
      }
    }
    yield put({
      type: NETWORK.SIGNAL_OK,
      ssid: '',
      signal: 0
    });
  } catch (e) {
    console.error(e);
  }
}

function* doSet(action) {
  try {
    const { data } = action;
    let ret = 0;
    const mask = netmask2CIDR(data.netmask.value);
    yield call(execNmcli, 'nmcli con delete default');
    const { error, stdout, stderr } = yield call(execNmcli,
      `LANG=eng nmcli dev wifi connect '${data.ssid.value}' password '${
        data.password.value
        }' name default`);
    if (error) {
      console.log(`exec error: ${error}`);
      ret = -1;
    }
    if (stderr) {
      console.log(`stderr: ${stderr}`);
      ret = -1;
    }
    if (stdout) {
      console.log(`stdout: ${stdout}`);
      if (!/successfully activated/.test(stdout)) {
        console.log(`connect failed`);
        ret = -1;
      }
    }
    if (ret < 0) {
      yield call(execNmcli, 'nmcli con delete default');
      yield put({
        type: NETWORK.SET_FAIL,
        message: '网络设置失败'
      });
      return;
    }
    yield call(execNmcli, 'nmcli con down default');
    const cmd = `nmcli con mod default ipv4.method manual ipv4.address ${
      data.ipAddress.value
      }/${mask} ipv4.gateway ${data.gateway.value}`;
    const { error: err, stdout: sOut, stderr: sErr } = yield call(execNmcli, cmd);
    if (err) {
      console.log(`exec error: ${err}`);
      ret = -1;
    }
    if (sErr) {
      console.log(`stderr: ${sErr}`);
      ret = -1;
    }

    if (ret === 0) {
      yield put({
        type: NETWORK.SET_OK,
        data
      });
      yield call(execNmcli, 'nmcli con up default');
      yield put(saveConfigs(section, data));
    } else {
      yield call(execNmcli, 'nmcli con delete default');
      yield put({
        type: NETWORK.SET_FAIL,
        message: '网络设置失败'
      });
    }
  } catch (e) {
    console.error(e);
  }
}

function* doScan() {
  try {
    const ssidList = [];
    const { error, stdout, stderr } = yield call(execNmcli, 'nmcli -f ssid -t -e no dev wifi');
    if (error) {
      console.log(`exec error: ${error}`);
      return;
    }
    if (stdout) {
      const lines = stdout.toString().split('\n');
      for (let i = 0; i < lines.length - 1; i += 1) {
        const line = lines[i].toString();
        ssidList.push(line);
      }
      yield put({
        type: NETWORK.SCAN_OK,
        data: lodash.uniq(ssidList)
      });
    }
    if (stderr) {
      console.log(stderr);
    }
  } catch (e) {
    console.error(e);
  }
}

function* networkFail(action) {
  try {
    yield put(setNewNotification('error', action.message || '无线网络错误'));
  } catch (e) {
    console.error(e);
  }
}

function execNmcli(command) {
  return exec(command)
    .then(resp => resp)
    .catch(e => ({ error: e }));
}



