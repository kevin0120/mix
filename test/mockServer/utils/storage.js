const fs = require('fs'); // 引入fs模块

const storageDir = './test/mockServer/storage';

function write(key, name, data) {
  return new Promise((resolve, reject) => {
    fs.writeFile(`${storageDir}/${key}/${name}.json`, data, {}, (err) => {
      if (err) {
        reject(err);
      }
      resolve();
    });
  });
}

function read(key, name) {
  return new Promise((resolve, reject) => {
    fs.readFile(`${storageDir}/${key}/${name}.json`, 'utf8', (err, data) => {
      if (err) {
        reject(err);
      }
      resolve({
        key,
        name,
        data: JSON.parse(data.replace('\\n', ''))
      });
    });
  });
}

function list(key) {
  const files = fs.readdirSync(`${storageDir}/${key}/`);
  return files.map(f => f.split('.')[0]);
}

module.exports = {
  read, write, list
};
