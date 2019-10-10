// @flow
import type { QueryOptions } from 'winston';
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

const cpFormate = config => {
  const { level, message, meta, timestamp } = config;
  return `[${timestamp()}]  (${level})  (${message}) :  ${JSON.stringify(meta)}\n`;
};

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
          timestamp: true
        }),
        new winston.transports.File({
          name: 'alarm-file',
          filename: path.join(homeDir, 'logs/pretty.log'),
          tailable: true,
          json: false,
          maxFiles: 10,
          maxsize: 65535,
          timestamp: () => Date(),
          formatter: cpFormate,
          zippedArchive: true,
          prettyPrint: true
        })
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

export function Info(msg: string | Error, meta: ?Object) {
  if (gLogger === null) return;
  gLogger.info(msg, { meta });
}

export function Warn(msg: string | Error, meta: ?Object) {
  if (gLogger === null) return;
  gLogger.warn(msg, { meta });
}

export function Debug(msg: string | Error, meta: ?Object) {
  if (gLogger === null) return;
  gLogger.debug(msg, { meta });
}

export function Maintenance(msg: string | Error, meta: ?Object) {
  if (gLogger === null) return;
  gLogger.maintenance(msg, { meta });
}

export function lError(msg: string | Error, meta: ?Object) {
  if (gLogger === null) return;
  if (typeof msg === 'string') {
    gLogger.error(msg, { meta });
  }
  if (msg instanceof Error) {
    gLogger.error(msg.message, { meta });
  }
}

export function Query(options: QueryOptions, f?: (err: Error, results: any) => void) {
  if (gLogger === null) {
    throw new Error('logger does not exist');
  }
  return gLogger.query(options, f);
}
