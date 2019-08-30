// @flow
import { primaryColor, dangerColor, infoColor, warningColor, roseColor } from './common/theme/palette';

const winston = require('winston');
const fs = require('fs');
const fse = require('fs-extra');
const path = require('path');
const os = require('os');

export const getAppDirectory = () => {
  switch (process.platform) {
    case 'darwin':
      return process.execPath.substring(0, process.execPath.indexOf('.app') + 4);
    case 'linux':
    case 'win32':
      return path.join(os.homedir(), '.controlPanel_sa');
    default:
      return path.join(os.homedir(), '.controlPanel_sa');
  }
};

const homeDir = getAppDirectory();

const dir = path.join(homeDir, 'logs');

let gLogger = null;

const cpFormate = info => `${info.timestamp} ${info.level}: ${info.message}`;

export function CreateDailyLogger() {
  const isExist = fs.existsSync(dir);
  if (!isExist) {
    fse.mkdirpSync(dir);
  }
  if (gLogger === null) {
    gLogger = new winston.Logger({
      level: 'info',
      transports: [
        new winston.transports.File({
          name: 'info',
          filename: path.join(homeDir, 'logs/event.log'),
          level: 'info',
          tailable: true,
          maxFiles: 10,
          maxsize: 65535,
          zippedArchive: true,
          formatter: cpFormate
        })
        // new transports.File({
        //   name: "alarm-file",
        //   filename: './logs/alarm.log',
        //   level: 'error',
        //   tailable: true,
        //   maxFiles: 10,
        //   maxsize: 65535,
        //   formatter: cpFormate,
        //   zippedArchive: true
        // }),
      ],
      exitOnError: false,
      levels: { error: 0, warn: 1, maintenance: 2, info: 3, debug: 4 }
    });
    winston.addColors({
      error: dangerColor,
      warn: warningColor,
      info: primaryColor,
      debug: infoColor,
      maintenance: infoColor,
      query: roseColor
    });
  }
}

export function Info(msg, meta) {
  if (gLogger === null) return;
  gLogger.info(msg, { meta });
}

export function Warn(msg, meta) {
  if (gLogger === null) return;
  gLogger.warn(msg, { meta });
}

export function Debug(msg, meta) {
  if (gLogger === null) return;
  gLogger.debug(msg, { meta });
}

export function Maintenance(msg, meta) {
  if (gLogger === null) return;
  gLogger.maintenance(msg, { meta });
}

export function lError(msg: mixed, meta) {
  if (gLogger === null) return;
  if (typeof msg === 'string') {
    gLogger.error(msg, { meta });
  }
  if (msg instanceof Error) {
    gLogger.error(msg.message, { meta });
  }
}

export function Query(options, f) {
  if (gLogger === null) {
    throw new Error('logger does not exist');
  }
  return gLogger.query(options, f);
}
