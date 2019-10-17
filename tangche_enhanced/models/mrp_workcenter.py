# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    workcenter_group_ids = fields.Many2many('mrp.workcenter.group', 'mrp_workcenter_rel', 'workcenter_id', 'group_id',
                                     string="Workcenter Group", copy=False)

    @api.multi
    def write(self, vals):
        ret = super(MrpWorkCenter, self).write(vals)
        if 'name' not in vals:
            return ret
        for workcenter in self:
            group_id = workcenter.workcenter_group_ids and workcenter.workcenter_group_ids[0]
            if not group_id:
                continue
            group_id.write(vals.get('name'))
        return ret

    @api.model
    def create(self, vals):
        ret = super(MrpWorkCenter, self).create(vals)
        if ret:
            val = {
                'code': ret.code,
                'name': ret.name,
                'sa_workcenter_ids': [(4, ret.id, None)]
            }
            self.env['mrp.workcenter.group'].sudo().create(val)  # 创建一个工位后，自动创建一个工位组
        return ret