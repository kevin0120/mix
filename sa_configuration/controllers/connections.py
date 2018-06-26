# -*- coding: utf-8 -*-

from odoo import http, fields,api, SUPERUSER_ID
import json
from odoo.http import request,Response


class HMIConnections(http.Controller):
    @http.route('/api/v1/hmi.connections/<string:serial_no>', type='http', methods=['GET', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def _get_connections(self, serial_no=None):
        hmi_id = None
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if serial_no:
            hmi = env['maintenance.equipment'].search([('serial_no', '=', serial_no)],limit=1)
            if hmi:
                hmi_id = hmi.ids[0]
            else:
                body = json.dumps({'msg': "hmi not existed"})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                return Response(body, status=404, headers=headers)
        workercenter_id = env['mrp.workcenter'].search([('hmi_id', '=', hmi_id)])
        if not workercenter_id:
            body = json.dumps({'msg': "Workcenter not found, plz add hmi to One workcenter"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        if not all([workercenter_id.masterpc_id, workercenter_id.controller_ids, workercenter_id.io_id, workercenter_id.rfid_id]):
            body = json.dumps({'msg': "Workcenter Cofiguration is not Complete"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        masterpc_connection = workercenter_id.masterpc_id.connection_ids.filtered(lambda r: r.protocol == 'http')[0] if workercenter_id.masterpc_id.connection_ids else None
        io_connection = workercenter_id.io_id.connection_ids.filtered(lambda r: r.protocol == 'modbustcp')[0] if workercenter_id.io_id.connection_ids else None
        rfid_connection = workercenter_id.rfid_id.connection_ids.filtered(lambda r: r.protocol == 'rawtcp')[0] if workercenter_id.rfid_id.connection_ids else None


        _controllers = []
        for controller in workercenter_id.controller_ids:
            _controllers.append({
                "serial_no": controller.serial_no,
                "connection": False
            })

        val = {
            'info': {
                'workcenter': workercenter_id.name,
                'worksegment': workercenter_id.worksegment_id.name
            },
            'masterpc': {'serial_no': workercenter_id.masterpc_id.serial_no, 'connection':masterpc_connection.name_get()[0][1] if masterpc_connection else False},
            # 'controller': {'serial_no': workercenter_id.controller_id.serial_no, 'connection': False},
            "controllers": _controllers,
            'io': {'serial_no': workercenter_id.io_id.serial_no, 'connection': io_connection.name_get()[0][1] if io_connection else False},
            'rfid': {'serial_no': workercenter_id.rfid_id.serial_no, 'connection': rfid_connection.name_get()[0][1] if rfid_connection else False}
        }
        body = json.dumps(val)
        headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
        return Response(body, status=200, headers=headers)