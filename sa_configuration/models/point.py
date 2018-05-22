# -*- coding: utf-8 -*-


from odoo import api, fields, models


class Points(models.Model):
    _name = 'point.point'

    x_offset = fields.Integer('x axis offset from left(px)', default=0)

    y_offset = fields.Integer('y axis offset from top(px)', default=0)

    res_model = fields.Char("related Model", default='ir.attachment')

    res_field = fields.Char('Related filed')

    res_id = fields.Integer('Related Model ID')