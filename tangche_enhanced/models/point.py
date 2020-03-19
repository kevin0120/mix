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

    @api.model
    def default_get(self, fields):
        res = super(OperationPoints, self).default_get(fields)
        if 'picking_type_id' not in res:
            res.update({
                'picking_type_id': self._get_default_picking_type()
            })

        operation_id = self.env.context.get('default_operation_id')
        if operation_id:
            operation = self.env['mrp.routing.workcenter'].sudo().browse(operation_id)
            if 'max_redo_times' in fields:
                res.update({'max_redo_times': operation.max_redo_times})
            if 'sequence' in fields and operation.operation_point_ids:
                res.update({'sequence': max(operation.operation_point_ids.mapped('sequence')) + 1})
        # parent_qcp_id = self.env.context.get('default_parent_qcp_id')
        # if parent_qcp_id:
        #     qcp_id = self.env['sa.quality.point'].sudo().browse(parent_qcp_id)
        #     if 'sa_operation_ids' in fields:
        #         res.update({'sa_operation_ids': [(6, 0, qcp_id.sa_operation_ids.ids)]})
        return res
