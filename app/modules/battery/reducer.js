import { BATTERY } from './action';
import { genReducers } from '../util';

const initBattery = {
  percentage: -1
};

const reducers = {
  [BATTERY.CHECK_OK]: setBatteryPercentage
};

function setBatteryPercentage(state, action) {
  return {
    percentage: action.percentage
  };
}

export default genReducers(reducers, initBattery);
