# -*- coding: utf-8 -*-

from odoo import models, fields, api
import time

class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    worksheet_img = fields.Binary(
        'Worksheet', related='operation_id.worksheet_img', readonly=True)

    result_ids = fields.One2many('operation.result', 'workorder_id',string='Operation Results')

    sent = fields.Boolean('Have sent to aiis', default=False)

    @api.model
    def create(self, vals):
        order = super(MrpWorkorder, self).create(vals)
        return order

    @api.multi
    def _create_checks(self):
        ret_vals = []
        for wo in self:
            production = wo.production_id
            points = self.env['quality.point'].search([('workcenter_id', '=', wo.workcenter_id.id),
                                                       ('operation_id','=', wo.operation_id.id),  # 定位到某个作业的质量控制点
                                                       ('picking_type_id', '=', production.picking_type_id.id),
                                                       '|', ('product_id', '=', production.product_id.id),
                                                       '&', ('product_id', '=', False), ('product_tmpl_id', '=', production.product_id.product_tmpl_id.id)])
            for point in points:
                vals = {'workorder_id': wo.id,
                         'production_id': production.id,
                         'workcenter_id': wo.workcenter_id.id,
                         'assembly_line_id':production.assembly_line_id.id,
                         'point_id': point.id,
                         'product_id': production.product_id.id,
                         'consu_product_id': wo.consu_product_id.id,
                         'time': production.date_planned_start or fields.Datetime.now(),
                         'control_date': fields.Datetime.now()}

                # if point.check_execute_now():
                map(lambda i: ret_vals.append(vals), range(point.times))

        self.env['operation.result'].sudo().bulk_create(ret_vals)