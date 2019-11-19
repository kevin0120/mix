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
