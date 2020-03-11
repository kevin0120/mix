# -*- coding: utf-8 -*-
from odoo import http
import json
from odoo.http import request, Response
from odoo.models import api, SUPERUSER_ID


class SanyiEnhancedEquipment(http.Controller):
    @http.route('/api/v1/equipments/<string:sn>/status', type='json', methods=['PUT', 'OPTIONS'], auth='none',
                cors='*', csrf=False)
    def _update_equipment_status(self, sn=None):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        domain = [('serial_no', '=', sn)]

        equipment = env['maintenance.equipment'].search(domain, limit=1)
        if not equipment:
            return Response(json.dumps({'msg': 'Equipment: {0} Is Not Found'.format(sn)}),
                            headers={'content-type': 'application/json'},
                            status=404)

        req_vals = request.jsonrequest
        status = req_vals['status']
        equipment.write({
            'status': status,
        })

        return Response(status=200)
