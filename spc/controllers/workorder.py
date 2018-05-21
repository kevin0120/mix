# -*- coding: utf-8 -*-
from odoo import http, fields,api, SUPERUSER_ID
import json
from odoo.http import request,Response

DEFAULT_LIMIT = 80

DEFAULT_ORDER_BY = 'production_date DESC'


class ApiMrpWorkorder(http.Controller):
    @http.route(['/api/v1/mrp.workorders/<string:order_id>', '/api/v1/mrp.workorders'], type='http', methods=['GET'], auth='none', cors='*', csrf=False)
    def _get_workorders(self, order_id=None, **kw):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if order_id:
            order = env['mrp.workorder'].search([('id', '=', order_id)])[0]
            if not order:
                body = json.dumps({'msg': 'Can not found workorder'})
                return Response(body, headers=[('Content-Type', 'application/json'), ('Content-Length', len(body))],
                                status=404)
            ret = {
                'id': order.id,
                'hmi': {'id': order.workcenter_id.hmi_id.id, 'uuid': order.workcenter_id.hmi_id.serial_no},
                'worksheet': order.worksheet_img,
                'pset': order.operation_id.program_id.code,
                'nut_total': order.consu_product_qty,
                'vin': order.production_id.vin,
                'knr': order.production_id.knr,
                'result_ids': order.result_ids.ids,
                'status': order.state  # pending, ready, process, done, cancel
            }
            body = json.dumps(ret)
            return Response(body, headers=[('Content-Type', 'application/json'), ('Content-Length', len(body))],
                            status=200)
        domain =[]
        if 'masterpc' in kw:
            masterpc_uuid = kw['masterpc']
            workcenter_id = env['mrp.workcenter'].search([('masterpc_id.serial_no', '=', masterpc_uuid)], limit=1)
            if not workcenter_id:
                body = json.dumps({'msg':'Can not found Workcenter'})
                return Response(body, headers=[('Content-Type', 'application/json'), ('Content-Length', len(body))], status=405)
            domain += [('workcenter_id', 'in', workcenter_id.ids)]  # 添加查询域
        if 'limit' in kw.keys():
            limit = int(kw['limit'])
        else:
            limit = DEFAULT_LIMIT
        if 'order' in kw.keys():
            order_by = kw['order'] + ' DESC'
        else:
            order_by = DEFAULT_ORDER_BY

        workorder_ids = env['mrp.workorder'].search(domain, limit=limit, order=order_by)
        _ret = list()
        for order in workorder_ids:
            _ret.append({
                'id': order.id,
                'hmi': {'id':workcenter_id.hmi_id.id, 'uuid': workcenter_id.hmi_id.serial_no},
                'worksheet': order.worksheet_img,
                'pset': order.operation_id.program_id.code,
                'nut_total': order.consu_product_qty,
                'vin': order.production_id.vin,
                'knr': order.production_id.knr,
                'result_ids': order.result_ids.ids,
                'status': order.state  # pending, ready, process, done, cancel
            })
        if len(_ret) == 0:
            body = json.dumps([])
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        body = json.dumps(_ret)
        return Response(body, headers=[('Content-Type', 'application/json'), ('Content-Length', len(body))], status=200)

