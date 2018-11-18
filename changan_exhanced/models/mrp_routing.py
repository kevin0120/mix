# -*- coding: utf-8 -*-

from odoo import models, fields, api

import json

import requests as Requests

from requests import ConnectionError, RequestException

MASTER_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter'


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    def _push_mrp_routing_workcenter(self, url, mrp_routing_workcenters):
        for operation_id in mrp_routing_workcenters:
            bom_id = self.env['mrp.bom'].search([('operation_ids', 'in', operation_id.ids)], limit=1)
            if not bom_id:
                continue
            _points = []
            for point in operation_id.operation_point_ids:
                _points.append({
                    'sequence': point.sequence,
                    'group_sequence': point.group_sequence,
                    'offset_x': point.x_offset,
                    'offset_Y': point.y_offset,
                    'max_redo_times': point.max_redo_times
                })

            val = {
                "id": operation_id.id,
                "workcenter_id": operation_id.workcenter_id.id,
                "job": int(operation_id.op_job_id.code) if operation_id.op_job_id else 0,
                "max_op_time": operation_id.max_op_time,
                "name": u"[{0}]{1}@{2}/{3}".format(operation_id.name, operation_id.group_id.code,
                                                   operation_id.workcenter_id.name,
                                                   operation_id.routing_id.name),
                "img": u'data:{0};base64,{1}'.format('image/png',
                                                     operation_id.worksheet_img) if operation_id.worksheet_img else "",
                "product_id": bom_id.product_id.id if bom_id else 0,
                "product_type": bom_id.product_id.vehicle_type_code if bom_id else "",
                "workcenter_code": operation_id.workcenter_id.code if operation_id.workcenter_id else "",
                'vehicleTypeImg': u'data:{0};base64,{1}'.format('image/png',
                                                                bom_id.product_id.image_small) if bom_id.product_id.image_small else "",
                "points": _points
            }
            try:
                ret = Requests.put(url, data=json.dumps(val), headers={'Content-Type': 'application/json'})
                if ret.status_code == 200:
                    operation_id.write({'sync_download_time': fields.Datetime.now()})  ### 更新发送结果
            except ConnectionError:
                break
            except RequestException as e:
                print(e)
        return True

    @api.multi
    def button_send_mrp_routing_workcenter(self):
        self.ensure_one()
        master = self.workcenter_id.masterpc_id if self.workcenter_id else None
        if not master:
            return True
        connections = master.connection_ids.filtered(lambda r: r.protocol == 'http') if master.connection_ids else None
        if not connections:
            return True
        url = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API) for connect in connections][0]

        self._push_mrp_routing_workcenter(url, self)
        return True
