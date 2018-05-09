# -*- coding: utf-8 -*-


DEFAULT_LIMIT = 80


api_data = {
  "info": {
    "termsOfService": "http://centronsys.com",
    "version": "1.0.0",
    "contact": {
      "email": "gubin@centronsys.com"
    },
    "description": "智能装配应用服务器RESTful",
    "title": "智能装配应用服务器RESTful"
  },
  "paths": {
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
        "description": "获取拧紧结果数据",
        "parameters": [
          {
            "name": "date_from",
            "description": "查询起始时间",
            "in": "query",
            "type": "string",
            "format": "date-time",
            "default": "1996-12-19T16:39:57-08:00"
          },
          {
            "name": "date_to",
            "description": "查询终止时间",
            "in": "query",
            "type": "string",
            "format": "date-time",
            "default": "1996-12-29T16:39:57+08:00"
          },
          {
            "name": "limit",
            "in": "query",
            "type": "integer",
            "default": 80
          }
        ],
        "tags": [
          "Result"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "获取结果数据"
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
        "description": "更新一条拧紧结果数据",
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
            "name": "body",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/result"
            }
          }
        ],
        "tags": [
          "Result"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "更新一条结果数据"
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
        "description": "获取一条拧紧结果数据",
        "parameters": [
          {
            "name": "resultId",
            "description": "需要更新的结果的ID",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer"
          }
        ],
        "tags": [
          "Result"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "获取一条结果数据"
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
        "description": "",
        "parameters": [
          {
            "description": "UUID to filter by",
            "items": {
              "type": "string"
            },
            "in": "query",
            "type": "array",
            "name": "uuids"
          },
          {
            "description": "返回结果限定个数",
            "collectionFormat": "multi",
            "default": 80,
            "required": False,
            "in": "query",
            "type": "integer",
            "name": "limit"
          }
        ],
        "tags": [
          "Users"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "查询用户清单"
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
        "description": "获取某一用户信息",
        "parameters": [
          {
            "description": "用户唯一标示(胸卡信息)",
            "required": True,
            "in": "path",
            "type": "string",
            "name": "uuid"
          }
        ],
        "tags": [
          "Users"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "查询用户清单"
      }
    },
    "/res.users/batch_archived": {
      "put": {
        "responses": {
          "200": {
            "description": "用户清单",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/User"
              }
            }
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "description": "批量归档用户",
        "parameters": [
          {
            "description": "用户唯一标示(胸卡信息)",
            "schema": {
              "type": "array",
              "items": {
                "type": "string",
                "example": "112233"
              }
            },
            "required": True,
            "in": "body",
            "name": "body"
          }
        ],
        "tags": [
          "Users"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "批量归档用户"
      }
    },
    "/mrp.productions": {
      "post": {
        "description": "当AIIS收到FIS下发的装配任务，会调用此API将任务同步下发给ODOO.",
        "parameters": [
          {
            "schema": {
              "$ref": "#/definitions/mission"
            },
            "name": "body",
            "in": "body"
          }
        ],
        "tags": [
          "AIIS"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "下发装配任务",
        "consumes": [
          "application/json"
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/ResponseBody"
            }
          },
          "400": {
            "description": "失败",
            "schema": {
              "$ref": "#/definitions/ResponseBody"
            }
          }
        }
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
    "User": {
      "xml": {
        "name": "User"
      },
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "example": "顾斌"
        },
        "status": {
          "type": "string",
          "description": "User Status",
          "enum": [
            "active",
            "archived"
          ],
          "example": "active"
        },
        "id": {
          "type": "integer",
          "format": "int64",
          "example": 1
        },
        "login": {
          "type": "string",
          "example": "gubin@empower.cn"
        },
        "uuid": {
          "type": "string",
          "example": "112233"
        }
      }
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
    "ResultDetail": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "example": 1
        },
        "product_id": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          }
        },
        "consu_product_id": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          }
        },
        "workorder_id": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          }
        },
        "workcenter_id": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/OdooMany2One"
          }
        },
        "measure_result": {
          "type": "string",
          "enum": [
            "nok",
            "ok",
            "none"
          ],
          "example": "none"
        },
        "op_time": {
          "type": "integer",
          "example": 1
        }
      }
    },
    "result": {
      "type": "object",
      "properties": {
        "pset_strategy": {
          "type": "string",
          "enum": [
            "AD",
            "AW",
            "ADW",
            "LN",
            "AN",
            "AT"
          ],
          "description": "拧紧枪策略",
          "example": "AD"
        },
        "pset_m_max": {
          "type": "number",
          "description": "设定最大扭矩",
          "example": 4.34
        },
        "pset_m_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小扭矩"
        },
        "pset_m_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定扭矩阈值"
        },
        "pset_m_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标扭矩"
        },
        "pset_w_max": {
          "type": "number",
          "description": "设定最大角度",
          "example": 4.34
        },
        "pset_w_min": {
          "type": "number",
          "example": 1.34,
          "description": "设定最小角度"
        },
        "pset_w_threshold": {
          "type": "number",
          "example": 1.44,
          "description": "设定角度阈值"
        },
        "pset_w_target": {
          "type": "number",
          "example": 1.44,
          "description": "设定目标角度"
        },
        "measure_result": {
          "type": "string",
          "enum": [
            "ok",
            "nok"
          ],
          "description": "测量结果",
          "example": "ok"
        },
        "control_date": {
          "type": "string",
          "format": "date-time",
          "description": "拧紧时间",
          "example": "1996-12-19T16:39:57+08:00"
        },
        "op_time": {
          "type": "integer",
          "description": "当前操作",
          "example": 1
        },
        "cur_objects": {
          "items": {
            "type": "string",
            "example": "cur.json"
          },
          "type": "array",
          "description": "个次操作波形对象列表"
        },
        "measure_torque": {
          "type": "number",
          "description": "实际扭矩",
          "example": 3.224
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
        "user_id": {
          "type": "integer",
          "example": 1,
          "description": "当前操作用户"
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
        "vin": {
          "type": "string",
          "description": "车辆识别号",
          "example": "LSV2A8CA7JN508198"
        },
        "factory_name": {
          "type": "string",
          "description": "订单工厂代号",
          "example": "C6"
        },
        "prs": {
          "items": {
            "$ref": "#/definitions/PR"
          },
          "type": "array"
        },
        "pin": {
          "type": "number",
          "description": "订单车身pin码",
          "example": 6473537
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
        "equipment_name": {
          "type": "string",
          "description": "设备名",
          "example": "SR1J"
        },
        "lnr": {
          "type": "string",
          "description": "流水号",
          "example": "0001"
        }
      }
    },
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
    "Order": {
      "xml": {
        "name": "Order"
      },
      "type": "object",
      "properties": {
        "status": {
          "enum": [
            "placed",
            "approved",
            "delivered"
          ],
          "type": "string",
          "description": "Order Status"
        },
        "shipDate": {
          "type": "string",
          "format": "date-time"
        },
        "complete": {
          "default": False,
          "type": "boolean"
        },
        "petId": {
          "type": "integer",
          "format": "int64"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "quantity": {
          "type": "integer",
          "format": "int32"
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
      "flow": "implicit",
      "type": "oauth2",
      "authorizationUrl": "http://petstore.swagger.io/oauth/dialog",
      "scopes": {
        "write:pets": "modify pets in your account",
        "read:pets": "read your pets"
      }
    }
  }
}