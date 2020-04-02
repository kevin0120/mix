# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError
from odoo.addons import decimal_precision as dp

import pprint
import logging
import json
import math

MASTER_WROKORDERS_API = '/rush/v1/workorders'

_logger = logging.getLogger(__name__)


class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    @api.multi
    def _create_bulk_cosume_lines(self):
        consume_sudo = self.env['mrp.wo.consu.line'].sudo()
        for order in self:
            step_ids = order.operation_id.sa_step_ids.mapped('step_id').filtered(lambda qcp: qcp and qcp.test_type != 'tightening_point')
            for idx, step in enumerate(step_ids):
                val = {
                    'sequence': idx + 1,
                    'workorder_id': order.id,
                    'point_id': step.id,
                    'production_id': order.production_id.id,
                    'test_type_id': step.test_type_id.id,
                    'product_id': step.product_id.id,
                    'team_id': step.team_id.id
                }
                ret = consume_sudo.create(val)
                for sub_idx, operation_point in enumerate(step.operation_point_ids):
                    val = {
                        'sequence': idx + 1 + sub_idx,
                        'workorder_id': order.id,
                        'test_type_id': operation_point.test_type_id.id,
                        'production_id': order.production_id.id,
                        'point_id': operation_point.qcp_id.id,
                        'operation_point_id': operation_point.id,
                        'product_id': operation_point.product_id.id,
                        'team_id': step.team_id.id,
                        'tool_id': operation_point.tool_id.id if operation_point.tool_id else False,
                        'program_id': step.program_id.id,
                        'parent_consu_order_line_id': ret.id,
                    }
                    consume_sudo.create(val)

    @api.one
    def button_push_workorder(self):
        order = self
        masterpcs = order.mapped('workcenter_id.equipment_ids').filtered(
            lambda equip: equip.category_name == 'MasterPC')
        for master in masterpcs:
            need_send_orders = self
            connections = master.connection_ids.filtered(lambda r: r.protocol == 'http')
            if not connections:
                continue
            connect = connections[0]
            url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API)
            push_workorder = self.env['workorder.push']
            ret = push_workorder._post_workorder_to_masterpc(url, need_send_orders)
            if ret:
                msg = u'推送工单 masterpc:{0} 成功!'.format(url)
                self.env.user.notify_info(msg)
        return True


class MrpProduction(models.Model):
    _inherit = 'mrp.production'

    @api.one
    def button_push_production(self):
        for order in self.workorder_ids:
            order.button_push_workorder()
