# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    worksheet_img = fields.Binary('worksheet_img')
