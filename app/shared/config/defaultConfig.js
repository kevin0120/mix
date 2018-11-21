
export const defaultConfigs = {
  version: 'v0.1',
  base: {
    userInfo: {
      uuid: '1234'
    },
    psetPointDiameter: 30 // 编辑页面点位大小,未来不使用
  },
  page: {
    // 配置页面不同的导航页
    network: {
      ssid: {
        displayOrder: 1,
        value: '',
        displayTitle: 'Configuration.network.SSID'
      },
      password: {
        displayOrder: 100,
        value: '',
        displayTitle: 'Configuration.network.PWD'
      },
      ipAddress: {
        displayOrder: 200,
        value: '192.168.1.5',
        displayTitle: 'Configuration.network.Addr'
      },
      netmask: {
        displayOrder: 300,
        value: '255.255.255.0',
        displayTitle: 'Configuration.network.Mask'
      },
      gateway: {
        displayOrder: 400,
        value: '192.168.1.1',
        displayTitle: 'Configuration.network.Gateway'
      }
    },
    odooConnection: {
      odooUrl: {
        displayOrder: 1,
        value: 'http://192.168.4.3:8069/api/v1',
        displayTitle: 'Configuration.connections.Odoo'
      },
      hmiSn: {
        displayOrder: 2,
        value: '1122334455667788',
        displayTitle: 'Configuration.connections.HMI'
      },
      aiisUrl: {
        displayOrder: 3,
        value: 'http://192.168.4.3:9092/aiis/v1',
        displayTitle: 'Configuration.connections.AIIS'
      }
    },
    modbus: {
      in: [
        {
          bit: 0,
          io: 'in',
          function: '',
          label: ''
        },
        {
          bit: 1,
          io: 'in',
          function: '',
          label: ''
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
          function: '',
          label: ''
        },
        {
          bit: 4,
          io: 'in',
          function: 'RESET',
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
      out: [
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
          function: 'LED_BEEP',
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
    device: {
      scanner: {
        vendorId: 3118,
        mode: 'HID' // HID or BT_HID
        // vendorId: 1504
      },
      rfid: {
        host: '192.168.1.120',
        port: 2112
      }
    }
  },
  // 作业配置
  operationSettings: {
    opMode: 'order', // 作业模式:        op 或 order
    controllerMode: 'pset', // 拧紧模式:        job 或 pset
    workMode: 'auto', // 工作模式:        auto 或 scanner 或 manual
    flowTriggers: ['carID'], // 工作流程触发条件:  carType:车型代码 carID:vin/knr/longpin

    // 作业前检测(order mode only)
    preCheck: false,

    // 强制放行配置
    byPass: {
      enable: true,
      type: 'sleep' // sleep or press
    },

    // 结果对话框
    enableResultDialog: true
  },

  systemSettings: {
    enableDebugInfo: false,
    showSwitchMode: false, // 切换pset or job 模式
    defaultControllerSN: '0001',
    authEnable: true,
    switchAutoManual: false,
    oeeFuncEnable: false,
    modbusEnable: true,
    rfidEnabled: false,
    psetContinueMode: false
  }
};
