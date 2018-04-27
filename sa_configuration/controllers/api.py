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
