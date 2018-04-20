# -*- coding: utf-8 -*-

from odoo import models, fields, api


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    def _default_get_program_set(self):
        routing_id = self._context.get('default_routing_id',False)
        return self.env['mrp.routing'].browse(routing_id).program_set

    program_set = fields.Char(string='程序号', default=_default_get_program_set)


class MrpRouting(models.Model):
    """ Specifies routings of work centers """
    _inherit = 'mrp.routing'

    operation_count = fields.Integer(string='Operations', compute='_compute_operation_count')

    worksheet_img = fields.Binary('worksheet_img')

    program_set = fields.Char(string='程序号')

    @api.depends('operation_ids')
    def _compute_operation_count(self):
        for routing in self:
            routing.operation_count = len(routing.operation_ids)

    @api.multi
    def action_sa_view_operation(self):
        action = self.env.ref('sa_configuration.sa_mrp_routing_workcenter_action').read()[0]
        workcenter_id = self.env.ref('sa_configuration.cunrong_default_workcenter').id
        ids = self.ids
        # bom specific to this variant or global to template
        action['context'] = {
            'default_routing_id': ids[0],
            'default_workcenter_id': workcenter_id
        }
        action['domain'] = [('routing_id', 'in', [self.ids])]
        return action
