import { POWER } from './action';
import { watch } from '../indexSaga';

const { ipcRenderer } = require('electron');

export default watch({
  [POWER.SHUTDOWN]: shutdown
});

function shutdown() {
  ipcRenderer.send('asynchronous-message', 'shutdown');
}
