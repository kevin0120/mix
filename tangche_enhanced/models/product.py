# -*- coding: utf-8 -*-

from odoo import models, fields, api


class ProductTemplate(models.Model):
    _inherit = 'product.template'

    sa_type = fields.Selection(
        selection_add=[('carriage', 'Train Carriage'), ('consu', 'Consumable Component'), ('socket', 'Bolt Socket')],
        default='carriage')


class ProductProduct(models.Model):
    _inherit = 'product.product'

    socket = fields.Many2one('product.product', domain=[('sa_type', '=', 'socket')])
