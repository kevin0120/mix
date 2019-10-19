# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from dateutil import relativedelta
import datetime
import json
import urllib

import requests as Requests

from requests import ConnectionError, RequestException, exceptions

from odoo import api, fields, models, _
from odoo.exceptions import ValidationError


class MrpWorkAssembly(models.Model):
    _name = 'mrp.assemblyline'
    _description = 'Work Assembly Line'
    _order = "id"

    name = fields.Char('Assembly Line', copy=False)
    code = fields.Char('Reference', copy=False, required=True)
    worksegment_count = fields.Integer('Work Sections', compute='_compute_worksegments_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    worksegment_ids = fields.One2many('mrp.worksection', 'workassembly_id', 'Work Sections', copy=False)

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Assembly Line is allowed')]

    @api.multi
    @api.depends('worksegment_ids')
    def _compute_worksegments_count(self):
        for line in self:
            line.worksegment_count = len(line.worksegment_ids)

    @api.multi
    @api.depends('name', 'code')
    def name_get(self):
        res = []
        for line in self:
            name = u"[{0}] {1}".format(line.code, line.name)
            res.append((line.id, name))
        return res


class MrpWorkSegment(models.Model):
    _name = 'mrp.worksection'
    _description = 'Work Section'
    _order = "id"

    name = fields.Char('Work Section Name', copy=False)
    code = fields.Char('Work Section Reference', copy=False, required=True)
    workassembly_id = fields.Many2one('mrp.assemblyline', string='Work Assembly Line')
    workcenter_count = fields.Integer('Work Centers', compute='_compute_workcenters_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    workcenter_ids = fields.One2many('mrp.workcenter', 'worksegment_id', 'Work Centers', copy=False)

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Section is allowed')]

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

    type = fields.Selection([('operate', 'Operate'),
                             ('rework', 'Rework')], default='operate')

    qc_workcenter_id = fields.Many2one('mrp.workcenter', string='Quality Check Work Center')

    worksegment_id = fields.Many2one('mrp.worksection', copy=False)

    sa_workcentergroup_ids = fields.Many2many('mrp.workcenter', 'mrp_workcenter_rel', 'workcenter_id', 'group_id',
                                              string="Workcenters Group", copy=False)

    # @api.multi
    # def _update_create_workcenter_group_tool_by_tool(self):
    #     tool_category_ids = self.env.ref('sa_base.equipment_Gun') + self.env.ref('sa_base.equipment_Wrench')
    #     for workcenter_id in self:
    #         already_equipment_group_ids = self.env['mrp.workcenter.group.tool'].search(
    #             [('workcenter_id', '=', workcenter_id.id)])
    #         if not already_equipment_group_ids:
    #             continue
    #         already_equipment_ids = already_equipment_group_ids.mapped('tool_id')
    #         tool_ids = workcenter_id.equipment_ids.filtered(
    #             lambda r: r.category_id in tool_category_ids and r not in already_equipment_ids)
    #         for tool in tool_ids:
    #             for wg in workcenter_id.sa_workcentergroup_ids:
    #                 val = {
    #                     "workgroup_id": wg.id,
    #                     "workcenter_id": tool.workcenter_id.id,
    #                     "tool_id": tool.id,
    #                 }
    #                 self.env['mrp.workcenter.group.tool'].sudo().create(val)

    @api.multi
    def write(self, vals):
        ret = super(MrpWorkCenter, self).write(vals)
        return ret
