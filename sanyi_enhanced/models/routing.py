# -*- coding: utf-8 -*-

from odoo import models, fields, api, _, SUPERUSER_ID

from odoo.exceptions import UserError, ValidationError

import requests as Requests

from requests import ConnectionError, RequestException, exceptions
import logging
import json
import pprint

from odoo.addons.sa_base.models.mrp_bom import MASTER_ROUTING_API

_logger = logging.getLogger(__name__)


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    @api.multi
    def _push_mrp_routing_workcenter(self, url):
        self.ensure_one()
        operation_id = self
        bom_ids = self.env['mrp.bom'].search([('routing_id.sa_operation_ids', 'in', operation_id.ids)])
        if not bom_ids:
            _logger.debug("_push_mrp_routing_workcenter, BOM:{0}".format(pprint.pformat(bom_ids.ids, indent=4)))
            msg = "Can Not Found MRP BOM Within The Operation:{0}".format(operation_id.name)
            _logger.error(msg)
            raise ValidationError(msg)
        tightening_step_ids = operation_id.sa_step_ids.filtered(lambda step: step.test_type == 'tightening')
        if not tightening_step_ids:
            msg = "Can Not Found Tightening Step For Operation:{0}".format(operation_id.name)
            _logger.error(msg)
            raise ValidationError(msg)
        for tightening_step_id in tightening_step_ids:
            _points = []
            operation_point_ids = tightening_step_id.operation_point_ids
            for point in operation_point_ids:
                _points.append({
                    'sequence': point.sequence,
                    'group_sequence': point.group_sequence,
                    'offset_x': point.x_offset,
                    'offset_y': point.y_offset,
                    'max_redo_times': point.max_redo_times,
                    'tool_sn': point.tool_id.serial_no if point.tool_id else "",
                    'controller_sn': '',
                    'norm_torque': point.norm,
                    'norm_angle': point.norm_degree,
                    'pset': point.program_id.code if point.program_id else False,
                    'pset_name': point.program_id.name if point.program_id else "",
                    'consu_product_id': point.product_id.id if point.product_id.id else 0,
                    'nut_no': point.name,  # 螺栓编号为拧紧点上的名称
                })

            for bom_id in bom_ids:
                val = {
                    # "id": operation_id.id,
                    'tightening_step_ref': tightening_step_id.ref or tightening_step_id.name,  # fixme：考虑ref已经为空
                    'tightening_step_name': tightening_step_id.name,
                    "workcenter_id": operation_id.workcenter_id.id,
                    "job": int(operation_id.op_job_id.code) if operation_id.op_job_id else 0,
                    "max_op_time": operation_id.max_op_time,
                    "name": u"[{0}]{1}@{2}/{3}".format(operation_id.name, operation_id.group_id.code,
                                                       operation_id.workcenter_id.name,
                                                       operation_id.routing_id.name),
                    "img": u'data:{0};base64,{1}'.format('image/png',
                                                         tightening_step_id.worksheet_img) if tightening_step_id.worksheet_img else "",
                    "product_id": bom_id.product_id.id if bom_id else 0,
                    "product_type": bom_id.product_id.default_code if bom_id else "",
                    "workcenter_code": operation_id.workcenter_id.code if operation_id.workcenter_id else "",
                    'product_type_image': u'data:{0};base64,{1}'.format('image/png',
                                                                        bom_id.product_id.image_small) if bom_id.product_id.image_small else "",
                    "points": _points
                }
                try:
                    ret = Requests.put(url, data=json.dumps(val), headers={'Content-Type': 'application/json'},
                                       timeout=1)
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
        try:
            for operation in self:
                if not operation.workcenter_ids:
                    continue
                for workcenter_id in operation.workcenter_ids:
                    connects = workcenter_id.get_workcenter_masterpc_http_connect()
                    if not len(connects):
                        _logger.error(
                            "Can Not Found MasterPC Connect Info For Work Center:{0}".format(workcenter_id.name))
                        continue
                    connect = connects[0]
                    url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_ROUTING_API)
                    operation._push_mrp_routing_workcenter(url)
            self.env.user.notify_info(u'同步工艺成功')
            return True
        except Exception as e:
            self.env.user.notify_warning(u'同步工艺失败:{0}'.format(str(e)))
