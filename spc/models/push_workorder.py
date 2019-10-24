# -*- coding: utf-8 -*-

import logging

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.models import AbstractModel
import requests as Requests
from requests import ConnectionError, RequestException
import json
import logging

logger = logging.getLogger(__name__)

MASTER_WROKORDERS_API = '/rush/v1/workorders'
headers = {'Content-Type': 'application/json'}
_logger = logging.getLogger(__name__)

ORDER_LIMIT = 80

ORDER_PRODUCTION_ORDER_BY = 'production_date ASC'


def str_time_to_rfc3339(s_time):
    sp = s_time.split(' ')
    return sp[0] + 'T' + sp[1] + 'Z'


class PushWorkorder(AbstractModel):
    _name = "workorder.push"

    def _post_workorder_to_masterpc(self, url, orders):
        r = list()
        for order in orders:
            _consumes = list()
            for consu in order.consu_work_order_line_ids:

                _consumes.append({
                    "sequence": consu.bom_line_id.operation_point_id.sequence,
                    "group_sequence": consu.bom_line_id.operation_point_id.group_sequence,
                    'max_redo_times': consu.bom_line_id.operation_point_id.max_redo_times,
                    'offset_x': consu.bom_line_id.operation_point_id.x_offset,
                    'offset_y': consu.bom_line_id.operation_point_id.y_offset,
                    "pset": consu.bom_line_id.program_id.code if consu.bom_line_id.program_id.code else "0",
                    "nut_no": consu.product_id.default_code,
                    "gun_sn": consu.bom_line_id.gun_id.serial_no if consu.bom_line_id.gun_id.serial_no else "",
                    "controller_sn": consu.bom_line_id.controller_id.serial_no if consu.bom_line_id.controller_id.serial_no else "",
                    # 'tolerance_min': _qcps.tolerance_min if _qcps else 0.0,
                    # 'tolerance_max': _qcps.tolerance_max if _qcps else 0.0,
                    # 'tolerance_min_degree': _qcps.tolerance_min_degree if _qcps else 0.0,
                    # 'tolerance_max_degree': _qcps.tolerance_max_degree if _qcps else 0.0,
                })

            vals = {
                'order': {'name': order.name, 'origin': order.origin},
                'workcenter': {'name': order.workcenter_id.name, 'code': order.workcenter_id.code} if order.workcenter_id else None,
                'img_op_code': order.operation_id.code,
                'max_op_time': order.operation_id.max_op_time,
                'track_no': order.track_no,
                'status': order.state,  # pending, ready, process, done, cancel
                'assembly_line': order.production_id.assembly_line_id.code,
                'consumes': _consumes,
                'finished_product': order.product_id.default_code,
                'update_time': str_time_to_rfc3339(order.production_date),
                'job': order.operation_id.op_job_id.code if order.operation_id.op_job_id else "0"
            }
            r.append(vals)
        try:
            logger.debug("try to push workorder to masterpc:{0}".format(url))
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
        limit = self.env['ir.config_parameter'].sudo().get_param('sa.wo.push.num', default=80)
        orders = self.env['mrp.workorder'].sudo().search(domain, limit=int(limit), order=ORDER_PRODUCTION_ORDER_BY)
        masterpcs = orders.mapped('workcenter_id.equipment_ids').filtered(
            lambda equip: equip.category_name == 'MasterPC')
        for master in masterpcs:
            need_send_orders = orders.filtered(lambda r: r.workcenter_id.id == master.workcenter_id.id)
            if not need_send_orders:
                continue
            connections = master.connection_ids.filtered(lambda r: r.protocol == 'http')
            if not connections:
                continue
            url = \
                ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API) for connect in
                 connections][0]
            ret = self._post_workorder_to_masterpc(url, need_send_orders)
        return True
