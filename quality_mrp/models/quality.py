# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import fields, models


class QualityCheck(models.Model):
    _inherit = "sa.quality.check"

    production_id = fields.Many2one(
        'mrp.production', 'Production Order', check_company=True)


class QualityAlert(models.Model):
    _inherit = "sa.quality.alert"

    production_id = fields.Many2one(
        'mrp.production', "Production Order", check_company=True)
