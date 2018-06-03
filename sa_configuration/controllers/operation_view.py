# -*- coding: utf-8 -*-

from odoo import http, fields,api, SUPERUSER_ID
import json
from odoo.http import request,Response


class OperationView(http.Controller):
    @http.route('/api/v1/mrp.routing.workcenter/<int:operation_id>/edit', type='http', methods=['PUT', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def _edit_points(self, operation_id=None):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        operation = env['mrp.routing.workcenter'].search([('id', '=', operation_id)],limit=1)
        if not operation:
            body = json.dumps({'msg': "Operation %d not existed" % operation_id})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        else:
            req_vals = request.jsonrequest
            points = req_vals['points'] if 'points' in req_vals else None
            img = req_vals['img'] if 'img' in req_vals else None
            if img:
                ret = operation.write({'worksheet_img': img})
                if not ret:
                    body = json.dumps({'msg': "Operation %d upload image fail" % operation_id})
                    headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                    return Response(body, status=405, headers=headers)
            if not points:
                body = json.dumps({'msg': "Edit point success"})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=200, headers=headers)
            if not isinstance(points, list):
                body = json.dumps({'msg': "Body must be point array"})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=405, headers=headers)
            for val in points:
                point_id = env['point.point'].search([('res_id', '=', operation_id),
                                                      ('res_model', '=', 'mrp.routing.workcenter'),
                                                      ('sequence', '=', val.sequence)])
                if not point_id:
                    _val = val.update({'res_id': operation_id, 'res_model': 'mrp.routing.workcenter', 'res_field': 'worksheet_img'})
                    env['point.point'].create(_val)
                else:
                    env['point.point'].write(val)
                body = json.dumps({'msg': "Edit point success"})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=200, headers=headers)

    @http.route(['/api/v1/mrp.routing.workcenter/<int:operation_id>', '/api/v1/mrp.routing.workcenter'], type='http', methods=['PUT', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def _get_operations(self, operation_id=None):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if operation_id:
            operation = env['mrp.routing.workcenter'].search([('id', '=', operation_id)], limit=1)
            if not operation:
                body = json.dumps({'msg': "Operation %d not existed" % operation_id})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=404, headers=headers)
            else:
                val = {
                    "id": operation_id,
                    "name": operation.name,
                    "img": operation.worksheet_img
                }
                body = json.dumps(val)
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=200, headers=headers)
        else:
            ### 获取作业清单
            operations = env['mrp.routing.workcenter'].search([])
            vals = []
            for operation in operations:
                vals.append({
                    'id': operation.id,
                    'name': operation.name
                })
            body = json.dumps(vals)
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=200, headers=headers)