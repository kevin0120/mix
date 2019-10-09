import { call, put } from 'redux-saga/effects';
import { BATTERY, batteryCheckOK } from './action';
import { watchWorkers } from '../util';

const util = require('util');
const exec = util.promisify(require('child_process').exec);

const workers = {
  [BATTERY.CHECK]: doCheckBattery
};

export default watchWorkers(workers);

function* doCheckBattery() {
  try {
    // upower -i $(upower -e | grep BAT) | grep --color=never -E "percentage"
    const { error, stdout, stderr } = yield call(
      execute,
      'upower -i $(upower -e | grep BAT) | grep --color=never -E "percentage"'
    );
    if (error) {
      console.log(error);
      return;
    }
    if (stderr) {
      console.log(stderr);
      return;
    }
    if (stdout) {
      const result = /\d+/.exec(stdout);
      let percentage = -1;
      if (result) {
        percentage = parseInt(result['0'], 10);
      }
      yield put(batteryCheckOK(percentage));
    }
  } catch (e) {
    console.error(e);
  }
}

function execute(command) {
  return exec(command)
    .then(resp => resp)
    .catch(e => ({ error: e }));
}
