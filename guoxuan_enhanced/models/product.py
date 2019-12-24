# -*- coding: utf-8 -*-

from odoo import models, fields, api


class ProductTemplate(models.Model):
    _inherit = 'product.template'

    sa_type = fields.Selection(selection_add=[('battery', 'Battery Package')],
                               default='battery')
