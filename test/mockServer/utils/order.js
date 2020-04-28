const { genStep, stepTypes } = require('./step');
const { setDir, write, list } = require('./storage');

const demoStepConfigs = [
  { testType: stepTypes.scanner },
  { testType: stepTypes.material, config: { consume_product: '1111' } },
  { testType: stepTypes.text },
  { testType: stepTypes.screw },
  { testType: stepTypes.passFail },
  { testType: stepTypes.measure }
];

setDir('../storage');
generateDemoOrder('demoOrder');

function generateDemoOrder(name) {
  const pointConfigs = list('order');
  const pointConfigsCount = pointConfigs.length;
  const points = genOrder(demoStepConfigs);
  write('order', name || `demoPoint_${pointConfigsCount + 1}`, JSON.stringify(points, null, '\t'));
}

function genOrder(stepConfigs, {
  status = 'todo',
  product_type_image = '',
  product_code = 'demo_product_code',
  track_code = '1x0001'
} = {}) {
  const now = new Date();
  const id = now.valueOf();
  const steps = stepConfigs.map((c, idx) => genStep(c.testType, idx + 1, c.config));
  return {
    id,
    code: `WO${id}`,
    track_code,
    product_code,
    date_planned_start: now.toDateString(),
    date_planned_complete: now.toDateString(),
    status,
    product_type_image,
    payload: {
      ...genOrderPayload(),
      ...orderPayloadTangChe()
    },
    steps
  };
}

function genOrderPayload() {
  return {
    "products": [{
      "url": "string",
      "code": "string"
    }],
    worksheet: {
      "name": "工艺指导书1",
      "revision": "v1",
      "url": "http://www.pdf995.com/samples/pdf.pdf"
    },
    environments: [{
      text: 'string',
      test_type: 'text',
      code: '1111',
      sequence: 1,
      desc: 'string'
    }],
    components: [{
      is_key: true,
      code: '1111'
    }],
    operation: {
      code: 'G0C',
      resources: {
        equipments: [
          'xx0011'
        ],
        users: [
          '112233'
        ]
      },
      desc: '测试工艺作业'
    },
    workcenter: {
      code: 'TA2-26L-01',
      locations: [
        {
          equipment_sn: '1234567',
          io_input: 1,
          io_output: 1,
          product_code: 'C1780589'
        }
      ]
    }
  };
}

function orderPayloadTangChe() {
  return {
    'SYSTEMTYPE': '3', // 系统类型
    'WIPORDERTYPE': '108', // 订单类型
    'MOMDISPOSITIONS': [
      {
        'DTID': 'null',
        'DISPOSITION': 'null',
        'LINESEQUENCENO': 'null',
        'SEQUENCENO': 'null',
        'SPECIFICATIONDESC': 'null',
        'CHARACTERISTIC': 'null',
        'CHARACTERISTICDESC': 'null',
        'CHARACTERISTICTYPE': 'null',
        'UPPERCONTROLLIMIT': 'null',
        'LOWERCONTROLLIMIT': 'null',
        'TARGETVALUE': 'null',
        'UOMCODE': 'null',
        'TESTVALUE': 'null',
        'MEANVALUE': 'null',
        'ATTRIBUTE': 'null'
      }
    ], // 订单类型
    'MOMCONFIG': {
      'PLABEL': 'null',
      'PPRODUCTNO': 'M003000000529',
      'PSERIALNO': 'null',
      'CHILDCONFIGS': {
        'CHILDCONFIG': [{
          'CGBOMID': 11652,
          'CPRODUCTID': 100224207,
          'CPRODUCTNO': 'CNR0000392139',
          'TRACKINGCODE': '3',
          'INSTALLATIONPOSITION': '',
          'INSTALLATIONDESC': '',
          'BOMGWDESC': '',
          'DUMMY1': '',
          'DUMMY2': '',
          'DUMMY3': '',
          'DUMMY4': '',
          'DUMMY5': '',
          'DUMMY6': ''
        }]
      }
    },
    'RESOURCEGROUP': 'null',
    'STARTEMPLOYEE': 'luogan',
    'RESOURCENAMES': {
      'RESOURCENAME': [
        'null',
        'null'
      ]
    },
    'PARENTWIPORDERNO': '083000166996',
    'PARENTWIPORDERTYPE': '101'
  };
}
