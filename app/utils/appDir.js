const fs = require('fs');
const fse = require('fs-extra');
const path = require('path');
const os = require('os');

export const getAppDirectory = () => {
  let dir = path.join(os.homedir(), '.controlPanel_sa');
  if (
    process.env.NODE_ENV === 'development' ||
    process.env.DEBUG_PROD === 'true' ||
    process.env.NODE_ENV === 'test'
  ) {
    dir = path.join(os.homedir(), '.controlPanel_sa');
  }
  if (process.platform === 'darwin') {
    dir = process.execPath.substring(0, process.execPath.indexOf('.app') + 4);
  }
  const isExist = fs.existsSync(dir);
  if (!isExist) {
    fse.mkdirpSync(dir);
  }
  return dir;
};
