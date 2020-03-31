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
    _inherit = "workorder.push"

    @api.multi
    def workerorder_push(self):
        domain = [('sent', '=', False)]
        limit = self.env['ir.config_parameter'].sudo().get_param('sa.wo.push.limit', default=80)
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