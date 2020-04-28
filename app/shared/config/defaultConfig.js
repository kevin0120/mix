// @flow
import { ioInputs, ioOutputs } from '../../modules/io/constants';

const defaultConfigs = {
  version: 'v0.1',
  // appName: '智能装配系统-终端',
  appName: '唐山客车智能装配系统-终端',
  devices: {
    io: {
      sn: '11',
      enable: true,
      config: { input_num: 0, output_num: 16 },
      inputs: {
        [ioInputs.resetKey]: 0,
        [ioInputs.byPass]: 1,
        [ioInputs.modeSelect]: 2
      },
      outputs: {
        [ioOutputs.white]: 0,
        [ioOutputs.yellow]: 1,
        [ioOutputs.green]: 2,
        [ioOutputs.red]: 3,
        [ioOutputs.unlock]: 13
      }
    }
  },

  adminKey: {
    enable: false,
    io_sn: '11',
    input: 3
  },
  pages: {
    app: {
      // manual: ['user', 'admin'],
      working: ['user', 'admin'],
      manualwork: ['user', 'admin'],
      // viewer: ['user', 'admin'],
      // order: ['user', 'admin'],
      preference: {
        __role: ['admin'],
        Net: ['user', 'admin'],
        IO: ['user', 'admin'],
        Connect: ['user', 'admin'],
        help: ['user', 'admin']
      },
      event: ['admin'],
      result: ['admin']
      // curve: ['user', 'admin'],
      // help: ['user', 'admin']
    },
    pages: {
      login: ['user', 'admin']
    }
  },
  authorization: {
    localUsers: {
      userName: {
        password: 'password',
        role: 'admin',
        uid: '1',
        uuid: '1',
        avatar: ''
      },
      admin: {
        password: 'admin',
        role: 'admin',
        uid: '1',
        uuid: '1',
        avatar: ''
      }
    },
    verify: 'online', // 'online'
    maxUsers: 0 // 0:no limit
  },

  page: {
    // 配置页面不同的导航页
    network: {
      ssid: {
        displayOrder: 1,
        value: '',
        displayTitle: 'Configuration.network.SSID',
        isPWD: false
      },
      password: {
        displayOrder: 100,
        value: '',
        displayTitle: 'Configuration.network.PWD',
        isPWD: true
      },
      ipAddress: {
        displayOrder: 200,
        value: '192.168.1.5',
        displayTitle: 'Configuration.network.Addr',
        isPWD: false
      },
      netmask: {
        displayOrder: 300,
        value: '255.255.255.0',
        displayTitle: 'Configuration.network.Mask',
        isPWD: false
      },
      gateway: {
        displayOrder: 400,
        value: '192.168.1.1',
        displayTitle: 'Configuration.network.Gateway',
        isPWD: false
      }
    },
    odooConnection: {
      odooUrl: {
        displayOrder: 1,
        value: 'http://10.1.1.47:8069/api/v1',
        // value: 'http://127.0.0.1:8069/api/v1',
        displayTitle: 'Configuration.connections.Odoo'
      },
      hmiSn: {
        displayOrder: 2,
        value: '1560c527ac7f4f9f90e4900b50717456',
        displayTitle: 'Configuration.connections.HMI'
      }
    }
  },
  system: {
    workcenter: {
      // code: '1',
      workcenterType: 'normal', // 'normal', 'rework'
      rework_workcenter: 'qrk',
      hardware: {
        scanner: {
          vendorId: 3118,
          mode: 'HID' // HID or BT_HID
          // mode: 'BT_HID',
          // vendorId: 1504
        },
        io: {
          serialNo: '111',
          connection: 'modbustcp://192.168.1.122:502/0'
        },
        controllers: [
          {
            serialNo: 'b2a8d23b44ad4836be0be8ebca72b8a1',
            guns: [{ serialNo: '1111' }]
          }
        ]
      }
    },
    connections: {
      // rush: 'http://192.168.4.32:8082',
      // rush: 'http://192.168.4.18:8082',
      // rush: 'http://192.168.3.118:8082',
      // rush: 'http://192.168.4.96:8082',
      // rush: 'http://192.168.5.2:8082',
      rush: 'http://127.0.0.1:8082',
      // rush: 'http://10.1.1.47:8082',
      // rush: 'http://192.168.3.240:8082',
      // rush: 'http://0da41704.ngrok.io:8082',//4.188//4.219//192.168.4.188//10.1.1.65//192.168.4.247//0.2
      rfid: 'tcp://127.0.0.1:2112',
      aiis: 'http://127.0.0.1:9092'
    },
    endpoints: {
      operationResources:
        'http://172.26.214.80:8091/a/login2authc?type=mtShow&userId=5f19fbe973a44b99826488b3be02d41a&orderNumber=0000198338122', // 生产过程MES
      operationTracing:
        'http://172.26.214.80:8091/a/product/traceable?orderNumber=0000198338122' // 生产信息追溯
    }
  },
  // 作业配置
  operationSettings: {
    opMode: 'op', // 作业模式:        op 或 order
    controllerMode: 'job', // 拧紧模式:        job 或 pset
    workMode: 'auto', // 工作模式:        auto 或 manual 或 scanner
    flowTriggers: ['carID', 'carType'], // 工作流程触发条件:  carType:车型代码 carID:vin/knr/longpin

    // 作业前检测(order mode only)
    preCheck: false, // 开工检查

    maxReworkCnt: 1, // 默认的最大返工点计数

    // 强制放行配置
    byPass: {
      enable: true,
      type: 'sleep' // sleep or press
    },

    // 空车job
    emptyCarJob: 250,

    // 结果对话框
    enableResultDialog: true,

    // 手动模式下是否启用freestyle
    manualFreestyle: true,

    // 启用ak2
    enableAk2: true,

    regExp: '(C6\\d{12})' // rfid正則表達式,SVW
  },
  cvinetweb: {
    url: 'http://192.168.1.10:8080/CVINetWeb'
  },
  systemSettings: {
    restrictMode: true, // 严格模式，关键点必须为pass
    enableDebugInfo: false,
    showSwitchMode: false, // 切换pset or job 模式
    authEnable: true, // 认证
    switchAutoManual: false,
    oeeFuncEnable: false,
    modbusEnable: false,
    rfidEnabled: false,
    andonEnable: false,
    enableFocus: false,
    viewerEnable: true,
    viewer: {
      image: 'editable',
      file: true
    },
    curveEnable: true,
    reportStart: true, // 报开工
    reportFinish: true, // 报完工
    orderSimulate: false, // 产前模拟
    strictOrderSimulate: false, // true:严格产前模拟，不通过无法开工
    strictReportStart: false, // true:严格报开工，不通过无法开工
    strictResume: false, // true:严格恢复，不通过无法恢复作业
    canDoAnotherStep: true,
    canRedoOrders: true
  },
  debugSettings: {
    disableOrderTriggerLimit: false
  }
};

export default defaultConfigs;
