# -*- coding: utf-8 -*-


from odoo import api, fields, models, SUPERUSER_ID, _

from odoo.exceptions import UserError, ValidationError

from odoo.addons import decimal_precision as dp

import json


class OperationPointsGoup(models.Model):
    _name = 'operation.point.group'

    _order = "sequence"

    sequence = fields.Integer('sequence', default=1)

    name = fields.Char('Operation Point Group')

    # operation_point_ids_domain = fields.Char(
    #     compute="_compute_operation_point_ids_domain",
    #     readonly=True,
    #     store=False,
    # )

    operation_id = fields.Many2one('mrp.routing.workcenter', ondelete='cascade', index=True)

    operation_point_ids = fields.One2many('operation.point', 'group_id',
                                           string='Points', copy=False)

    # @api.constrains(operation_point_ids)
    # def _constraint_operation_point_ids(self):
    #     point_ids = self.operation_point_ids.ids
    #     if len(point_ids) != len(set(point_ids)):
    #         raise ValidationError(u'作业点设定中存在重复项')

    @api.multi
    def name_get(self):
        res = []
        for point in self:
            res.append((point.id, _('[%s] %s') % (point.operation_id.name, point.name)))
        return res


class OperationPoints(models.Model):
    _name = 'operation.point'

    _order = "sequence"

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    sequence = fields.Integer('sequence', default=1)

    name = fields.Char('Operation Point Name', default='Bolt Point')

    group_id = fields.Many2one('operation.point.group')

    group_sequence = fields.Integer('Group Sequence for Multi Spindle')

    product_id = fields.Many2one('product.product', 'Product')

    product_qty = fields.Float('Product Quantity', default=1.0, digits=dp.get_precision('Product Unit of Measure'))

    x_offset = fields.Float('x axis offset from left(%)', default=0.0, digits=dp.get_precision('POINT_OFFSET'))

    y_offset = fields.Float('y axis offset from top(%)', default=0.0, digits=dp.get_precision('POINT_OFFSET'))

    program_id = fields.Many2one('controller.program',  string='程序号', ondelete='cascade')

    operation_id = fields.Many2one('mrp.routing.workcenter', ondelete='cascade', index=True)

    max_redo_times = fields.Integer('Operation Max Redo Times', default=3)  # 此项重试业务逻辑在HMI中实现

    @api.multi
    def name_get(self):
        res = []
        for point in self:
            res.append((point.id, _('[%s] %s') % (point.operation_id.name, point.sequence)))
        return res

    @api.model
    def default_get(self, fields):
        res = super(OperationPoints, self).default_get(fields)

        operation_id = self.env.context.get('default_operation_id')
        if operation_id:
            operation = self.env['mrp.routing.workcenter'].sudo().browse(operation_id)
            if 'max_redo_times' in fields:
                res.update({'max_redo_times': operation.max_redo_times})
            if 'sequence' in fields and operation.operation_point_ids:
                res.update({'sequence': max(operation.operation_point_ids.mapped('sequence')) + 1})
        return res

    @api.multi
    def unlink(self):
        if self.env.uid != SUPERUSER_ID:
            raise UserError(_(u"Only SuperUser can delete program"))
        return super(OperationPoints, self).unlink()

    @api.one
    def toggle_active(self):
        bom_line_id = self.env['mrp.bom.line'].search([('operation_point_id', '=', self.id)])
        if bom_line_id:
            bom_line_id.toggle_active()
        return super(OperationPoints, self).toggle_active()