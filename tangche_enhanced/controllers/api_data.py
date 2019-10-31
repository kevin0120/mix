# -*- coding: utf-8 -*-


DEFAULT_LIMIT = 80

api_data = {
  "info": {
    "termsOfService": "http://centronsys.com",
    "version": "1.0.0",
    "title": "智能装配应用服务器RESTful",
    "description": "智能装配应用服务器RESTful",
    "contact": {
      "email": "gubin@centronsys.com"
    }
  },
  "paths": {
    "/res.users": {
      "get": {
        "responses": {
          "200": {
            "description": "用户清单",
            "schema": {
              "items": {
                "$ref": "#/definitions/User"
              },
              "type": "array"
            }
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "items": {
              "type": "string"
            },
            "type": "array",
            "description": "UUID to filter by",
            "name": "uuids",
            "in": "query"
          },
          {
            "description": "返回结果限定个数",
            "default": 80,
            "required": False,
            "name": "limit",
            "in": "query",
            "type": "integer",
            "collectionFormat": "multi"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Users"
        ],
        "summary": "查询用户清单",
        "consumes": [
          "application/json"
        ],
        "description": ""
      }
    },
    "/ts002/workorders": {
      "post": {
        "responses": {
          "201": {
            "description": "成功创建工单下发给工位控制器",
            "schema": {
              "items": {
                "$ref": "#/definitions/ACKRESP"
              }
            }
          },
          "400": {
            "description": "下发参数有误",
            "schema": {
              "items": {
                "$ref": "#/definitions/ACKRESP"
              }
            }
          }
        },
        "description": "唐客下发工单",
        "parameters": [
          {
            "required": True,
            "in": "body",
            "name": "body",
            "schema": {
              "$ref": "#/definitions/ts002Order"
            }
          }
        ],
        "tags": [
          "TS002"
        ]
      }
    },
    "/maintenance/requests/try": {
      "post": {
        "responses": {
          "201": {
            "description": "创建维护请求成功",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          },
          "204": {
            "description": "无需创建维护请求,未到达下次创建的条件"
          },
          "404": {
            "description": "枪序列号未找到",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          },
          "405": {
            "description": "无效的输入,缺少参数",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          }
        },
        "parameters": [
          {
            "required": True,
            "in": "body",
            "name": "body",
            "schema": {
              "$ref": "#/definitions/try_maintenance_request"
            }
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Maintenance"
        ],
        "summary": "根据发送拧紧次数创建维护请求",
        "consumes": [
          "application/json"
        ],
        "description": "根据发送拧紧次数创建维护请求"
      }
    },
    "/operation.results/{resultId}/curves_add": {
      "patch": {
        "responses": {
          "200": {
            "description": "成功更新了结果数据",
            "schema": {
              "items": {
                "$ref": "#/definitions/ResultDetail"
              }
            }
          },
          "404": {
            "description": "resultId 未找到"
          }
        },
        "parameters": [
          {
            "description": "需要更新的结果的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "resultId"
          },
          {
            "schema": {
              "$ref": "#/definitions/curve"
            },
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Result"
        ],
        "summary": "为一条结果添加波形",
        "consumes": [
          "application/json"
        ],
        "description": "更新一条拧紧结果数据"
      }
    },
    "/hmi.connections/{serial_no}": {
      "get": {
        "responses": {
          "200": {
            "description": "连接信息清单",
            "schema": {
              "$ref": "#/definitions/hmi_connections"
            }
          },
          "404": {
            "description": "uuid not found"
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "description": "HMI唯一标示",
            "default": 1122334455667788,
            "required": True,
            "in": "path",
            "type": "string",
            "name": "serial_no"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "HMI"
        ],
        "summary": "查询HMI连接信息",
        "consumes": [
          "application/json"
        ],
        "description": "获取hmi信息"
      }
    },
    "/mrp.productions/{vin}": {
      "get": {
        "responses": {
          "200": {
            "description": "生产订单",
            "schema": {
              "$ref": "#/definitions/Production"
            }
          },
          "404": {
            "description": "vin not found"
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "required": True,
            "type": "string",
            "description": "VIN",
            "name": "vin",
            "in": "path"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Manufacture"
        ],
        "consumes": [
          "application/json"
        ],
        "description": "获取某一用户信息"
      }
    },
    "/mrp.workorders": {
      "get": {
        "responses": {
          "200": {
            "description": "获取工单",
            "schema": {
              "items": {
                "$ref": "#/definitions/WorkOrder"
              },
              "type": "array"
            }
          },
          "404": {
            "description": "MasterPC not found"
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "required": False,
            "type": "string",
            "description": "MasterPC UUID",
            "name": "masterpc",
            "in": "query"
          },
          {
            "required": False,
            "type": "string",
            "description": "HMi序列号",
            "name": "hmi",
            "in": "query"
          },
          {
            "required": False,
            "type": "string",
            "description": "工位序列号",
            "name": "workcenter",
            "in": "query"
          },
          {
            "required": False,
            "type": "string",
            "description": "Long PIN or VIN or KNR",
            "name": "code",
            "in": "query"
          },
          {
            "description": "返回结果的条数限制",
            "default": 80,
            "required": False,
            "in": "query",
            "type": "integer",
            "name": "limit"
          },
          {
            "description": "返回结果的排序字段,默认是降序",
            "default": "production_date",
            "required": False,
            "in": "query",
            "type": "string",
            "name": "order"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Manufacture"
        ],
        "consumes": [
          "application/json"
        ],
        "description": "获取某一用户信息"
      }
    },
    "/operation.results": {
      "put": {
        "responses": {
          "204": {
            "description": "成功更新了结果数据"
          },
          "404": {
            "description": "未找到记录"
          },
          "405": {
            "description": "无效的输入"
          }
        },
        "parameters": [
          {
            "schema": {
              "items": {
                "$ref": "#/definitions/Batchresult"
              },
              "type": "array"
            },
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Result"
        ],
        "summary": "批量修改结果数据",
        "consumes": [
          "application/json"
        ],
        "description": "获取拧紧结果数据"
      },
      "get": {
        "responses": {
          "200": {
            "description": "成功获取结果数据",
            "schema": {
              "items": {
                "$ref": "#/definitions/ResultDetail"
              }
            }
          },
          "404": {
            "description": "未找到记录"
          },
          "405": {
            "description": "无效的输入"
          }
        },
        "parameters": [
          {
            "description": "查询起始时间",
            "format": "date-time",
            "default": "2017-06-19T16:39:57-08:00",
            "in": "query",
            "type": "string",
            "name": "date_from"
          },
          {
            "description": "查询终止时间",
            "format": "date-time",
            "default": "2018-12-29T16:39:57+08:00",
            "in": "query",
            "type": "string",
            "name": "date_to"
          },
          {
            "default": 80,
            "type": "integer",
            "name": "limit",
            "in": "query"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Result"
        ],
        "summary": "获取结果数据",
        "consumes": [
          "application/json"
        ],
        "description": "获取拧紧结果数据"
      }
    },
    "/logo": {
      "get": {
        "responses": {
          "200": {
            "description": "图片信息",
            "schema": {
              "$ref": "#/definitions/Logo"
            }
          }
        },
        "tags": [
          "LOGO"
        ]
      }
    },
    "/maintenance/requests": {
      "post": {
        "responses": {
          "201": {
            "description": "创建维护请求成功",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          },
          "404": {
            "description": "枪序列号未找到",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          },
          "405": {
            "description": "无效的输入,缺少参数",
            "schema": {
              "$ref": "#/definitions/common_resp"
            }
          }
        },
        "parameters": [
          {
            "required": True,
            "in": "body",
            "name": "body",
            "schema": {
              "$ref": "#/definitions/maintenance_request"
            }
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Maintenance"
        ],
        "summary": "创建维护请求",
        "consumes": [
          "application/json"
        ],
        "description": "创建维护请求"
      }
    },
    "/mrp.workorders/{order_id}": {
      "get": {
        "responses": {
          "200": {
            "description": "获取工单",
            "schema": {
              "items": {
                "$ref": "#/definitions/WorkOrder"
              },
              "type": "array"
            }
          },
          "404": {
            "description": "MasterPC not found"
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "required": True,
            "type": "string",
            "description": "MasterPC UUID",
            "name": "order_id",
            "in": "path"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Manufacture"
        ],
        "consumes": [
          "application/json"
        ],
        "description": "获取某一用户信息"
      }
    },
    "/mrp.routing.workcenter/{operation_id}": {
      "get": {
        "description": "获取作业详情",
        "responses": {
          "200": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OperationDetail"
            }
          }
        },
        "parameters": [
          {
            "description": "需要获取的作业的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "operation_id"
          }
        ],
        "tags": [
          "Operation"
        ]
      }
    },
    "/mrp.routing.workcenter/{operation_id}/edit": {
      "put": {
        "description": "通过此api可以对某个作业的图片和点位进行编辑。参数中图片传入base64编码的字符串。 对于点位列表,如果传入的点位有sequence则会进行更新（如果对应sequence的点位不存在则会新增），如果不传sequence则会新增。对于那些不在点位列表中的作业点会进行删除操作（比如当前某作业有3个点，调用api时只传了前两个，那么第三个点会被删除。同理如果传了一个空的点位列表则会删除该作业的所有点）",
        "responses": {
          "200": {
            "description": "成功更新了作业图片点位"
          },
          "405": {
            "description": "无效参数"
          }
        },
        "parameters": [
          {
            "description": "需要修改的作业的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "operation_id"
          },
          {
            "schema": {
              "$ref": "#/definitions/operation_edit"
            },
            "description": "图片点位信息",
            "name": "body",
            "in": "body"
          }
        ],
        "tags": [
          "Operation"
        ]
      }
    },
    "/res.users/batch_archived": {
      "put": {
        "responses": {
          "200": {
            "description": "用户清单",
            "schema": {
              "items": {
                "$ref": "#/definitions/User"
              },
              "type": "array"
            }
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "required": True,
            "in": "body",
            "description": "用户唯一标示(胸卡信息)",
            "name": "body",
            "schema": {
              "items": {
                "type": "string",
                "example": "112233"
              },
              "type": "array"
            }
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Users"
        ],
        "summary": "批量归档用户",
        "consumes": [
          "application/json"
        ],
        "description": "批量归档用户"
      }
    },
    "/operation.results/{resultId}": {
      "put": {
        "responses": {
          "204": {
            "description": "成功更新了结果数据"
          },
          "404": {
            "description": "resultId 未找到"
          }
        },
        "parameters": [
          {
            "description": "需要更新的结果的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "resultId"
          },
          {
            "schema": {
              "$ref": "#/definitions/result"
            },
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Result"
        ],
        "summary": "更新一条结果数据",
        "consumes": [
          "application/json"
        ],
        "description": "更新一条拧紧结果数据"
      },
      "get": {
        "responses": {
          "200": {
            "description": "成功更新了结果数据",
            "schema": {
              "items": {
                "$ref": "#/definitions/ResultDetail"
              }
            }
          },
          "404": {
            "description": "resultId 未找到"
          }
        },
        "parameters": [
          {
            "description": "需要更新的结果的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "resultId"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Result"
        ],
        "summary": "获取一条结果数据",
        "consumes": [
          "application/json"
        ],
        "description": "获取一条拧紧结果数据"
      }
    },
    "/mrp.productions": {
      "post": {
        "description": "当AIIS收到FIS下发的装配任务，会调用此API将任务同步下发给ODOO.",
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "schema": {
              "$ref": "#/definitions/mission"
            }
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "AIIS"
        ],
        "summary": "下发装配任务",
        "consumes": [
          "application/json"
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/Production"
            }
          },
          "204": {
            "description": "生产订单已存在"
          },
          "400": {
            "description": "失败",
            "schema": {
              "$ref": "#/definitions/ResponseBody"
            }
          }
        }
      },
      "get": {
        "description": "获取生产订单清单",
        "parameters": [
          {
            "items": {
              "default": "LSV2A8CA7JN508198",
              "type": "string"
            },
            "type": "array",
            "name": "vins",
            "in": "query"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Manufacture"
        ],
        "summary": "获取生产订单清单",
        "consumes": [
          "application/json"
        ],
        "responses": {
          "200": {
            "description": "成功",
            "schema": {
              "items": {
                "$ref": "#/definitions/Production"
              },
              "type": "array"
            }
          },
          "404": {
            "description": "未找到",
            "schema": {
              "$ref": "#/definitions/ResponseBody"
            }
          }
        }
      }
    },
    "/res.users/{uuid}": {
      "get": {
        "responses": {
          "200": {
            "description": "用户清单",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "404": {
            "description": "uuid not found"
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "parameters": [
          {
            "required": True,
            "type": "string",
            "description": "用户唯一标示(胸卡信息)",
            "name": "uuid",
            "in": "path"
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Users"
        ],
        "summary": "查询用户清单",
        "consumes": [
          "application/json"
        ],
        "description": "获取某一用户信息"
      }
    },
    "/mrp.routing.workcenter": {
      "get": {
        "description": "获取作业清单",
        "tags": [
          "Operation"
        ],
        "responses": {
          "200": {
            "description": "成功更新了结果数据",
            "schema": {
              "items": {
                "$ref": "#/definitions/Operation"
              }
            }
          }
        },
        "produces": [
          "application/json"
        ]
      }
    }
  },
  "host": "127.0.0.1",
  "schemes": [
    "http"
  ],
  "tags": [
    {
      "name": "Operation",
      "description": "工序作业"
    },
    {
      "name": "Result",
      "description": "工序操作结果"
    },
    {
      "name": "HMI",
      "description": "HMI设备"
    },
    {
      "name": "Manufacture",
      "description": "生产相关"
    },
    {
      "name": "Users",
      "description": "用户相关"
    },
    {
      "name": "Maintenance",
      "description": "设备维护"
    },
    {
      "name": "TS002",
      "description": "唐山客车项目特有API"
    }
  ],
  "definitions": {
    "Batchresult": {
      "type": "object",
      "properties": {
        "pset_m_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定扭矩阈值"
        },
        "pset_m_max": {
          "type": "number",
          "description": "设定最大扭矩",
          "example": 4.34
        },
        "control_date": {
          "example": "2018-05-19T16:39:57+08:00",
          "type": "string",
          "description": "拧紧时间",
          "format": "date-time"
        },
        "pset_w_max": {
          "type": "number",
          "description": "设定最大角度",
          "example": 4.34
        },
        "user_id": {
          "type": "integer",
          "example": 1,
          "description": "当前操作用户"
        },
        "one_time_pass": {
          "type": "string",
          "example": "fail",
          "description": "是否一次成功（pass/fail）"
        },
        "pset_strategy": {
          "enum": [
            "AD",
            "AW",
            "ADW",
            "LN",
            "AN",
            "AT"
          ],
          "type": "string",
          "description": "拧紧枪策略",
          "example": "AD"
        },
        "pset_w_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定角度阈值"
        },
        "cur_objects": {
          "items": {
            "$ref": "#/definitions/curve"
          },
          "type": "array",
          "description": "个次操作波形对象列表"
        },
        "pset_m_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标扭矩"
        },
        "pset_m_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小扭矩"
        },
        "final_pass": {
          "type": "string",
          "example": "pass",
          "description": "是否最终成功（pass/fail）"
        },
        "measure_degree": {
          "type": "number",
          "description": "实际角度",
          "example": 2.44
        },
        "measure_t_don": {
          "type": "number",
          "description": "拧紧过程花费时间",
          "example": 3.22
        },
        "measure_torque": {
          "type": "number",
          "description": "实际扭矩",
          "example": 3.224
        },
        "id": {
          "type": "integer",
          "example": 11635,
          "description": "修改了结果的id"
        },
        "measure_result": {
          "enum": [
            "ok",
            "nok"
          ],
          "type": "string",
          "description": "测量结果",
          "example": "ok"
        },
        "op_time": {
          "type": "integer",
          "description": "当前操作",
          "example": 1
        },
        "pset_w_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小角度"
        },
        "pset_w_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标角度"
        }
      }
    },
    "ACKRESP": {
      "type": "object",
      "properties": {
        "msg": {
          "type": "string",
          "description": "响应返回信息"
        },
        "extra": {
          "type": "string",
          "description": "响应额外返回信息"
        }
      }
    },
    "hmi_connections": {
      "type": "object",
      "properties": {
        "info": {
          "$ref": "#/definitions/HMILocationInfo"
        },
        "masterpc": {
          "$ref": "#/definitions/EquipConnection"
        },
        "rfid": {
          "$ref": "#/definitions/EquipConnection"
        },
        "controllers": {
          "items": {
            "$ref": "#/definitions/EquipConnection"
          },
          "type": "array"
        },
        "io": {
          "$ref": "#/definitions/EquipConnection"
        }
      }
    },
    "point": {
      "type": "object",
      "properties": {
        "y_offset": {
          "type": "integer",
          "description": "上偏移"
        },
        "x_offset": {
          "type": "integer",
          "description": "左偏移"
        },
        "sequence": {
          "type": "integer",
          "description": "序号"
        }
      }
    },
    "ts002OperationResource": {
      "required": [
        "users",
        "equipments"
      ],
      "type": "object",
      "description": "作业相关资源",
      "properties": {
        "equipments": {
          "items": {
            "type": "string",
            "example": "xx0011"
          },
          "type": "array",
          "description": "作业相关使用的工具"
        },
        "users": {
          "items": {
            "type": "string",
            "example": "112233"
          },
          "type": "array",
          "description": "作业相关操作人员"
        }
      }
    },
    "ts002FinishedProduct": {
      "required": [
        "code",
        "url"
      ],
      "type": "object",
      "description": "产成品",
      "properties": {
        "url": {
          "type": "string",
          "description": "产成品三维模型路径"
        },
        "code": {
          "type": "string",
          "description": "产成品料号"
        }
      }
    },
    "point_edit": {
      "type": "object",
      "properties": {
        "y_offset": {
          "type": "integer",
          "description": "上偏移"
        },
        "x_offset": {
          "type": "integer",
          "description": "左偏移"
        },
        "sequence": {
          "type": "integer",
          "description": "id"
        }
      }
    },
    "ts002Operation": {
      "required": [
        "code",
        "resources"
      ],
      "type": "object",
      "description": "工艺作业",
      "properties": {
        "code": {
          "type": "string",
          "description": "工艺代码",
          "example": "G0C"
        },
        "resources": {
          "$ref": "#/definitions/ts002OperationResource"
        },
        "desc": {
          "type": "string",
          "description": "工艺相关描述",
          "example": "测试工艺作业"
        }
      }
    },
    "mission": {
      "type": "object",
      "properties": {
        "pin": {
          "type": "number",
          "description": "订单车身pin码",
          "example": 6473537
        },
        "vin": {
          "type": "string",
          "description": "车辆识别号",
          "example": "LSV2A8CA7JN508198"
        },
        "lnr": {
          "type": "string",
          "description": "流水号",
          "example": "0001"
        },
        "assembly_line": {
          "type": "string",
          "description": "装配流水线id",
          "example": "01"
        },
        "prs": {
          "items": {
            "$ref": "#/definitions/PR"
          },
          "type": "array"
        },
        "year": {
          "type": "number",
          "description": "订单年份",
          "example": 2018
        },
        "pin_check_code": {
          "type": "number",
          "description": "pin校验位",
          "example": 5
        },
        "date_planned_start": {
          "format": "date-time",
          "type": "string",
          "description": "生产订单日期",
          "example": "2018-05-19T16:39:57+08:00"
        },
        "equipment_name": {
          "type": "string",
          "description": "设备名",
          "example": "SR1J"
        },
        "factory_name": {
          "type": "string",
          "description": "订单工厂代号",
          "example": "C6"
        },
        "model": {
          "type": "string",
          "description": "车型代码",
          "example": "BR24J3"
        }
      }
    },
    "User": {
      "xml": {
        "name": "User"
      },
      "type": "object",
      "properties": {
        "status": {
          "enum": [
            "active",
            "archived"
          ],
          "type": "string",
          "description": "User Status",
          "example": "active"
        },
        "uuid": {
          "type": "string",
          "example": "112233"
        },
        "image_small": {
          "type": "string"
        },
        "login": {
          "type": "string",
          "example": "gubin@empower.cn"
        },
        "id": {
          "type": "integer",
          "example": 1,
          "format": "int64"
        },
        "name": {
          "type": "string",
          "example": "顾斌"
        }
      }
    },
    "ResponseBody": {
      "type": "object",
      "properties": {
        "jsonrpc": {
          "type": "string",
          "description": "jsonrpc版本",
          "example": "2.0"
        },
        "id": {
          "type": "number",
          "example": 1
        },
        "result": {
          "type": "object",
          "description": "返回结果"
        }
      }
    },
    "Production": {
      "type": "object",
      "properties": {
        "product_id": {
          "type": "integer",
          "example": 1
        },
        "result_ids": {
          "items": {
            "type": "integer",
            "example": 1
          },
          "type": "array"
        },
        "assembly_line_id": {
          "type": "integer",
          "example": 1
        },
        "vin": {
          "type": "string",
          "example": "456464"
        },
        "knr": {
          "type": "string",
          "example": "234242423424"
        },
        "id": {
          "type": "integer",
          "example": 1
        }
      }
    },
    "ResultDetail": {
      "type": "object",
      "properties": {
        "product_id": {
          "type": "integer",
          "example": 1
        },
        "consu_product_id": {
          "type": "integer",
          "example": 1
        },
        "workorder_id": {
          "type": "integer",
          "example": 1
        },
        "workcenter_id": {
          "type": "integer",
          "example": 1
        },
        "measure_result": {
          "enum": [
            "nok",
            "ok",
            "none"
          ],
          "type": "string",
          "example": "none"
        },
        "op_time": {
          "type": "integer",
          "example": 1
        },
        "id": {
          "type": "integer",
          "example": 1
        }
      }
    },
    "ts002Component": {
      "required": [
        "code",
        "is_key"
      ],
      "type": "object",
      "description": "作业所需部件信息",
      "properties": {
        "is_key": {
          "type": "boolean",
          "description": "是否为关键部件"
        },
        "code": {
          "type": "string",
          "description": "部件代码",
          "example": "1111"
        }
      }
    },
    "OperationDetail": {
      "type": "object",
      "properties": {
        "points": {
          "items": {
            "$ref": "#/definitions/point"
          },
          "type": "array"
        },
        "id": {
          "type": "integer"
        },
        "img": {
          "type": "string",
          "description": "作业图片信息(base64)"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "common_resp": {
      "type": "object",
      "properties": {
        "msg": {
          "type": "string"
        }
      }
    },
    "HMILocationInfo": {
      "type": "object",
      "properties": {
        "worksegment": {
          "type": "string",
          "description": "工段名称"
        },
        "workcenter": {
          "type": "string",
          "description": "工位名称"
        }
      }
    },
    "ts002Worksheet": {
      "required": [
        "name",
        "revision",
        "url"
      ],
      "type": "object",
      "description": "工艺指导书",
      "properties": {
        "url": {
          "type": "string",
          "description": "工艺指导书存放的远程路径",
          "example": "http://baidu.com"
        },
        "name": {
          "type": "string",
          "description": "工艺指导书名称",
          "example": "工艺指导书1"
        },
        "revision": {
          "type": "string",
          "description": "工艺版本",
          "example": "v1"
        }
      }
    },
    "PR": {
      "type": "object",
      "properties": {
        "pr_value": {
          "type": "string",
          "description": "pr设定值",
          "example": "G0C"
        },
        "pr_group": {
          "type": "string",
          "description": "pr群组",
          "example": "GSP"
        }
      }
    },
    "operation_edit": {
      "type": "object",
      "properties": {
        "points": {
          "items": {
            "$ref": "#/definitions/point_edit"
          },
          "type": "array"
        },
        "img": {
          "type": "string",
          "description": "作业图片信息(base64)"
        }
      }
    },
    "Consume": {
      "type": "object",
      "properties": {
        "controller_sn": {
          "type": "string",
          "description": "控制器序列号",
          "example": "0001"
        },
        "tolerance_min_degree": {
          "type": "number",
          "description": "最小测量角度",
          "example": 170
        },
        "sequence": {
          "type": "integer",
          "description": "序号",
          "example": 1
        },
        "gun_sn": {
          "type": "string",
          "description": "枪的序列号",
          "example": "0002"
        },
        "max_redo_times": {
          "type": "integer",
          "description": "最大重试次数",
          "format": "int32"
        },
        "nut_no": {
          "type": "string",
          "description": "螺栓编号",
          "example": "1234"
        },
        "pset": {
          "type": "integer",
          "description": "程序号",
          "example": 1
        },
        "result_ids": {
          "items": {
            "type": "integer"
          },
          "type": "array"
        },
        "tolerance_min": {
          "type": "number",
          "description": "最小测量扭矩",
          "example": 1
        },
        "tolerance_max": {
          "type": "number",
          "description": "最大测量扭矩",
          "example": 2
        },
        "tolerance_max_degree": {
          "type": "number",
          "description": "最大测量角度",
          "example": 190
        },
        "offset_x": {
          "type": "integer",
          "description": "左偏移"
        },
        "offset_y": {
          "type": "integer",
          "description": "上偏移"
        },
        "group_sequence": {
          "type": "integer",
          "description": "组序号"
        }
      }
    },
    "ts002Order": {
      "required": [
        "code",
        "track_code",
        "product_code",
        "workcenter",
        "worksheet",
        "product",
        "operation",
        "steps"
      ],
      "type": "object",
      "properties": {
        "product": {
          "$ref": "#/definitions/ts002FinishedProduct"
        },
        "code": {
          "type": "string",
          "description": "生产任务单的唯一标识(Work Order)",
          "example": "MO0001"
        },
        "track_code": {
          "type": "string",
          "description": "车架号",
          "example": "XM205Z03"
        },
        "environments": {
          "items": {
            "$ref": "#/definitions/ts002Environment"
          },
          "type": "array"
        },
        "operation": {
          "$ref": "#/definitions/ts002Operation"
        },
        "product_code": {
          "type": "string",
          "description": "产成品 车架类型号",
          "example": "M000001780589"
        },
        "date_planned_complete": {
          "example": "2019-10-17T12:20:30+08:00",
          "type": "string",
          "description": "计划开工时间",
          "format": "date-time"
        },
        "date_planned_start": {
          "example": "2019-10-16T11:20:30+08:00",
          "type": "string",
          "description": "计划开工时间",
          "format": "date-time"
        },
        "worksheet": {
          "$ref": "#/definitions/ts002Worksheet"
        },
        "workcenter": {
          "type": "string",
          "description": "工位号",
          "example": "TA2-26L-01"
        },
        "steps": {
          "items": {
            "$ref": "#/definitions/ts002WorkStep"
          },
          "type": "array"
        },
        "components": {
          "items": {
            "$ref": "#/definitions/ts002Component"
          },
          "type": "array"
        }
      }
    },
    "ts002WorkStep": {
      "required": [
        "code",
        "sequence",
        "test_type",
        "failure_msg"
      ],
      "type": "object",
      "description": "环境",
      "properties": {
        "code": {
          "type": "string",
          "description": "工步代码",
          "example": "1111"
        },
        "tolerance_max": {
          "type": "number",
          "description": "如果是measure类型，此字段为必填"
        },
        "text": {
          "type": "string",
          "description": "如果是text类型，此字段为必填"
        },
        "desc": {
          "type": "string",
          "description": "环境检查的描述"
        },
        "tightening_total": {
          "type": "integer",
          "description": "拧紧点总数,如果是tightening类型，此字段为必填",
          "example": 16
        },
        "target": {
          "type": "number",
          "description": "如果是measure类型，此字段为必填"
        },
        "tolerance_min": {
          "type": "number",
          "description": "如果是measure类型，此字段为必填"
        },
        "test_type": {
          "default": "text",
          "enum": [
            "text",
            "passfail",
            "measure",
            "tightening",
            "register_consumed_materials"
          ],
          "type": "string",
          "example": "text",
          "description": "检查类型,如果是拧紧的话,根据code找相应的作业进行拧紧点的填充"
        },
        "sequence": {
          "type": "integer",
          "example": 1
        },
        "consume_product": {
          "type": "string",
          "description": "此工步消耗的物料,如果是register_consumed_materials类型，此字段为必填"
        },
        "failure_msg": {
          "type": "string",
          "description": "如果失败弹出的信息"
        }
      },
      "examples": [
        {
          "code": "1111",
          "test_type": "text",
          "sequence": 1,
          "text": "测试文本工步",
          "desc": "测试文本工步描述",
          "failure_msg": "测试文本工步失败消息"
        },
        {
          "code": "222",
          "test_type": "tightening",
          "sequence": 2,
          "desc": "测试拧紧工步描述",
          "failure_msg": "测试拧紧工步失败消息",
          "tightening_total": 4
        }
      ]
    },
    "curve": {
      "type": "object",
      "properties": {
        "file": {
          "type": "string",
          "example": "opration.json"
        },
        "op": {
          "type": "integer",
          "example": 1
        }
      }
    },
    "ts002Environment": {
      "required": [
        "code",
        "text",
        "sequence",
        "test_type"
      ],
      "type": "object",
      "description": "环境,此检查类型必须在正式作业的前面完成。",
      "properties": {
        "text": {
          "type": "string"
        },
        "test_type": {
          "default": "text",
          "enum": [
            "text",
            "passfail"
          ],
          "type": "string",
          "example": "text",
          "description": "检查类型"
        },
        "code": {
          "type": "string",
          "description": "环境检查项编码",
          "example": "1111"
        },
        "sequence": {
          "type": "integer",
          "example": 1
        },
        "desc": {
          "type": "string",
          "description": "环境检查的描述"
        }
      }
    },
    "maintenance_request": {
      "type": "object",
      "properties": {
        "serial_no": {
          "type": "string",
          "description": "拧紧枪序列号"
        },
        "type": {
          "enum": [
            "corrective",
            "calibration",
            "preventive"
          ],
          "type": "string",
          "description": "维护类型"
        }
      }
    },
    "EquipConnection": {
      "type": "object",
      "properties": {
        "connection": {
          "type": "string"
        },
        "serial_no": {
          "type": "string"
        }
      }
    },
    "HMI": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer"
        },
        "uuid": {
          "type": "string"
        }
      }
    },
    "try_maintenance_request": {
      "type": "object",
      "properties": {
        "sin_last_service": {
          "type": "number",
          "description": "自从上次服务后拧紧枪拧紧次数"
        },
        "serial_no": {
          "type": "string",
          "description": "拧紧枪序列号"
        },
        "times": {
          "type": "number",
          "description": "拧紧枪拧紧次数"
        }
      }
    },
    "Logo": {
      "type": "object",
      "properties": {
        "logo": {
          "type": "string",
          "description": "作业图片(base64),图片限制大小1024*1024"
        }
      }
    },
    "Operation": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer"
        },
        "name": {
          "type": "string"
        }
      }
    },
    "WorkOrder": {
      "type": "object",
      "properties": {
        "status": {
          "enum": [
            "pending",
            "ready",
            "process",
            "done",
            "cancel"
          ],
          "type": "string",
          "description": "Order Status"
        },
        "pin_check_code": {
          "type": "string",
          "example": 3334
        },
        "pin": {
          "type": "string",
          "example": 12345
        },
        "worksheet": {
          "type": "string",
          "description": "作业图片"
        },
        "max_op_time": {
          "type": "integer",
          "description": "节拍时间",
          "format": "int32"
        },
        "vin": {
          "type": "string"
        },
        "update_time": {
          "type": "string",
          "example": "2018-05-19T16:39:57+08:00",
          "format": "date-time"
        },
        "lnr": {
          "type": "string",
          "description": "流水号",
          "example": "0001"
        },
        "job": {
          "type": "integer",
          "format": "int64"
        },
        "vehicleTypeImg": {
          "type": "string",
          "description": "车型图片"
        },
        "assembly_line": {
          "type": "string",
          "description": "装配流水线id",
          "example": "01"
        },
        "hmi": {
          "$ref": "#/definitions/HMI"
        },
        "long_pin": {
          "type": "string"
        },
        "year": {
          "type": "integer",
          "example": 2018
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "factory_name": {
          "type": "string",
          "description": "订单工厂代号",
          "example": "C6"
        },
        "model": {
          "type": "string",
          "description": "车型代码",
          "example": "SK234"
        },
        "knr": {
          "type": "string"
        },
        "consumes": {
          "items": {
            "$ref": "#/definitions/Consume"
          },
          "type": "array"
        },
        "equipment_name": {
          "type": "string",
          "description": "设备名",
          "example": "SR1J"
        }
      }
    },
    "image": {
      "type": "object",
      "properties": {
        "content": {
          "type": "string",
          "description": "图片内容"
        },
        "points": {
          "items": {
            "$ref": "#/definitions/point"
          },
          "type": "array"
        }
      }
    },
    "result": {
      "type": "object",
      "properties": {
        "pset_m_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定扭矩阈值"
        },
        "pset_m_max": {
          "type": "number",
          "description": "设定最大扭矩",
          "example": 4.34
        },
        "control_date": {
          "example": "2018-05-19T16:39:57+08:00",
          "type": "string",
          "description": "拧紧时间",
          "format": "date-time"
        },
        "pset_w_max": {
          "type": "number",
          "description": "设定最大角度",
          "example": 4.34
        },
        "user_id": {
          "type": "integer",
          "example": 1,
          "description": "当前操作用户"
        },
        "one_time_pass": {
          "type": "string",
          "example": "fail",
          "description": "是否一次成功（pass/fail）"
        },
        "quality_state": {
          "enum": [
            "exception",
            "pass",
            "fail"
          ],
          "type": "string",
          "description": "质量检测结果",
          "example": "pass"
        },
        "pset_strategy": {
          "enum": [
            "AD",
            "AW",
            "ADW",
            "LN",
            "AN",
            "AT"
          ],
          "type": "string",
          "description": "拧紧枪策略",
          "example": "AD"
        },
        "pset_w_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定角度阈值"
        },
        "pset_w_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小角度"
        },
        "cur_objects": {
          "items": {
            "$ref": "#/definitions/curve"
          },
          "type": "array",
          "description": "个次操作波形对象列表"
        },
        "pset_m_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标扭矩"
        },
        "pset_m_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小扭矩"
        },
        "final_pass": {
          "type": "string",
          "example": "pass",
          "description": "是否最终成功（pass/fail）"
        },
        "measure_degree": {
          "type": "number",
          "description": "实际角度",
          "example": 2.44
        },
        "measure_t_don": {
          "type": "number",
          "description": "拧紧过程花费时间",
          "example": 3.22
        },
        "measure_torque": {
          "type": "number",
          "description": "实际扭矩",
          "example": 3.224
        },
        "measure_result": {
          "enum": [
            "ok",
            "nok"
          ],
          "type": "string",
          "description": "测量结果",
          "example": "ok"
        },
        "op_time": {
          "type": "integer",
          "description": "当前操作",
          "example": 1
        },
        "exception_reason": {
          "type": "string",
          "example": "unknown",
          "description": "异常原因"
        },
        "pset_w_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标角度"
        }
      }
    }
  },
  "basePath": "/api/v1",
  "swagger": "2.0",
  "externalDocs": {
    "url": "http://saturn.com",
    "description": "Saturn Project"
  },
  "securityDefinitions": {
    "api_key": {
      "type": "apiKey",
      "name": "api_key",
      "in": "header"
    }
  }
}
