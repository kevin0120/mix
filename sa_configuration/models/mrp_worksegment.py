# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from dateutil import relativedelta
import datetime

from odoo import api, exceptions, fields, models, _


class MrpWorkSegment(models.Model):
    _name = 'mrp.worksegament'
    _description = 'Work Segment'
    _order = "id"

    name = fields.Char('Segament', copy=False)
    code = fields.Char('Reference', copy=False)
    workcenter_count = fields.Integer('Work Centers', compute='_compute_workcenters_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    workcenter_ids = fields.One2many('mrp.workcenter', 'worksegment_id', 'Work Centers', copy=False)

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Segment is allowed')]

    @api.multi
    @api.depends('workcenter_ids')
    def _compute_workcenters_count(self):
        for segment in self:
            segment.workcenter_count = len(segment.workcenter_ids)

    @api.multi
    @api.depends('name', 'code')
    def name_get(self):
        res = []
        for segment in self:
            name = u"[{0}] {1}".format(segment.code, segment.name)
            res.append((segment.id, name))
        return res


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    worksegment_id = fields.Many2one('mrp.worksegament', ondelete='cascade', copy=False)
    hmi_id = fields.Many2one('maintenance.equipment', ondelete='cascade', copy=False)