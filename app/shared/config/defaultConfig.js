// @flow

const defaultConfigs = {
  version: 'v0.1',
  devices: {
    io: {
      enable: true,

    }
  },
  pages: {
    app: {
      working: ['user', 'admin'],
      // viewer: ['user', 'admin'],
      // order: ['user', 'admin'],
      preference: {
        Net:['user', 'admin'],
        IO:['user', 'admin'],
        Connect:['user', 'admin'],
        help:['user','admin']
      },
      event: ['user', 'admin'],
      // result: ['user', 'admin'],
      // curve: ['user', 'admin'],
      // help: ['user', 'admin']
    },
    pages: {
      login: ['user', 'admin']
    }
  },
  authorization: {
    localUsers: {
      'userName': {
        password: 'password',
        role: 'admin',
        uid: 'demoUID',
        uuid: 'uuid',
        avatar: ''
      },
      'ming': {
        password: '123',
        role: 'admin',
        uid: 'XiaoMing',
        uuid: '87d1c74e9000',
        avatar: ''
      },
      'hong': {
        password: '123',
        role: 'admin',
        uid: 'XiaoHong',
        uuid: '677e0f4f9000',
        avatar: ''
      },
    },
    // verify: 'local',//'online
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
        value:
          '192.168.1.5',
        displayTitle:
          'Configuration.network.Addr',
        isPWD:
          false
      }
      ,
      netmask: {
        displayOrder: 300,
        value:
          '255.255.255.0',
        displayTitle:
          'Configuration.network.Mask',
        isPWD:
          false
      }
      ,
      gateway: {
        displayOrder: 400,
        value:
          '192.168.1.1',
        displayTitle:
          'Configuration.network.Gateway',
        isPWD:
          false
      }
    }
    ,
    odooConnection: {
      odooUrl: {
        displayOrder: 1,
        value:
          'http://192.168.127.118:8069/api/v1',
        displayTitle:
          'Configuration.connections.Odoo'
      }
      ,
      hmiSn: {
        displayOrder: 2,
        value:
          '1560c527ac7f4f9f90e4900b50717456',
        displayTitle:
          'Configuration.connections.HMI'
      }
    }
    ,
    modbus: {
      in:
        [
          {
            bit: 0,
            io: 'in',
            function: '',
            label: '模式选中'
          },
          {
            bit: 1,
            io: 'in',
            function: '',
            label: '强制放行'
          },
          {
            bit: 2,
            io: 'in',
            function: '',
            label: ''
          },
          {
            bit: 3,
            io: 'in',
            function: 'BYPASS',
            label: ''
          },
          {
            bit: 4,
            io: 'in',
            function: 'MODE_SELECT',
            label: '复位钥匙'
          },
          {
            bit: 5,
            io: 'in',
            function: '',
            label: ''
          },
          {
            bit: 6,
            io: 'in',
            function: '',
            label: ''
          },
          {
            bit: 7,
            io: 'in',
            function: '',
            label: ''
          }
        ],
      out:
        [
          {
            bit: 0,
            io: 'out',
            function: 'LED_WHITE',
            label: '白灯'
          },
          {
            bit: 1,
            io: 'out',
            function: 'LED_YELLOW',
            label: '黄灯'
          },
          {
            bit: 2,
            io: 'out',
            function: 'LED_GREEN',
            label: '绿灯'
          },
          {
            bit: 3,
            io: 'out',
            function: 'LED_RED',
            label: '红灯'
          },
          {
            bit: 4,
            io: 'out',
            function: 'BEEP',
            label: '蜂鸣'
          },
          {
            bit: 5,
            io: 'out',
            function: '',
            label: ''
          },
          {
            bit: 6,
            io: 'out',
            function: '',
            label: ''
          },
          {
            bit: 7,
            io: 'out',
            function: '',
            label: ''
          }
        ]
    }
  },
  system: {
    workcenter: {
      code: '1',
      rework_workcenter:
        'qrk',
      hardware:
        {
          scanner: {
            vendorId: 3118,
            mode:
              'HID' // HID or BT_HID
            // mode: 'BT_HID',
            // vendorId: 1504
          },
          io: {
            serialNo: '111',
            connection:
              'modbustcp://192.168.1.122:502/0'
          },
          controllers: [
            {
              serialNo: 'b2a8d23b44ad4836be0be8ebca72b8a1',
              guns: [
                { serialNo: '1111' }
              ]
            }
          ]
        }

    },
    connections: {
      rush: 'http://192.168.4.146:8082',// 4.188//4.219//192.168.4.188//10.1.1.65//192.168.4.247//0.2
      // rush: 'http://0da41704.ngrok.io:8082',//4.188//4.219//192.168.4.188//10.1.1.65//192.168.4.247//0.2
      rfid: 'tcp://127.0.0.1:2112',
      aiis: 'http://127.0.0.1:9092'
    }
  },
// 作业配置
  operationSettings: {
    opMode: 'op', // 作业模式:        op 或 order
    controllerMode:
      'job', // 拧紧模式:        job 或 pset
    workMode:
      'auto', // 工作模式:        auto 或 manual 或 scanner
    flowTriggers:
      ['carID', 'carType'], // 工作流程触发条件:  carType:车型代码 carID:vin/knr/longpin

    // 作业前检测(order mode only)
    preCheck:
      false, // 开工检查

    // 强制放行配置
    byPass:
      {
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
    curveEnable: true
  }
};

export default defaultConfigs;
