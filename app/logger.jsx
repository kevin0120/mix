const { Logger, transports } = require('winston');

const fs = require('fs');
const fse = require('fs-extra');
const path = require('path');
const os = require('os');

export const getAppDirectory = () => {
  switch (process.platform) {
    case 'darwin':
      return process.execPath.substring(
        0,
        process.execPath.indexOf('.app') + 4
      );
    case 'linux':
    case 'win32':
      return path.join(os.homedir(), '.controlPanel_sa');
    default:
  }
};

const homeDir = getAppDirectory();

const dir = path.join(homeDir, 'logs');

let gLogger = null;

const cpFormate = info => {
  return `${info.timestamp} ${info.level}: ${info.message}`;
};

export function CreateDailyLogger() {
  const isExist = fs.existsSync(dir);
  if (!isExist) {
    fse.mkdirpSync(dir);
  }
  if (gLogger === null) {
    gLogger = new Logger({
      level: 'info',
      transports: [
        new transports.File({
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
      exitOnError: false
    });
  }
}

export function Info(msg) {
  if (gLogger === null) return;
  gLogger.info(msg);
}

export function Warn(msg) {
  if (gLogger === null) return;
  gLogger.warn(msg);
}

export function Error(msg) {
  if (gLogger === null) return;
  gLogger.error(msg);
}

export function Query(options, f) {
  if (gLogger === null) {
    return;
  }
  return gLogger.query(options, f);
}
