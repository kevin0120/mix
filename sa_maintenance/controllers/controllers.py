# -*- coding: utf-8 -*-
import werkzeug


from odoo import http, fields,api, SUPERUSER_ID
import json
from odoo.http import request, Response

from odoo.api import Environment

from odoo import SUPERUSER_ID
from odoo import registry as registry_get


class SaMaintenance(http.Controller):
    @http.route([
        "/maintenance/requests/<int:ticket_id>",
        "/maintenance/requests/<int:ticket_id>/<token>"
    ], type='http', auth="public")
    def get_maintenance_requests(self, db, ticket_id, token=None):
        registry = registry_get(db)
        with registry.cursor() as cr:
            env = Environment(cr, SUPERUSER_ID, {})
            Ticket = False
            if token:
                Ticket = env['maintenance.request'].sudo().search(
                    [('id', '=', ticket_id), ('access_token', '=', token)])
            else:
                Ticket = env['maintenance.request'].browse(ticket_id)
            if not Ticket:
                return request.not_found()

            return werkzeug.utils.redirect('/web?db=%s#id=%s&view_type=form&model=maintenance.request' % (db, ticket_id))

    @http.route(["/maintenance/requests"], type='json', methods=['POST', 'OPTIONS'], auth="public")
    def _create_maintenance_requests(self, **kw):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if 'serial_no' not in kw or 'description' not in kw:
            body = json.dumps({'msg': "payload must contain serial number!!!!"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        gun_id = env['maintenance.equipment'].search([('serial_no', '=', kw["serial_no"])], limit=1)
        if not gun_id:
            body = json.dumps({'msg': "not found equipment(gun)!!!!"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        vals = {
            "name": "Maintenance:{0}@{1}".format(gun_id.serial_no, fields.Date.context_today()),
            "maintenance_type": "corrective",
            "equipment_id": gun_id.id,
            "description": kw['description']
        }

        ret = env['maintenance.request'].create(vals)
        if not ret:
            body = json.dumps({'msg': "create maintenance request fail!!!!"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=400, headers=headers)
        else:
            body = json.dumps({'msg': "create maintenance request success!!!!"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=201, headers=headers)


