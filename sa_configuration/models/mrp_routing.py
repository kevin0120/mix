# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    def _get_program_set(self):
        return self.routing_id.program_set

    program_set = fields.Integer(string='程序号', related='routing_id.program_set', store=True)


class MrpRouting(models.Model):
    """ Specifies routings of work centers """
    _inherit = 'mrp.routing'

    worksheet_img = fields.Binary('worksheet_img')

    program_set = fields.Integer(string='程序号')
