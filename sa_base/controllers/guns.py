# -*- coding: utf-8 -*-

from odoo import http, fields, api, SUPERUSER_ID
import json
from odoo.http import request, Response


class tools(http.Controller):
    @http.route('/api/v1/tools', type='http', methods=['GET', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def _get_guns(self, **kw):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        domain = [('category_name', 'in', ('tightening_gun', 'tightening_wrench'))]
        if 'serial' in kw.keys():
            domain += [('serial_no', '=', kw['serial'])]

        guns = env['maintenance.equipment'].search(domain)
        gun_values = []
        for g in guns:
            gun_values.append({
                "id": g.id,
                "serial": g.serial_no if g.serial_no else ""
            })

        body = json.dumps(gun_values)
        headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
        return Response(body, status=200, headers=headers)
