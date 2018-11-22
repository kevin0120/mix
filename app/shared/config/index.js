import { defaultConfigs } from './defaultConfig';
import {getAppDirectory} from '../../logger';

const path = require('path');
const settings = require('electron-settings');

settings.setPath(path.join(getAppDirectory(), 'setting.json'));

if (
  process.env.NODE_ENV === 'development' ||
  process.env.DEBUG_PROD === 'true'
) {
  settings.deleteAll();
}

if (Object.keys(settings.getAll()).length === 0) {
  settings.setAll(defaultConfigs);
}

const configs = settings.getAll();
export default configs;
