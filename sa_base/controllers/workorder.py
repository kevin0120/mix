# -*- coding: utf-8 -*-

from odoo import http, fields, api, SUPERUSER_ID
import json
from odoo.exceptions import UserError
from odoo.tools import exception_to_unicode
from odoo.http import request, Response
from jsonschema import validators, validate, ValidationError

DEFAULT_LIMIT = 80

pending_req_order_schema = {
    "type": "object",
    "properties": {
        "except_type": {"type": "string"},
        "except_code": {"type": "string"},
        "order_name": {"type": "string"},
        "workcenter_code": {"type": 'string'},
        "end_time": {"type": "string"},  # datetime
    },
    "required": ["except_type"]
}


class WorkOrderController(http.Controller):
    @http.route('/api/v1/workorder/<string:order_name>/resume', type='http', methods=['PUT', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def put_workorder_resume(self, order_name, **kwargs):
        domain = [('name', '=', order_name)]
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)

        order = env['mrp.workorder'].search(domain, limit=1)
        if not order:
            body = json.dumps({'msg': "order: {} Can Not Be Found".format(order_name)})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        return Response(status=200)

    @http.route('/api/v1/workorder/<string:order_name>/pending', type='http', methods=['PUT', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def put_workorder_pending(self, order_name, **kw):
        domain = [('name', '=', order_name)]
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        try:
            order = env['mrp.workorder'].search(domain, limit=1)
            if not order:
                raise ValidationError("order: {} Can Not Be Found".format(order_name))
            req_vals = request.jsonrequest
            validate(req_vals, pending_req_order_schema)

            except_code = req_vals.get('except_code', None)
            end_time = req_vals.get('end_time', None)
            loss_id = self.env['mrp.workcenter.productivity.loss'].search([('name', '=', except_code)], limit=1)
            if not loss_id:
                raise UserError(u'未找到相关的异常代码:{}'.format(except_code))
            body = json.dumps(req_vals)
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=200, headers=headers)
        except Exception as e:
            body = json.dumps({'msg': "Validate Error: {}".format(exception_to_unicode(e))})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
