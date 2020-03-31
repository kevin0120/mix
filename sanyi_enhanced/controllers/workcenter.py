# -*- coding: utf-8 -*-
from odoo import http
import json
from odoo.http import request, Response
from odoo.models import api, SUPERUSER_ID


class SanyiEnhanced(http.Controller):
    @http.route('/api/v1/workcenters/<string:workcenter_code>', type='http', methods=['GET', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def _get_guns(self, workcenter_code, **kw):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        domain = [('code', '=', workcenter_code)]

        workcenter = env['mrp.workcenter'].search(domain, limit=1)
        if not workcenter:
            return Response(json.dumps({'msg': 'WorkCenter: {0} Is Not Found'.format(workcenter_code)}),
                            headers={'content-type': 'application/json'},
                            status=404)
        values = workcenter._get_workcenter_data()
        body = json.dumps(values)
        headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
        return Response(body, status=200, headers=headers)
