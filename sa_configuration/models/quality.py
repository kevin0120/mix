# -*- coding: utf-8 -*-
from odoo import fields,models,api
from odoo.exceptions import ValidationError
import json


class QualityPoint(models.Model):
    _inherit = "quality.point"

    operation_id = fields.Many2one('mrp.routing.workcenter', 'Operation')

    operation_id_domain = fields.Char(
        compute="_compute_operation_id_domain",
        readonly=True,
        store=False,
    )

    _sql_constraints = [('product_operation_uniq', 'unique(product_id,operation_id)', 'Only one quailty point per product operation is allowed')]

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
