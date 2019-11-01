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

    @staticmethod
    def pack_step_payload(steps, type_tightening_point_id):
        payloads = []
        for step in steps.filtered(lambda t: t.test_type_id != type_tightening_point_id):
            val = {
                "sequence": step.sequence,
                'max_redo_times': step.bom_line_id.operation_point_id.max_redo_times,
                'offset_x': step.bom_line_id.operation_point_id.x_offset,
                'offset_y': step.bom_line_id.operation_point_id.y_offset,
                "pset": step.bom_line_id.program_id.code if step.bom_line_id.program_id.code else "0",
                "nut_no": step.product_id.default_code,
                "gun_sn": step.bom_line_id.gun_id.serial_no if step.bom_line_id.gun_id.serial_no else "",
                "controller_sn": step.bom_line_id.controller_id.serial_no if step.bom_line_id.controller_id.serial_no else "",
            }
            for tightening_point_step in step.child_ids:
                operation_point_id = tightening_point_step.operation_point_id
                val = {
                    'sequence': operation_point_id.sequence,
                    'group_sequence': operation_point_id.group_sequence,
                    'bolt_number': operation_point_id.product_id.default_code,
                    'program_id': operation_point_id.program_id.code,
                }
            payloads.append(val)

        return payloads

    @api.model
    def _post_workorder_to_masterpc(self, url, orders):
        r = list()
        type_tightening_point_id = self.env.ref('quality.test_type_tightening_point').id
        for order in orders:
            _steps = self.pack_step_payload(order.consu_work_order_line_ids, type_tightening_point_id)

            vals = {
                'order': {'name': order.name, 'origin': order.origin},
                'workcenter': {'name': order.workcenter_id.name,
                               'code': order.workcenter_id.code} if order.workcenter_id else None,
                'img_op_code': order.operation_id.code,
                'max_op_time': order.operation_id.max_op_time,
                'track_no': order.track_no,
                'status': order.state,  # pending, ready, process, done, cancel
                'assembly_line': order.production_id.assembly_line_id.code,
                'work_steps': _steps,
                'finished_product': order.product_id.default_code,
                'update_time': str_time_to_rfc3339(order.production_date),
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
        pass
