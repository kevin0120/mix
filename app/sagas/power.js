import {POWER} from '../actions/actionTypes';
import {watch} from './utils';

const { ipcRenderer } = require('electron');

export default watch({
  [POWER.SHUTDOWN]:shutdown
});

function shutdown(){
  ipcRenderer.send('asynchronous-message', 'shutdown');
}
