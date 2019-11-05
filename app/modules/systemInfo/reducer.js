import { genReducers } from '../util';
import { SYSTEM_INFO } from './action';

const eSetting = require('electron-settings');


const initSystemState = {
  workcenter: eSetting.get('system.workcenter.code')
};

const systemReducers = {
  [SYSTEM_INFO.SET_WORKCENTER]: (state, { workcenter }) => ({
    ...state,
    workcenter
  })
};

export default genReducers(systemReducers, initSystemState);
