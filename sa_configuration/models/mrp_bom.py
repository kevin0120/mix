# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.exceptions import ValidationError


class MrpBom(models.Model):
    _inherit = 'mrp.bom'

    _sql_constraints = [('product_uniq', 'unique(product_id, active)', u'当前产品只有唯一激活的BOM'),
                        ('product_tmpl_uniq', 'unique(product_tmpl_id, active)', u'当前产品类型只有唯一激活的BOM')]

    @api.constrains('product_id', 'product_tmpl_id')
    def _product_tmpl_product_constraint(self):
        if self.product_id.product_tmpl_id.id != self.product_tmpl_id.idsa_mrp_bom_form_view:
            raise ValidationError('The product template "%s" is invalid on product with name "%s"' % (self.product_tmpl_id.name, self.product_id.name))