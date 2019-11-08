# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.addons.sa_base.models.mrp_worksegment import DELETE_ALL_MASTER_WROKORDERS_API
import logging

_logger = logging.getLogger(__name__)


class MrpWorkCenterLoc(models.Model):
    _inherit = 'mrp.workcenter.loc'

    equipment_id = fields.Many2one('maintenance.equipment', 'Remote IO Module For Control',
                                   domain=[('category_name', '=', 'IO')])

    io_output = fields.Integer('IO Output For Picking Indicate')

    io_input = fields.Integer('IO Output For Picking Confirm')

    _sql_constraints = [
        ('equipment_input_uniq', 'unique (equipment_id, io_input)', 'The Equipment With Input Must be unique!'),
        ('equipment_output_uniq', 'unique (equipment_id, io_output)', 'The Equipment With Output Must be unique!')
    ]


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
