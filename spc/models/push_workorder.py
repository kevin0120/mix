# -*- coding: utf-8 -*-

import logging

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.models import AbstractModel
import requests as Requests
from requests import ConnectionError, RequestException
import json

MASTER_WROKORDERS_API = '/rush/v1/workorders'
headers = {'Content-Type': 'application/json'}
_logger = logging.getLogger(__name__)

ORDER_LIMIT = 80

ORDER_ORDER_BY = 'production_date ASC'


def str_time_to_rfc3339(s_time):
    sp = s_time.split(' ')
    return sp[0] + 'T' + sp[1] + 'Z'

class PushWorkorder(AbstractModel):
    _name = "workorder.push"

    def _post_workorder_to_masterpc(self, url, orders):
        r = list()
        for workorder in orders:
            points = self.env['point.point'].sudo().search_read(
                domain=[('res_model', '=', 'mrp.routing.workcenter'), ('res_id', '=', workorder.operation_id.id),
                        ('res_field', '=', 'worksheet_img')],
                fields=['x_offset', 'y_offset'])

            # 工单中的消耗品列表
            _consumes = list()
            for consu in workorder.consu_bom_line_ids:
                # 定位消耗品的qcp
                _qcps = self.env['quality.point'].search_read(
                    domain=[('bom_line_id', '=', consu.bom_line_id.id), ('operation_id', '=', workorder.operation_id.id)],
                    fields=['tolerance_min', 'tolerance_max', 'tolerance_min_degree', 'tolerance_max_degree'])

                _consumes.append({
                    "seq": consu.sequence,
                    "pset": consu.bom_line_id.program_id.code,
                    "nut_no": consu.product_id.screw_type_code,
                    "gun_sn": consu.bom_line_id.gun_id.serial_no,
                    "controller_sn": consu.bom_line_id.controller_id.serial_no,
                    'tolerance_min': _qcps[0]['tolerance_min'],
                    'tolerance_max': _qcps[0]['tolerance_max'],
                    'tolerance_min_degree': _qcps[0]['tolerance_min_degree'],
                    'tolerance_max_degree': _qcps[0]['tolerance_max_degree'],
                    "result_ids": consu.result_ids.ids
                })

            vals = {
                'id': workorder.id,
                'hmi': {'id': workorder.workcenter_id.hmi_id.id, 'uuid': workorder.workcenter_id.hmi_id.serial_no},
                'worksheet': {'content': workorder.worksheet_img, "points": points},
                # 'pset': workorder.operation_id.program_id.code,
                'max_redo_times': workorder.operation_id.max_redo_times,
                'max_op_time': workorder.operation_id.max_op_time,
                # 'nut_total': workorder.consu_product_qty,
                'vin': workorder.production_id.vin,
                'knr': workorder.production_id.knr,
                'long_pin': workorder.production_id.long_pin,
                # 'result_ids': workorder.result_ids.ids,
                'status': workorder.state,  # pending, ready, process, done, cancel
                
                'equipment_name': workorder.production_id.equipment_name,
                'factory_name': workorder.production_id.factory_name,
                'year': workorder.production_id.year,
                'pin': workorder.production_id.pin,
                'pin_check_code': workorder.production_id.pin_check_code,
                'assembly_line': workorder.production_id.assembly_line_id.code,
                'lnr': workorder.production_id.lnr,
                # 'nut_no': workorder.consu_product_id.screw_type_code,
                'consumes': _consumes,
                'update_time': str_time_to_rfc3339(workorder.production_date)
            }
            r.append(vals)
        try:
            ret = Requests.post(url, data=json.dumps(r), headers=headers)
            if ret.status_code == 201:
                orders.write({'sent': True})
                return True
        except ConnectionError:
            _logger.debug(u'masterpc:{0} 链接失败'.format(url))
            return False
        return True


    @api.multi
    def workerorder_push(self):
        domain = [('sent', '=', False)]
        orders = self.env['mrp.workorder'].sudo().search(domain, limit=ORDER_LIMIT, order=ORDER_ORDER_BY)
        masterpcs = orders.mapped('workcenter_id.masterpc_id')
        for master in masterpcs:
            need_send_orders = orders.filtered(lambda r: r.workcenter_id.masterpc_id.id == master.id)
            if not need_send_orders:
                continue
            connections = master.connection_ids.filtered(lambda r: r.protocol == 'http')
            if not connections:
                continue
            url = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API) for connect in connections][0]
            ret = self._post_workorder_to_masterpc(url, need_send_orders)
        return True


