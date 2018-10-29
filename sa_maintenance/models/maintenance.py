# -*- coding: utf-8 -*-

from odoo import models, fields, api
import odoo.addons.decimal_precision as dp


class MaintenanceCheckPointCategory(models.Model):
    _name = 'maintenance.cp.category'

    name = fields.Char('Check Point Name')
    code = fields.Char('Check Point Code')

    test_type = fields.Selection([
        ('passfail', 'Pass - Fail'),
        ('measure', 'Measure')], string="Test Type",
        default='passfail', required=True)


class MaintenanceCheckPoint(models.Model):
    _name = 'maintenance.cp'
    _description = 'Equipment Checklist for each maintenance request '

    equipment_id = fields.Many2one('maintenance.equipment', string='Equipment', index=True)

    category_id = fields.Many2one('maintenance.cp.category')
    description = fields.Html('Maintenance Check Point Description')

    test_type = fields.Selection(related='category_id.test_type', readonly=True)

    norm = fields.Float('Norm', digits=dp.get_precision('Maintenance Tests'))  # TDE RENAME ?
    tolerance_min = fields.Float('Min Tolerance', digits=dp.get_precision('Maintenance Tests'))
    tolerance_max = fields.Float('Max Tolerance', digits=dp.get_precision('Maintenance Tests'))

    @api.onchange('norm')
    def onchange_norm(self):
        if self.tolerance_max == 0.0:
            self.tolerance_max = self.norm


class MaintenanceCheckPointAction(models.Model):
    _name = 'maintenance.cp.action'

    request_id = fields.Many2one('maintenance.request', string='Request', index=True)

    point_id = fields.Many2one('maintenance.cp', 'Check Point')
    measure = fields.Float(digits=dp.get_precision('Maintenance Tests'))

    measure_success = fields.Selection([
        ('none', 'No measure'),
        ('pass', 'Pass'),
        ('fail', 'Fail')], string="Measure Success", compute="_compute_measure_success",
        readonly=True, store=True)

    @api.one
    @api.depends('measure')
    def _compute_measure_success(self):
        if self.point_id.test_type == 'passfail':
            self.measure_success = 'none'
        else:
            if self.measure < self.point_id.tolerance_min or self.measure > self.point_id.tolerance_max:
                self.measure_success = 'fail'
            else:
                self.measure_success = 'pass'


class MaintenanceRequest(models.Model):
    _inherit = 'maintenance.request'

    check_point_action_ids = fields.One2many('maintenance.cp.action', 'request_id')

    @api.model
    def create(self, vals):
        ret = super(MaintenanceRequest, self).create(vals)
        if 'equipment_id' in vals and vals.get('maintenance_type', 'corrective') == 'preventive':
            equipment_id = vals.get('equipment_id')
            check_point_ids = self.env['maintenance.cp'].sudo().search([('equipment_id', '=', equipment_id)])

        return ret


class MaintenanceEquipment(models.Model):
    _inherit = 'maintenance.equipment'

    check_point_ids = fields.One2many('maintenance.cp', 'equipment_id')

