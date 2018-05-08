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
    "/operation.results/{resultId}": {
      "put": {
        "responses": {
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
          "result"
        ],
        "produces": [
          "application/json"
        ],
        "consumes": [
          "application/json"
        ],
        "summary": "更新一条结果数据"
      }
    },
    "/results": {
      "put": {
        "responses": {
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Pet not found"
          },
          "405": {
            "description": "Validation exception"
          }
        },
        "description": "",
        "parameters": [
          {
            "schema": {
              "$ref": "#/definitions/Result"
            },
            "description": "Pet object that needs to be added to the store",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "tags": [
          "pet"
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "security": [
          {
            "petstore_auth": [
              "write:pets",
              "read:pets"
            ]
          }
        ],
        "operationId": "updatePet",
        "consumes": [
          "application/json",
          "application/xml"
        ],
        "summary": "Update an existing pet"
      },
      "post": {
        "responses": {
          "201": {
            "description": "成功创建一条结果",
            "schema": {
              "$ref": "#/definitions/Result"
            }
          },
          "405": {
            "description": "Invalid input"
          }
        },
        "description": "",
        "parameters": [
          {
            "schema": {
              "$ref": "#/definitions/Result"
            },
            "description": "添加工序结果",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "tags": [
          "result"
        ],
        "produces": [
          "application/json"
        ],
        "operationId": "addResult",
        "consumes": [
          "application/json"
        ],
        "summary": "添加一条新的(拧紧)工序结果"
      },
      "get": {
        "responses": {
          "200": {
            "description": "结果清单",
            "schema": {
              "items": {
                "$ref": "#/definitions/Result"
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
            "description": "起始时间",
            "required": False,
            "type": "date",
            "name": "date_from",
            "in": "query"
          },
          {
            "description": "截止时间",
            "required": False,
            "type": "date",
            "name": "date_to",
            "in": "query"
          },
          {
            "description": "返回结果限定个数",
            "default": 80,
            "required": False,
            "in": "query",
            "type": "integer",
            "name": "limit"
          }
        ],
        "tags": [
          "result"
        ],
        "produces": [
          "application/json"
        ],
        "operationId": "getResultList",
        "consumes": [
          "application/json"
        ],
        "summary": "查询结果清单"
      }
    },
    "/pet/findByTags": {
      "get": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "items": {
                "$ref": "#/definitions/Pet"
              },
              "type": "array"
            }
          },
          "400": {
            "description": "Invalid tag value"
          }
        },
        "description": "Muliple tags can be provided with comma separated strings. Use         tag1, tag2, tag3 for testing.",
        "parameters": [
          {
            "description": "Tags to filter by",
            "items": {
              "type": "string"
            },
            "required": True,
            "collectionFormat": "multi",
            "in": "query",
            "type": "array",
            "name": "tags"
          }
        ],
        "tags": [
          "pet"
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "deprecated": True,
        "security": [
          {
            "petstore_auth": [
              "write:pets",
              "read:pets"
            ]
          }
        ],
        "operationId": "findPetsByTags",
        "summary": "Finds Pets by tags"
      }
    },
    "/user/createWithArray": {
      "post": {
        "responses": {
          "default": {
            "description": "successful operation"
          }
        },
        "description": "",
        "parameters": [
          {
            "schema": {
              "items": {
                "$ref": "#/definitions/User"
              },
              "type": "array"
            },
            "description": "List of user object",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Creates list of users with given input array",
        "operationId": "createUsersWithArrayInput"
      }
    },
    "/store/order": {
      "post": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "400": {
            "description": "Invalid Order"
          }
        },
        "description": "",
        "parameters": [
          {
            "schema": {
              "$ref": "#/definitions/Order"
            },
            "description": "order placed for purchasing the pet",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "store"
        ],
        "summary": "Place an order for a pet",
        "operationId": "placeOrder"
      }
    },
    "/store/order/{orderId}": {
      "delete": {
        "responses": {
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Order not found"
          }
        },
        "description": "For valid response try integer IDs with positive integer value.         Negative or non-integer values will generate API errors",
        "parameters": [
          {
            "description": "ID of the order that needs to be deleted",
            "format": "int64",
            "required": True,
            "minimum": 1,
            "in": "path",
            "type": "integer",
            "name": "orderId"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "store"
        ],
        "summary": "Delete purchase order by ID",
        "operationId": "deleteOrder"
      },
      "get": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/Order"
            }
          },
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Order not found"
          }
        },
        "description": "For valid response try integer IDs with value >= 1 and <= 10.         Other values will generated exceptions",
        "parameters": [
          {
            "minimum": 1,
            "name": "orderId",
            "in": "path",
            "format": "int64",
            "required": True,
            "type": "integer",
            "maximum": 10,
            "description": "ID of pet that needs to be fetched"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "store"
        ],
        "summary": "Find purchase order by ID",
        "operationId": "getOrderById"
      }
    },
    "/user/login": {
      "get": {
        "responses": {
          "200": {
            "headers": {
              "X-Rate-Limit": {
                "type": "integer",
                "description": "calls per hour allowed by the user",
                "format": "int32"
              },
              "X-Expires-After": {
                "type": "string",
                "description": "date in UTC when token expires",
                "format": "date-time"
              }
            },
            "description": "successful operation",
            "schema": {
              "type": "string"
            }
          },
          "400": {
            "description": "Invalid username/password supplied"
          }
        },
        "description": "",
        "parameters": [
          {
            "description": "The user name for login",
            "required": True,
            "type": "string",
            "name": "username",
            "in": "query"
          },
          {
            "description": "The password for login in clear text",
            "required": True,
            "type": "string",
            "name": "password",
            "in": "query"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Logs user into the system",
        "operationId": "loginUser"
      }
    },
    "/user": {
      "post": {
        "responses": {
          "default": {
            "description": "successful operation"
          }
        },
        "description": "This can only be done by the logged in user.",
        "parameters": [
          {
            "schema": {
              "$ref": "#/definitions/User"
            },
            "description": "Created user object",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Create user",
        "operationId": "createUser"
      }
    },
    "/user/{username}": {
      "put": {
        "responses": {
          "400": {
            "description": "Invalid user supplied"
          },
          "404": {
            "description": "User not found"
          }
        },
        "description": "This can only be done by the logged in user.",
        "parameters": [
          {
            "description": "name that need to be updated",
            "required": True,
            "type": "string",
            "name": "username",
            "in": "path"
          },
          {
            "schema": {
              "$ref": "#/definitions/User"
            },
            "description": "Updated user object",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Updated user",
        "operationId": "updateUser"
      },
      "delete": {
        "responses": {
          "400": {
            "description": "Invalid username supplied"
          },
          "404": {
            "description": "User not found"
          }
        },
        "description": "This can only be done by the logged in user.",
        "parameters": [
          {
            "description": "The name that needs to be deleted",
            "required": True,
            "type": "string",
            "name": "username",
            "in": "path"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Delete user",
        "operationId": "deleteUser"
      },
      "get": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "400": {
            "description": "Invalid username supplied"
          },
          "404": {
            "description": "User not found"
          }
        },
        "description": "",
        "parameters": [
          {
            "description": "The name that needs to be fetched. Use user1 for testing. ",
            "required": True,
            "type": "string",
            "name": "username",
            "in": "path"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Get user by user name",
        "operationId": "getUserByName"
      }
    },
    "/pet/findByStatus": {
      "get": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "items": {
                "$ref": "#/definitions/Pet"
              },
              "type": "array"
            }
          },
          "400": {
            "description": "Invalid status value"
          }
        },
        "description": "Multiple status values can be provided with comma separated strings",
        "parameters": [
          {
            "description": "Status values that need to be considered for filter",
            "items": {
              "default": "available",
              "enum": [
                "available",
                "pending",
                "sold"
              ],
              "type": "string"
            },
            "required": True,
            "collectionFormat": "multi",
            "in": "query",
            "type": "array",
            "name": "status"
          }
        ],
        "tags": [
          "pet"
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "security": [
          {
            "petstore_auth": [
              "write:pets",
              "read:pets"
            ]
          }
        ],
        "operationId": "findPetsByStatus",
        "summary": "Finds Pets by status"
      }
    },
    "/store/inventory": {
      "get": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "additionalProperties": {
                "type": "integer",
                "format": "int32"
              },
              "type": "object"
            }
          }
        },
        "description": "Returns a map of status codes to quantities",
        "parameters": [],
        "tags": [
          "store"
        ],
        "produces": [
          "application/json"
        ],
        "security": [
          {
            "api_key": []
          }
        ],
        "operationId": "getInventory",
        "summary": "Returns pet inventories by status"
      }
    },
    "/user/logout": {
      "get": {
        "responses": {
          "default": {
            "description": "successful operation"
          }
        },
        "description": "",
        "parameters": [],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Logs out current logged in user session",
        "operationId": "logoutUser"
      }
    },
    "/user/createWithList": {
      "post": {
        "responses": {
          "default": {
            "description": "successful operation"
          }
        },
        "description": "",
        "parameters": [
          {
            "schema": {
              "items": {
                "$ref": "#/definitions/User"
              },
              "type": "array"
            },
            "description": "List of user object",
            "required": True,
            "name": "body",
            "in": "body"
          }
        ],
        "produces": [
          "application/xml",
          "application/json"
        ],
        "tags": [
          "user"
        ],
        "summary": "Creates list of users with given input array",
        "operationId": "createUsersWithListInput"
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
    },
    "/pet/{petId}/uploadImage": {
      "post": {
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/ApiResponse"
            }
          }
        },
        "description": "",
        "parameters": [
          {
            "description": "ID of pet to update",
            "format": "int64",
            "required": True,
            "in": "path",
            "type": "integer",
            "name": "petId"
          },
          {
            "description": "Additional data to pass to server",
            "required": False,
            "type": "string",
            "name": "additionalMetadata",
            "in": "formData"
          },
          {
            "description": "file to upload",
            "required": False,
            "type": "file",
            "name": "file",
            "in": "formData"
          }
        ],
        "tags": [
          "pet"
        ],
        "produces": [
          "application/json"
        ],
        "security": [
          {
            "petstore_auth": [
              "write:pets",
              "read:pets"
            ]
          }
        ],
        "operationId": "uploadFile",
        "consumes": [
          "multipart/form-data"
        ],
        "summary": "uploads an image"
      }
    }
  },
  "host": "172.17.0.1:8069",
  "schemes": [
    "http"
  ],
  "tags": [
    {
      "name": "result",
      "description": "工序操作结果"
    },
    {
      "name": "HMI",
      "description": "HMI设备"
    },
    {
      "name": "Manufacture",
      "description": "作业相关"
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
        "username": {
          "type": "string"
        },
        "phone": {
          "type": "string"
        },
        "userStatus": {
          "type": "integer",
          "description": "User Status",
          "format": "int32"
        },
        "firstName": {
          "type": "string"
        },
        "lastName": {
          "type": "string"
        },
        "password": {
          "type": "string"
        },
        "id": {
          "type": "integer",
          "format": "int64"
        },
        "email": {
          "type": "string"
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
          "description": "当前操作"
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