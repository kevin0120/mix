import defaultConfigs from './defaultConfig';
import { getAppDirectory } from '../../logger';

const path = require('path');
const settings = require('electron-settings');
const fse = require('fs-extra');
const fs = require('fs');

const dir = getAppDirectory();
const isExist = fs.existsSync(dir);
if (!isExist) {
  fse.mkdirpSync(dir);
}

settings.setPath(path.join(dir, 'setting.json'));

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
