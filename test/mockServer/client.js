const WebSocket = require('@oznu/ws-connect');

// establish new websocket connection
const ws = new WebSocket('ws://192.168.4.221:8082/rush/v1/ws');

const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

ws.on('open', () => {
  console.log('open');
  ws.sendJson({
    type: 'WS_REG',
    sn: 0,
    data: {
      hmi_sn: '1560c527ac7f4f9f90e4900b50717457'
    }
  });
  console.log('ws open');
  ws.on('message', message => {
    console.log(message);
  });
});

rl.on('line', () => {
  const data = JSON.parse(
    '{"code":"MO0001","track_code":"XM205Z03","workcenter":"TA-01L","product_code":"M000001780589","date_start":"2019-11-07T11:31:17+08:00"}'
  );

  const json = {
    sn: 1,
    type: 'WS_ORDER_START_REQUEST',
    data: data
  };
  ws.sendJson({
    sn: 1,
    type: 'WS_ORDER_START_REQUEST',
    data: data
  });
  console.log('sending...');
  console.log(json);
});
