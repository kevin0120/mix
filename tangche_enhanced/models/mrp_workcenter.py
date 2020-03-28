# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.addons.sa_base.models.mrp_worksegment import DELETE_ALL_MASTER_WROKORDERS_API
import logging
import requests as Requests
from requests import ConnectionError, RequestException, exceptions

import json

_logger = logging.getLogger(__name__)

SYNC_ALL_WORKCENTER_EQUIPMENT_TOOLS = '/rush/v1/equipments/sync'  # 当前只同步拧紧工具


class MrpWorkCenterLoc(models.Model):
    _inherit = 'mrp.workcenter.loc'

    equipment_id = fields.Many2one('maintenance.equipment', 'Remote IO Module For Control',
                                   domain=[('category_name', '=', 'IO')])

    io_output = fields.Integer('IO Output For Picking Indicate', default=1)

    io_input = fields.Integer('IO Input For Picking Confirm', default=1)

    @api.multi
    @api.depends('equipment_id', 'io_output', 'io_input')
    def name_get(self):
        res = []
        for location in self:
            name = u"IN{0}/OUT{1}@{2}".format(location.io_input, location.io_output, location.equipment_id.name)
            res.append((location.id, name))
        return res

    # _sql_constraints = [
    #     ('equipment_input_uniq', 'unique (equipment_id, io_input)', 'The Equipment With Input Must be unique!'),
    #     ('equipment_output_uniq', 'unique (equipment_id, io_output)', 'The Equipment With Output Must be unique!')
    # ]


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    @api.model
    def create(self, vals):
        ret = super(MrpWorkCenter, self).create(vals)
        # 自动创建一个工位组
        val = {
            'code': ret.code,
            'name': ret.name,
            'sa_workcenter_ids': [(4, ret.id, None)]
        }
        self.env['mrp.workcenter.group'].sudo().create(val)
        return ret

    @api.multi
    def button_sync_workcenter_equipments(self):
        self.ensure_one()
        center = self
        connects = center.get_workcenter_masterpc_http_connect()
        if not len(connects):
            _logger.error("Can Not Found MasterPC Connect Info For Work Center:{0}".format(center.name))
            return
        connect = connects[0]
        if not connect:
            _logger.error(
                "Sync Operation Fail. Can Not Found Connect Info For WorkCenter: {0}".format(center.name))
            return
        sync_equipment_url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port,
                                                        SYNC_ALL_WORKCENTER_EQUIPMENT_TOOLS)
        equipment_ids = self.equipment_ids.filtered(
            lambda equipment: equipment.category_name in ['tightening_wrench', 'tightening_gun'])
        vals = []
        for equipment in equipment_ids:
            loc = equipment.location_id
            d = {
                'equipment_sn': equipment.serial_no if equipment else False,  # 拧紧工具,
                'location': {
                    'equipment_sn': loc.equipment_id.serial_no if loc.equipment_id else False,  # IO模块
                    'io_output': loc.io_output,
                    'io_input': loc.io_input
                },
                'type': equipment.category_name  # 设备类型
            }
            vals.append(d)
        try:
            ret = Requests.put(sync_equipment_url, data=json.dumps(vals), headers={'Content-Type': 'application/json'},
                               timeout=1)
            if ret.status_code == 200:
                # operation_id.write({'sync_download_time': fields.Datetime.now()})  ### 更新发送结果
                self.env.user.notify_info(u'同步工位工具成功')
        except ConnectionError as e:
            self.env.user.notify_warning(u'同步工位工具失败, 错误原因:{0}'.format(e.message))
        except RequestException as e:
            self.env.user.notify_warning(u'同步工位工具失败, 错误原因:{0}'.format(e.message))

    @api.multi
    def button_sync_operations(self):
        """
        只同步拧紧相关的工步信息到工位控制器，为了后续的方便使用
        """
        self.ensure_one()
        try:
            operation_obj_sudo = self.env['mrp.routing.workcenter'].sudo()
            for center in self:
                connects = center.get_workcenter_masterpc_http_connect()
                if not len(connects):
                    _logger.error("Can Not Found MasterPC Connect Info For Work Center:{0}".format(center.name))
                    continue
                connect = connects[0]
                if not connect:
                    _logger.error(
                        "Sync Operation Fail. Can Not Found Connect Info For WorkCenter: {0}".format(center.name))
                    continue
                delete_all_endpoint = 'http://{0}:{1}{2}'.format(connect.ip, connect.port,
                                                                 DELETE_ALL_MASTER_WROKORDERS_API)
                center._delete_workcenter_all_opertaions(delete_all_endpoint)
                operations = operation_obj_sudo.search([('workcenter_id', '=', center.id)])
                for operation in operations:
                    operation.button_send_mrp_routing_workcenter()
            return True
        except Exception as e:
            _logger.error("button_sync_operations Error", e)
            raise e


@api.multi
def write(self, vals):
    ret = super(MrpWorkCenter, self).write(vals)
    val = {}
    if 'code' in vals:
        val.update({'code': vals.get('code')})
    if 'name' in vals:
        val.update({'name': vals.get('name')})
    self.mapped('sa_workcentergroup_ids').sudo().write(val)
    return ret


@api.multi
def unlink(self):
    ret = super(MrpWorkCenter, self).unlink()
    self.mapped('sa_workcentergroup_ids').sudo().unlink()
    return ret
