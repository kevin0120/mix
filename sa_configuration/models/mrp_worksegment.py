# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from dateutil import relativedelta
import datetime
import json

from odoo import api, fields, models, _
from odoo.exceptions import ValidationError


class MrpWorkAssembly(models.Model):
    _name = 'mrp.assemblyline'
    _description = 'Work Assembly Line'
    _order = "id"

    name = fields.Char('Assembly Line', copy=False)
    code = fields.Char('Reference', copy=False)
    worksegment_count = fields.Integer('Work Segments', compute='_compute_worksegments_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    worksegment_ids = fields.One2many('mrp.worksegament', 'workassembly_id', 'Work Segments', copy=False)

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
    _name = 'mrp.worksegament'
    _description = 'Work Segment'
    _order = "id"

    name = fields.Char('Segament', copy=False)
    code = fields.Char('Reference', copy=False)
    workassembly_id = fields.Many2one('mrp.assemblyline', string='Work Assembly Line')
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

    worksegment_id = fields.Many2one('mrp.worksegament', copy=False)
    hmi_id = fields.Many2one('maintenance.equipment',  string='HMI', copy=False)
    masterpc_id = fields.Many2one('maintenance.equipment',  string='MasterPC', copy=False)

    io_id = fields.Many2one('maintenance.equipment',  string='Remote IO', copy=False)

    rfid_id = fields.Many2one('maintenance.equipment',  string='RFID', copy=False)

    controller_ids = fields.Many2many('maintenance.equipment', 'controller_center_rel', 'center_id', 'controller_id', string='Controller', copy=False)

    gun_ids = fields.Many2many('maintenance.equipment', 'gun_center_rel', 'center_id', 'gun_id',
                                      string='Screw Gun', copy=False)

    controller_ids_domain = fields.Char(
        compute="_compute_controller_ids_domain",
        readonly=True,
        store=False,
    )

    gun_ids_domain = fields.Char(
        compute="_compute_gun_ids_domain",
        readonly=True,
        store=False,
    )

    @api.constrains('controller_ids', 'gun_ids')
    def _constraint_equipments(self):
        self.ensure_one()
        workcenter_ids = self.env['mrp.workcenter'].sudo().search([('id', '!=', self.id)])
        for workcenter in workcenter_ids:
            # org_list = workcenter.controller_ids.ids
            # new_list = self.controller_ids.ids
            # new = len(new_list) if new_list else 0
            # org = len(org_list) if org_list else 0
            # org_list.extend(new_list)
            # if len(set(org_list)) != new + org:
            #     raise ValidationError('控制器设置重复')
            org_list = workcenter.gun_ids.ids
            new_list = self.gun_ids.ids
            new = len(new_list) if new_list else 0
            org = len(org_list) if org_list else 0
            org_list.extend(new_list)
            if len(set(org_list)) != new + org:
                raise ValidationError('拧紧枪设置重复')



    @api.multi
    @api.depends('masterpc_id')
    def _compute_controller_ids_domain(self):
        for rec in self:
            rec.controller_ids_domain = json.dumps([('id', 'in', rec.masterpc_id.child_ids.ids), ('category_name', '=', 'Controller')])
            rec.controller_ids = [(5,)]  # 去除所有的枪 重新设置

    @api.multi
    @api.depends('controller_ids')
    def _compute_gun_ids_domain(self):
        for rec in self:
            child_ids = rec.controller_ids.mapped('child_ids')
            rec.gun_ids_domain = json.dumps([('id', 'in', child_ids.ids), ('category_name', '=', 'Gun')])
            rec.gun_ids = [(5,)]  # 去除所有的枪 重新设置

    _sql_constraints = [('code_hmi', 'unique(hmi_id)', 'Only one HMI is allowed'),
                        ('code_rfid', 'unique(rfid_id)', 'Only one RFID is allowed'),
                        ('code_io', 'unique(io_id)', 'Only one Remote IO is allowed')]