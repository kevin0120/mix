const settings = window.require('electron-settings');

import { defaultConfigs } from './defaultConfig';

if (
  process.env.NODE_ENV === 'development' ||
  process.env.DEBUG_PROD === 'true'
) {
  settings.deleteAll();
}

if (Object.keys(settings.getAll()).length === 0) {
  settings.setAll(defaultConfigs);
  // console.log('set default configs');
}

const configs = settings.getAll();
export default configs;
