# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpWorkCenterLoc(models.Model):
    _inherit = 'mrp.workcenter.loc'

    io_output = fields.Integer('IO Output For Picking Indicate')

    io_input = fields.Integer('IO Output For Picking Confirm')


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
