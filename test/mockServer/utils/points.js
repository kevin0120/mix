const { setDir, write, list } = require('./storage');

const tools = ['xx4443', 'xx0011'];

const demoPoints = [
  {
    [tools[0]]: true,
    [tools[1]]: true
  }, {
    [tools[0]]: true,
    [tools[1]]: true
  }, {
    [tools[0]]: true,
    [tools[1]]: false
  }
];

function generatePoints(name) {
  const pointConfigs = list('points');
  const pointConfigsCount = pointConfigs.length;
  const points = pointsData(demoPoints);
  write('points', name || `demoPoint_${pointConfigsCount + 1}`, JSON.stringify(points));
}

setDir('../storage');
generatePoints();

function pointsData(groupConfigs) {
  const points = [];
  let seq = 1;
  groupConfigs.forEach((gc, idx) => {
    const keyNum = Object.values(gc).reduce((total, c) => {
      return c ? total + 1 : total;
    }, 0);
    Object.keys(gc).forEach((tool, pIdx) => {
      const isKey = gc[tool];
      if (!tool) {
        throw new Error(`invalid tool : ${idx}, ${pIdx}`);
      }
      const tempPoints = [...points].reverse();
      const point = tempPoints.find(p => p.tightening_tool === tool);
      if (point && (!point.is_key || Object.keys(groupConfigs[point.group_sequence]) > point.key_num)) {
        throw new Error(`previous point of tool(${tool}) is not key, group: ${idx}, point: ${pIdx
          }, prev group: ${point.group_sequence - 1}`
        );
      }
      points.push({
        'is_key': isKey,
        'pset': '1',
        'sequence': seq,
        'y': (idx + 1) * 10,
        'x': (pIdx + 1) * 10,
        'key_num': keyNum,
        'tightening_tool': tool,
        'max_redo_times': 3,
        'group_sequence': idx + 1,
        'nut_no': ''
      });
      seq += 1;
    });
  });
  return points;
}
