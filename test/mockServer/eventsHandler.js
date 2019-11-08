const { orderHandlers } = require('./dataHandlers/order');
const { deviceHandlers } = require('./dataHandlers/device');
const { scannerHandlers } = require('./dataHandlers/scanner');
const { readerHandlers } = require('./dataHandlers/reader');
const { resultHandlers } = require('./dataHandlers/result');
const { registerHandlers } = require('./dataHandlers/register');
const { toolHandlers } = require('./dataHandlers/tool');

module.exports = {
  ...orderHandlers,
  ...scannerHandlers,
  ...readerHandlers,
  ...deviceHandlers,
  ...resultHandlers,
  ...registerHandlers,
  ...toolHandlers
};
