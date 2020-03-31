# -*- coding: utf-8 -*-
from odoo import api, fields, models


class ResUsers(models.Model):
    _inherit = 'res.users'

    uuid = fields.Char(string='UUID', required=True)
    working_status = fields.Selection([('ready', 'Ready'),
                                      ('doing', 'Doing')])
