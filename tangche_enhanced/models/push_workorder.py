# -*- coding: utf-8 -*-

import logging

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.addons.tangche_enhanced.controllers.mrp_order_gateway import package_tightening_points, post_order_2_masterpc
from odoo.models import AbstractModel
from tenacity import RetryError
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

    @staticmethod
    def pack_step_payload(steps, type_tightening_id):
        payloads = []
        for step in steps.filtered(lambda t: t.test_type_id.id == type_tightening_id):
            ts = {}
            ts.update({'tightening_image_by_step_code': step.name or step.ref})
            val = package_tightening_points(step.operation_point_ids)
            ts.update({'tightening_points': val})  # 将拧紧点的包包裹进去
            payloads.append(ts)

        return payloads

    @api.model
    def _post_workorder_to_masterpc(self, url, orders):
        r = list()
        type_tightening_id = self.env.ref('quality.test_type_tightening').id
        for order in orders:
            _steps = self.pack_step_payload(order.consu_work_order_line_ids.mapped('point_id'), type_tightening_id)

            vals = {
                'code': order.name,
                'track_no': order.track_no,
                'product_code': order.product_id.default_code,
                'workcenter': order.workcenter_id.code if order.workcenter_id else None,
                'date_planned_start': str_time_to_rfc3339(order.date_planned_start) if order.date_planned_start else str_time_to_rfc3339(fields.Datetime.now()),
                'date_planned_complete': str_time_to_rfc3339(order.date_planned_finished) if order.date_planned_finished else str_time_to_rfc3339(fields.Datetime.now()),
                'worksheet': {
                    "name": "",
                    "revision": "",
                    "url": ""
                },
                'product': {
                    "code": order.product_id.default_code,
                    "url": ""
                },
                'operation': {
                    'code': order.operation_id.name or order.operation_id.ref,
                    'desc': '',
                    "resources": {
                        "users": [],
                        "equipments": [],
                    },
                    'components': [],
                    'environments': [],
                    'steps': _steps,
                },
            }
            r.append(vals)
        try:
            logger.debug("try to push workorder to masterpc:{0}".format(url))
            ret = post_order_2_masterpc(url, r)
            if ret.status_code == 201:
                orders.write({'sent': True})
                return True
        except RetryError:
            msg = u'推送工单 masterpc:{0} 链接失败'.format(url)
            _logger.debug(msg)
            raise UserError(msg)
        return True

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
            connect = connections[0]
            url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API)
            ret = self._post_workorder_to_masterpc(url, need_send_orders)
        return True
