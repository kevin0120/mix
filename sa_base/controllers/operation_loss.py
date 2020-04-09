# -*- coding: utf-8 -*-

from odoo import http, fields, api, SUPERUSER_ID
import json
from odoo.http import request, Response

DEFAULT_LIMIT = 80


class OperationLoss(http.Controller):
    @http.route('/api/v1/mrp.workcenter.productivity.loss', type='http', methods=['GET', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def _get_production_loss(self, **kw):
        domain = [('manual', '=', True)]
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if 'limit' in kw.keys():
            limit = int(kw['limit'])
        else:
            limit = DEFAULT_LIMIT
        losses = env['mrp.workcenter.productivity.loss'].search(domain, limit=limit)
        vals = []
        for loss in losses:
            vals.append({
                'loss_id': loss.id,
                'name': loss.name,
                'type:': loss.loss_type,
            })
        body = json.dumps(vals)
        headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
        return Response(body, status=200, headers=headers)
