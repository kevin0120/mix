# -*- coding: utf-8 -*-

from odoo import models, fields, api
from odoo.exceptions import ValidationError


class MrpProduction(models.Model):
    _inherit = 'mrp.production'

    equipment_name = fields.Char(string='Equipment Name')
    factory_name = fields.Char(string='Factory Name')
    year = fields.Integer(string='Year')
    pin = fields.Integer(string='PIN')
    vin = fields.Char(string='VIN')
    pin_check_code = fields.Integer(string='PIN check Code')
    assembly_line = fields.Char(string='Assembly Line ID')
    lnr = fields.Char(string='Line Number')
    knr = fields.Char(string='Vehicle Assembly Code', store=True, compute='_compute_mo_knr')

    _sql_constraints = [('vin_uniq', 'unique(vin)', 'Only one VIN per MO is allowed')]

    @api.depends('pin','pin_check_code')
    def _compute_mo_knr(self):
        ### 只会在创建记录时计算一次
        for mo in self:
            mo.knr = u'{0}{1}'.format(mo.pin, mo.pin_check_code)

    @api.constrains('year')
    def _constraint_mo_year(self):
        if len(str(self.year)) != 4:
            raise ValidationError(u'不是年份')



class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    worksheet_img = fields.Binary(
        'Worksheet', related='operation_id.worksheet_img', readonly=True)

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
                    for i in range(point.times):
                        self.env['quality.check'].create({'workorder_id': wo.id,
                                                          'production_id': production.id,
                                                      'point_id': point.id,
                                                      'team_id': point.team_id.id,
                                                      'product_id': production.product_id.id})