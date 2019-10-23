# -*- coding: utf-8 -*-

from odoo import models, fields, api, _, SUPERUSER_ID

from odoo.exceptions import UserError, ValidationError

import requests as Requests

from requests import ConnectionError, RequestException, exceptions
import logging
import json
import pprint

MASTER_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter'

_logger = logging.getLogger(__name__)


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    @api.depends('workcenter_id')
    def _compute_workcenter_group(self):
        for routing in self:
            if routing.workcenter_id and routing.workcenter_id.sa_workcentergroup_ids:
                routing.workcenter_group_id = routing.workcenter_id.sa_workcentergroup_ids[0]

    workcenter_id = fields.Many2one('mrp.workcenter', string='Prefer Work Center', copy=True,
                                    required=True)

    workcenter_group_id = fields.Many2one('mrp.workcenter.group', compute=_compute_workcenter_group, store=True,
                                          readonly=True)

    @api.multi
    def _push_mrp_routing_workcenter(self, url):
        self.ensure_one()
        operation_id = self
        bom_ids = self.env['mrp.bom'].search([('operation_ids', 'in', operation_id.ids), ('active', '=', True)])
        if not bom_ids:
            _logger.info('Cannot Found BOM 4 This Operation:{0}'.format(operation_id.name))
            return
        _points = []
        tightening_step_ids = operation_id.sa_step_ids.filtered(lambda step: step.test_type == 'tightening')
        operation_point_ids = tightening_step_ids.mapped('operation_point_ids')
        for point in operation_point_ids:
            _points.append({
                'sequence': point.sequence,
                'group_sequence': point.group_sequence,
                'offset_x': point.x_offset,
                'offset_y': point.y_offset,
                'max_redo_times': point.max_redo_times,
                'tool_sn': '',  # 默认模式下这里传送的枪的序列号是空字符串
                'controller_sn': '',
                # 'tolerance_min': qcp.tolerance_min,
                # 'tolerance_max': qcp.tolerance_max,
                # 'tolerance_min_degree': qcp.tolerance_min_degree,
                # 'tolerance_max_degree': qcp.tolerance_max_degree,
                'consu_product_id': point.product_id.id if point.product_id.id else 0,
                'nut_no': point.product_id.default_code if point.product_id else '',
            })

        for bom_id in bom_ids:
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
                "product_type": bom_id.product_id.default_code if bom_id else "",
                "workcenter_code": operation_id.workcenter_id.code if operation_id.workcenter_id else "",
                'vehicleTypeImg': u'data:{0};base64,{1}'.format('image/png',
                                                                bom_id.product_id.image_small) if bom_id.product_id.image_small else "",
                "points": _points
            }
            try:
                ret = Requests.put(url, data=json.dumps(val), headers={'Content-Type': 'application/json'}, timeout=1)
                if ret.status_code == 200:
                    # operation_id.write({'sync_download_time': fields.Datetime.now()})  ### 更新发送结果
                    self.env.user.notify_info(u'下发工艺成功')
            except ConnectionError as e:
                self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))
            except RequestException as e:
                self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))

        return True

    @api.multi
    def button_send_mrp_routing_workcenter(self):
        for operation in self:
            if not operation.workcenter_ids:
                continue
            for workcenter_id in operation.workcenter_ids:
                master = workcenter_id.equipment_ids.filtered(
                    lambda equipment: equipment.category_name == 'MasterPC')
                if not master:
                    _logger.warning("Send Operation To Workcenter: {0}, But Cannot Found MasterPC".format(
                        pprint.pformat(workcenter_id.name)))
                    continue
                connections = master.connection_ids.filtered(
                    lambda r: r.protocol == 'http') if master.connection_ids else None
                if not connections:
                    continue
                url = \
                    ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API) for connect in
                     connections][0]

                operation._push_mrp_routing_workcenter(url)
        return True
