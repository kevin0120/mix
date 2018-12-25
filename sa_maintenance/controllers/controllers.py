# -*- coding: utf-8 -*-
import werkzeug

from odoo.api import Environment
import odoo.http as http

from odoo.http import request
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

