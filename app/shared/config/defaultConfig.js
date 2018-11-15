
export const defaultConfigs = {
  version: 'v0.1',
  base:{
    userInfo: {
      uuid: '1234'
    },
    psetPointDiameter: 30 // 编辑页面点位大小,未来不使用
  },
  page:{
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
          function: 'RESET_KEY',
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
          function: 'LED_GREEN',
          label: '绿灯'
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
          function: 'LED_RED',
          label: '红灯'
        },
        {
          bit: 3,
          io: 'out',
          function: 'LED_BLUE',
          label: '蓝灯'
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
    },
  },
  system:{
    device: {
      scanner: {
        vendorId: 3118,
        mode: 'HID', // HID or BT_HID
        // vendorId: 1504
      },
      rfid: {
        host: '192.168.1.120',
        port: 2112
      },
    },
  },
  workflow:{
    name: "VW", // 工作流有: General or VW,
    bypass: {
      enable: true,
      output: 4, // 输出,输入哪一位
      wait: 5000
    },
    manualByPanel: {
      enable: false,
      input: 4, // 前几位
      wait: 5000
    },
    byPass: {
      enable: true,
      byPassJob: 250,
      // sleep/press
      byPassType: 'sleep',
      byPassPress: 2
    }
  },
  systemSettings: {
    enableDebugInfo: false,
    showSwitchMode: false, // 切换pset or job 模式
    defaultControllerSN: "0001",
    authEnable: true,
    switchAutoManual: false,
    oeeFuncEnable: false,
    modbusEnable: true,
    rfidEnabled: false,
    psetContinueMode: false,
  },
};
