const WebSocket = require('ws');
const handlers = require('./eventsHandler');
const { events } = require('./constants');

const wss = new WebSocket.Server({
  port: 8082
});

// eslint-disable-next-line import/order
const readline = require('readline');

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout
});

wss.on('connection', ws => {
  ws.on('message', json => {
    console.log(json);
    const data = JSON.parse(json);
    if (data.type) {
      const client = [...wss.clients][0];
      if (handlers[data.type]) {
        handlers[data.type](data, (rep, event = events.reply) => {
          const reply = `event:${event};${JSON.stringify(rep)}`;
          client.send(reply);
          console.log('replying with', reply);
        });
      } else {
        console.warn(`unmet data type ${data.type}`);
      }
    }
    console.log('   ');
  });
});
rl.on('line', input => {
  const client = [...wss.clients][0];
  if (!client) {
    return;
  }
  if (handlers[input]) {
    handlers[input]({}, (rep, event = events.reply) => {
      const reply = `iris-websocket-message:${event};0;${JSON.stringify(rep)}`;
      console.log('replying with', reply);
      console.log('   ');
      client.send(reply);
    });
  }
});
