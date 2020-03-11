# -*- coding: utf-8 -*-
from odoo import http
import json
from odoo.http import request, Response
from odoo.models import api, SUPERUSER_ID


class SanyiEnhancedUser(http.Controller):
    @http.route('/api/v1/users/<string:uuid>/status', type='json', methods=['PUT', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def _update_user_status(self, uuid=None):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        domain = [('uuid', '=', uuid)]

        user = env['res.users'].search(domain, limit=1)
        if not user:
            return Response(json.dumps({'msg': 'User: {0} Is Not Found'.format(uuid)}),
                            headers={'content-type': 'application/json'},
                            status=404)

        req_vals = request.jsonrequest
        status = req_vals['working_status']
        user.write({
            'working_status': status,
        })

        return Response(status=200)