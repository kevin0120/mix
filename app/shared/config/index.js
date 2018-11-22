import { defaultConfigs } from './defaultConfig';

const settings = require('electron-settings');

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
