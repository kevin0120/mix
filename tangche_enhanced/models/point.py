# -*- coding: utf-8 -*-

from odoo import models, fields, api


class OperationPoints(models.Model):
    _inherit = 'operation.point'

    socket = fields.Many2one('product.product', domain=[('sa_type', '=', 'socket')])

    @api.onchange('product_id')
    def onchange_product_id(self):
        for point in self:
            if point.product_id.socket:
                point.socket = point.product_id.socket
