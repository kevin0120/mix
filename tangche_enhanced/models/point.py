# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.exceptions import ValidationError


class OperationPoints(models.Model):
    _inherit = 'operation.point'

    socket = fields.Many2one('product.product', domain=[('sa_type', '=', 'socket')])

    parent_test_type = fields.Char('Parent Test Type', related='parent_qcp_id.test_type')

    # @api.model
    # def default_get(self, fields):
    #     res = super(OperationPoints, self).default_get(fields)
    #     parent_test_type = self.env.context.get('default_parent_test_type')
    #     if parent_test_type:
    #         res.update({'parent_test_type': parent_test_type})
    #     return res

    @api.onchange('product_id')
    def onchange_product_id(self):
        for point in self:
            if point.product_id.socket:
                point.socket = point.product_id.socket

    @api.constrains('tightening_tool_ids')
    def _constraint_tightening_tool_ids(self):
        context_parent_test_type = self.env.context.get('parent_test_type', False)
        for point in self:
            parent_test_type = point.parent_test_type or context_parent_test_type
            if not point.tightening_tool_ids:
                continue
            workcenter_ids = set(point.tightening_tool_ids.mapped('workcenter_id').ids)
            if parent_test_type == 'promiscuous_tightening':
                if len(workcenter_ids) != 1:
                    raise ValidationError(u'混杂模式下的拧紧工具必须是同一工位的')
                else:
                    continue
            if len(workcenter_ids) != len(point.tightening_tool_ids):
                raise ValidationError(u'不能对同一个拧紧点选择同一个工位上的拧紧工具')

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
