# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    routing_id = fields.Many2one(
        'mrp.routing', string='Operation', related='operation_id.routing_id', readonly=1)

    worksheet_img = fields.Binary(
        'Worksheet', related='routing_id.worksheet_img', readonly=True)

    @api.multi
    def _create_checks(self):
        for wo in self:
            production = wo.production_id
            points = self.env['quality.point'].search([('workcenter_id', '=', wo.workcenter_id.id),
                                                       ('operation_id','=', wo.operation_id.id),  # 定位到某个作业的质量控制点
                                                       ('picking_type_id', '=', production.picking_type_id.id),
                                                       '|', ('product_id', '=', production.product_id.id),
                                                       '&', ('product_id', '=', False), ('product_tmpl_id', '=', production.product_id.product_tmpl_id.id)])
            for point in points:
                if point.check_execute_now():
                    self.env['quality.check'].create({'workorder_id': wo.id,
                                                      'point_id': point.id,
                                                      'team_id': point.team_id.id,
                                                      'product_id': production.product_id.id})