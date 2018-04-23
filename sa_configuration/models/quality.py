# -*- coding: utf-8 -*-
from odoo import fields,models


class QualityPoint(models.Model):
    _inherit = "quality.point"

    operation_id = fields.Many2one('mrp.routing.workcenter', 'Operation')