import { POWER } from './action';
import { watchWorkers } from '../util';

const { ipcRenderer } = require('electron');

export default watchWorkers({
  [POWER.SHUTDOWN]: shutdown
});

function shutdown() {
  ipcRenderer.send('asynchronous-message', 'shutdown');
}
