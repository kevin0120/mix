# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    worksheet_img = fields.Binary(
        'Worksheet', related='operation_id.worksheet_img', readonly=True)