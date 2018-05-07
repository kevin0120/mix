# -*- coding: utf-8 -*-
from odoo import http,fields
import json
from odoo.http import request,Response
from odoo.tools import DEFAULT_SERVER_DATETIME_FORMAT
from dateutil import parser

DEFAULT_LIMIT = 80

NORMAL_RESULT_FIELDS_READ = ['workorder_id', 'id', 'product_id', 'consu_product_id', 'op_time', 'measure_result', 'workcenter_id']


class SPC(http.Controller):
    @http.route('/api/v1/operation.results/<int:result_id>', type='json', auth='none', cors='*', csrf=False)
    def _update_results(self, result_id):
        operation_result_id = request.env['operation.result'].sudo().browse(result_id)
        if not operation_result_id:
            body = {'msg': "result %d not existed"% result_id}
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=404, headers=headers)
            return response
        else:
            vals = request.jsonrequest
            if 'cur_objects' in vals:
                vals.update({
                    'cur_objects': json.dumps(vals['cur_objects'])
                })

            if 'control_date' in vals:
                _t = parser.parse(vals['control_date'])
                vals.update({
                    'control_date': fields.Datetime.to_string((_t - _t.utcoffset()))
                })
            ret = operation_result_id.sudo().write(vals)
            if ret:
                body = json.dumps(operation_result_id.read(fields=NORMAL_RESULT_FIELDS_READ)[0])
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                response = Response(body, status=200, headers=headers)
                return response

    @http.route(['/api/v1/operation.results', '/api/v1/operation.results/<int:result_id>'], type='http', auth='none', methods=['get'], cors='*', csrf=False)
    def _get_result_lists(self, result_id, **kw):
        domain = []
        if result_id:
            quality_checks = request.env['operation.results'].sudo().browse(result_id)
        else:
            if 'date_from' in kw.keys():
                domain += [('control_date', '>=', kw['date_from'])]
            if 'date_to' in kw.keys():
                domain += [('control_date', '<=', kw['date_to'])]
            if 'limit' in kw.keys():
                limit = int(kw['limit'])
            else:
                limit = DEFAULT_LIMIT
            quality_checks = request.env['operation.results'].sudo().search(domain, limit=limit)
        ret = quality_checks.fields_get()
        body = json.dumps(ret)
        return Response(body, headers=[('Content-Type', 'application/json'),('Content-Length', len(body))], status=200)
