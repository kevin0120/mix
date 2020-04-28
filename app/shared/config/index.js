import defaultConfigs from './defaultConfig';
import { getAppDirectory } from '../../utils/appDir';

const path = require('path');
const settings = require('electron-settings');

const dir = getAppDirectory();

settings.setPath(path.join(dir, 'setting.json'));

if (
  process.env.NODE_ENV === 'development'
  // || process.env.DEBUG_PROD === 'true'
  // || process.env.NODE_ENV === 'test'
) {
  settings.deleteAll();
}

if (Object.keys(settings.getAll()).length === 0) {
  settings.setAll(defaultConfigs,{
    prettify: true
  });
}

const configs = settings.getAll();
export default configs;
