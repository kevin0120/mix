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
