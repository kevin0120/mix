# -*- coding: utf-8 -*-

from odoo import http, fields, api, SUPERUSER_ID
import json
from werkzeug.exceptions import BadRequest
from odoo.http import request, Response


class TS009Controller(http.Controller):
    @http.route('/api/v1/rush/connections', auth='public')
    def _get_rush_connections(self, **kw):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        masterpcs = env['maintenance.equipment'].search([('category_name', '=', 'MasterPC')], limit=200)
        if not masterpcs:
            return BadRequest('Can Not Found Any Terminal To Connect')
        ret = []
        for masterpc in masterpcs:
            connections = masterpc.connection_ids.filtered(
                lambda r: r.protocol == 'http') if masterpc.connection_ids else None
            if not connections:
                raise None
            connect = connections[0]
            d = {
                "workcenter": masterpc.workcenter_id.display_name if masterpc.workcenter_id else masterpc.serial_no,  # 如果找不动工位名称用序列号代替
                "serial_no": masterpc.serial_no,
                "connection": 'ws://{0}:{1}'.format(connect.ip, connect.port)
            }
            ret.append(d)
        body = json.dumps(ret)
        headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
        return Response(body, status=200, headers=headers)
