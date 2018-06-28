# -*- coding: utf-8 -*-

from odoo import http, fields,api, SUPERUSER_ID
import json
from odoo.http import request,Response


class OperationView(http.Controller):
    @http.route('/api/v1/mrp.routing.workcenter/<int:operation_id>/edit', type='json', methods=['PUT', 'OPTIONS'], auth='none', cors='*', csrf=False)
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

            current_max_seq = len(operation.operation_point_ids)
            for val in points:
                if not val.has_key('x_offset') or not val.has_key('y_offset'):
                    # 忽略没有x_offset或y_offset的点
                    continue

                if val.has_key('sequence'):

                    point = env['operation.point'].search([('operation_id', '=', operation_id),
                                                          ('sequence', '=', val['sequence'])])
                    if not point:
                        # 新增点
                        current_max_seq = current_max_seq + 1
                        val.update({
                            'operation_id': operation_id,
                            'sequence': current_max_seq,
                            'x_offset': val['x_offset'],
                            'y_offset': val['y_offset']
                        })
                        env['operation.point'].create(val)
                    else:
                        # 更新点
                        point.write(val)
                else:
                    # 新增点
                    current_max_seq = current_max_seq + 1
                    val.update({
                        'operation_id': operation_id,
                        'sequence': current_max_seq,
                        'x_offset': val['x_offset'],
                        'y_offset': val['y_offset']
                    })
                    env['operation.point'].create(val)

            body = json.dumps({'msg': "Edit point success"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=200, headers=headers)

    @http.route(['/api/v1/mrp.routing.workcenter/<int:operation_id>', '/api/v1/mrp.routing.workcenter'], type='http', methods=['GET'], auth='none', cors='*', csrf=False)
    def _get_operations(self, operation_id=None):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if operation_id:
            operation = env['mrp.routing.workcenter'].search([('id', '=', operation_id)], limit=1)
            if not operation:
                body = json.dumps({'msg': "Operation %d not existed" % operation_id})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=404, headers=headers)
            else:
                _points = []
                for point in operation.operation_point_ids:
                    _points.append({
                        'seq': point.sequence,
                        'x_offset': point.x_offset,
                        'y_offset': point.y_offset
                    })

                val = {
                    "id": operation_id,
                    "name": u"[{0}]{1}@{2}/{3}".format(operation.name, operation.group_id.code, operation.workcenter_id.name, operation.routing_id.name),
                    "img": operation.worksheet_img,
                    "points": _points
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
                    'name': u"[{0}]{1}@{2}/{3}".format(operation.name, operation.group_id.code, operation.workcenter_id.name, operation.routing_id.name)
                })
            body = json.dumps(vals)
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=200, headers=headers)