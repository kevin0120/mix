const { pointsData } = require('./points');

const tools = ['xx4443', 'xx0011'];

const demoPoints = [
  {
    [tools[0]]: true,
    [tools[1]]: true
  }, {
    [tools[0]]: true
  }
];

const stepTypes = {
  input: 'input',
  scanner: 'register_byproducts',
  instruction: 'text',
  text: 'text',
  screw: 'tightening',
  material: 'register_consumed_materials',
  passFail: 'passfail',
  measure: 'measure',
  video: 'video',
};

function genStep(testType, sequence, {
  skippable = false,
  undoable = false,
  data = '',
  status = 'ready',
  image = '',
  text='',
  consume_product = null
} = {}) {
  const now = new Date();
  const id = now.valueOf();
  const code = `${testType}-${id}`;
  const target = Math.round(Math.random() * (10 - 1) + 1);
  let payload;
  if (stepPayloads[testType]) {
    payload = stepPayloads[testType]();
  }
  return {
    id,
    code,
    sequence,
    // responsible_user: 'someone',
    desc: `${testType} 工步描述`,
    test_type: testType,
    text,
    failure_msg: `${code}-failed`,
    tolerance_min: target - 1,
    tolerance_max: target + 1,
    target,
    consume_product,
    skippable,
    undoable,
    data,
    status,
    image,
    payload
  };
}

const stepPayloads = {
  [stepTypes.screw]() {
    const points = pointsData(demoPoints);
    return {
      tightening_total: points.length,
      tightening_points: points
    };
  }
  // [stepTypes.text]() {
  //   return {
  //     text: '测试文本工步'
  //   };
  // }
};

module.exports = {
  stepTypes,
  genStep
};
