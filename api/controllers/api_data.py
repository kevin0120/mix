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
    "/operation.results": {
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
    "/operation.results/{resultId}": {
      "put": {
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
            "required": True,
            "type": "string",
            "description": "MasterPC UUID",
            "name": "masterpc",
            "in": "query"
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
    }
  },
  "host": "172.17.0.1:8069",
  "schemes": [
    "http"
  ],
  "tags": [
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
      "description": "作业相关"
    },
    {
      "name": "Users",
      "description": "用户相关"
    }
  ],
  "definitions": {
    "Category": {
      "xml": {
        "name": "Category"
      },
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
          "type": "string"
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
    "Tag": {
      "xml": {
        "name": "Tag"
      },
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "name": {
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
        "nut_total": {
          "type": "integer",
          "format": "int32"
        },
        "pset": {
          "type": "integer",
          "format": "int64"
        },
        "hmi": {
          "$ref": "#/definitions/HMI"
        },
        "result_ids": {
          "items": {
            "$ref": "#/definitions/result_id"
          },
          "type": "array"
        },
        "worksheet": {
          "type": "string",
          "description": "作业图片"
        },
        "knr": {
          "type": "string"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "vin": {
          "type": "string"
        }
      }
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
    "mission": {
      "type": "object",
      "properties": {
        "pin_check_code": {
          "type": "number",
          "description": "pin校验位",
          "example": 5
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
        "pin": {
          "type": "number",
          "description": "订单车身pin码",
          "example": 6473537
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
        "assembly_line": {
          "type": "string",
          "description": "装配流水线id",
          "example": "01"
        },
        "model": {
          "type": "string",
          "description": "车型代码",
          "example": "BR24J3"
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
        }
      }
    },
    "ApiResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        },
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "type": {
          "type": "string"
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
        },
        "uuid": {
          "type": "string",
          "example": "112233"
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
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
        },
        "result_ids": {
          "items": {
            "type": "integer",
            "example": 1
          },
          "type": "array"
        },
        "assembly_line_id": {
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
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
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
        },
        "consu_product_id": {
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
        },
        "workorder_id": {
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
        },
        "workcenter_id": {
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          },
          "type": "array"
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
          "description": "拧紧时间",
          "type": "string",
          "example": "2018-05-19T16:39:57+08:00",
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
    "OdooMany2One": {
      "allOf": [
        {
          "type": "integer",
          "example": 1
        },
        {
          "type": "string",
          "example": "ceshi"
        }
      ]
    },
    "result_id": {
      "type": "integer"
    },
    "Result": {
      "xml": {
        "name": "Pet"
      },
      "required": [
        "name",
        "photoUrls"
      ],
      "type": "object",
      "properties": {
        "category": {
          "$ref": "#/definitions/Category"
        },
        "status": {
          "enum": [
            "available",
            "pending",
            "sold"
          ],
          "type": "string",
          "description": "pet status in the store"
        },
        "name": {
          "type": "string",
          "example": "doggie"
        },
        "tags": {
          "xml": {
            "wrapped": True,
            "name": "tag"
          },
          "items": {
            "$ref": "#/definitions/Tag"
          },
          "type": "array"
        },
        "photoUrls": {
          "xml": {
            "wrapped": True,
            "name": "photoUrl"
          },
          "items": {
            "type": "string"
          },
          "type": "array"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        }
      }
    }
  },
  "basePath": "/api/v1",
  "swagger": "2.0",
  "externalDocs": {
    "url": "http://swagger.io",
    "description": "Find out more about Swagger"
  },
  "securityDefinitions": {
    "api_key": {
      "type": "apiKey",
      "name": "api_key",
      "in": "header"
    },
    "petstore_auth": {
      "type": "oauth2",
      "flow": "implicit",
      "authorizationUrl": "http://petstore.swagger.io/oauth/dialog",
      "scopes": {
        "write:pets": "modify pets in your account",
        "read:pets": "read your pets"
      }
    }
  }
}