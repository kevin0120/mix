# -*- coding: utf-8 -*-
from odoo import fields,models,api
from odoo.exceptions import ValidationError
import odoo.addons.decimal_precision as dp
import json


class QualityPoint(models.Model):
    _inherit = "quality.point"

    operation_id = fields.Many2one('mrp.routing.workcenter', 'Operation')

    tolerance_min = fields.Float('Torque Min Tolerance', digits=dp.get_precision('Quality Tests'))
    tolerance_max = fields.Float('Torque Max Tolerance', digits=dp.get_precision('Quality Tests'))

    tolerance_min_degree = fields.Float('Degree Min Tolerance', digits=dp.get_precision('Quality Tests'))
    tolerance_max_degree = fields.Float('Degree Max Tolerance', digits=dp.get_precision('Quality Tests'))

    times = fields.Integer('Repeat times', default=1)

    operation_id_domain = fields.Char(
        compute="_compute_operation_id_domain",
        readonly=True,
        store=False,
    )

    _sql_constraints = [('product_operation_uniq', 'unique(product_id,operation_id)', 'Only one quailty point per product operation is allowed')]

    @api.onchange('operation_id')
    def _onchange_opeartion_id(self):
        self.ensure_one()
        bom_line_ids = self.env['mrp.bom.line'].search([('bom_id.product_id', '=', self.product_id.id), ('operation_id', '=', self.operation_id.id)])
        qtys = [bom_line_id.product_qty for bom_line_id in bom_line_ids]
        self.times = sum(qtys)

    @api.constrains('product_id', 'product_tmpl_id')
    def _product_tmpl_product_constraint(self):
        if self.product_id.product_tmpl_id.id != self.product_tmpl_id.id:
            raise ValidationError('The product template "%s" is invalid on product with name "%s"' % (self.product_tmpl_id.name, self.product_id.name))

    @api.multi
    @api.depends('operation_id', 'product_id','workcenter_id')
    def _compute_operation_id_domain(self):
        for rec in self:
            operation_ids = rec.product_id.bom_ids.mapped('routing_id.operation_ids').ids or []
            rec.operation_id_domain = json.dumps([('workcenter_id', '=', rec.workcenter_id.id), ('id', 'in', operation_ids)])


class QualityCheck(models.Model):
    _inherit = "quality.check"

    measure_degree = fields.Float('Measure Degree', default=0.0, digits=dp.get_precision('Quality Tests'),
                           track_visibility='onchange')

    @api.one
    @api.depends('measure','measure_degree')
    def _compute_measure_success(self):
        if self.point_id.test_type == 'passfail':
            self.measure_success = 'none'
        else:
            if self.measure < self.point_id.tolerance_min or self.measure > self.point_id.tolerance_max:
                self.measure_success = 'fail'
            elif self.measure_degree < self.point_id.tolerance_min_degree or self.measure_degree > self.point_id.tolerance_max_degree:
                self.measure_success = 'fail'
            else:
                self.measure_success = 'pass'