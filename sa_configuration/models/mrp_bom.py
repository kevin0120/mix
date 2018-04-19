# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpBom(models.Model):
    _inherit = 'mrp.bom'

    _sql_constraints = [('product_uniq', 'unique(product_id, active)', u'当前产品只有唯一激活的BOM'),
                        ('product_tmpl_uniq', 'unique(product_tmpl_id, active)', u'当前产品类型只有唯一激活的BOM')]