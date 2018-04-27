# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.exceptions import ValidationError


class MrpBom(models.Model):
    _inherit = 'mrp.bom'

    routing_group_id = fields.Many2one('mrp.routing.group', related='routing_id.group_id', readonly=True)

    _sql_constraints = [('product_routing_uniq', 'unique(product_id, routing_id, active)', u'当前产品只有唯一激活工艺组的BOM'),
                        ('product_tmpl_routing_uniq', 'unique(product_tmpl_id,routing_id, active)', u'当前产品类型只有唯一激活工艺组的BOM')]

    @api.constrains('product_id', 'product_tmpl_id')
    def _product_tmpl_product_constraint(self):
        if self.product_id.product_tmpl_id.id != self.product_tmpl_id.id:
            raise ValidationError('The product template "%s" is invalid on product with name "%s"' % (self.product_tmpl_id.name, self.product_id.name))