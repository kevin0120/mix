# -*- coding: utf-8 -*-

from odoo import models, fields, api, _


class MrpProduction(models.Model):
    _inherit = 'mrp.production'

    result_ids = fields.One2many('operation.result', 'production_id',string='Operation Results')