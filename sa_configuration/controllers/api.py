# -*- coding: utf-8 -*-
from odoo import http
import json
from odoo.http import request,Response

from api_data import api_data,DEFAULT_LIMIT
from result import Result,ResultList


class SaConfiguration(http.Controller):
    # 如果想使用此API,必须在配置文件中指定数据库方可使用
    @http.route('/api/v1/doc', type='http', auth='none', cors='*', csrf=False)
    def api_doc(self):
        return json.dumps(api_data)

    @http.route('/api/v1/assemble_missions', type='json', auth='none', cors='*', csrf=False)
    def assemble_mo_create(self):
        vals = request.jsonrequest
        vin = vals['vin']
        count = request.env['mrp.production'].sudo().search_count(
            [('vin', '=',vin)])
        if count > 0:
            return {'error': 400}

        vechile_code = vals['model']
        vals.pop('model')
        product_id = request.env['product.product'].sudo().search(
            [('vehicle_type_code', 'ilike',vechile_code)], limit=1)[0]
        if not product_id:
            return {'error': 405}
        vals.update({'name': u'{0}--V001--{1}-{2}-{3}={4}'.format(
            vals['equipment_name'],vals['factory_name'],vals['year'],vals['pin'],vals['pin_check_code'])})
        vals.update({'product_id': product_id.id,
                     'bom_id': product_id.active_bom_id.id,
                     'product_tmpl_id': product_id.product_tmpl_id.id,
                     'product_uom_id': product_id.active_bom_id.product_uom_id.id,
                     'routing_id': product_id.active_bom_id.product_uom_id.id})

        vals.pop('prs')
        production = request.env['mrp.production'].sudo().create(vals)
        production.button_plan()  ### 模拟点击安排,自动生成工单

        if production:
            return {'id': production.id}

    @http.route('/api/v1/results', type='http', auth='none', methods=['get'], cors='*', csrf=False)
    def _get_result_lists(self, **kw):
        ret = ResultList()
        fields = ['name']
        domain = []
        if 'date_from' in kw.keys():
            domain += [('','>=', kw['date_from'])]
        if 'date_to' in kw.keys():
          domain += [('', '<=', kw['date_to'])]
        if 'limit' in kw.keys():
            limit = int(kw['limit'])
        else:
            limit = DEFAULT_LIMIT
        quality_checks = request.env['quality.check'].sudo().search_read(domain, fields=fields, limit=limit)
        map(lambda check: ret.append_result(Result.convert_from_quailty_check(check)), quality_checks)
        body = ret.to_json()
        return Response(body, headers=[('Content-Type', 'application/json'),('Content-Length', len(body))])
