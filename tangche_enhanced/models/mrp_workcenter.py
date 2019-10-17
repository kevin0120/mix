# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    @api.model
    def create(self, vals):
        ret = super(MrpWorkCenter, self).create(vals)
        if ret:
            val = {
                'code': ret.code,
                'name': ret.name
            }
            self.env['mrp.workcenter.group'].sudo().create(val)  # 创建一个工位后，自动创建一个工位组
        return ret