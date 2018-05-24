# -*- coding: utf-8 -*-

from odoo import models, fields, api


# class ControllerResult(models.TransientModel):
#     _name = "controller.result"
#
#     workorder_id = fields.Many2one('mrp.workorder', 'Operation')
#     assembly_line_id = fields.Many2one('mrp.assemblyline', related='workorder_id.production_id.assembly_line_id',string='Assembly Line ID', store=True, readonly=True)
#     workcenter_id = fields.Many2one('mrp.workcenter')
#     worksegment_id = fields.Many2one('mrp.worksegament', related='workcenter_id.worksegment_id', store=True, readonly=True)
#     wizard_id = fields.Many2one('operation.result.wizard', ondelete='cascade')
#     pass_num = fields.Float('Amount of Pass')
#     fail_num = fields.Float('Amount of Fail')
#     amount_screw = fields.Float('Amount of Screw', compute='_compute_amount_of_screw', store=True)
#     pecent_pass = fields.Float('Pecent of Pass', digits=(16, 2), compute='_compute_amount_of_screw', store=True)
#
#     @api.multi
#     @api.depends('pass_num', 'fail_num')
#     def _compute_amount_of_screw(self):
#         for r in self:
#             r.amount_screw = r.pass_num + r.fail_num
#             r.pecent_pass = r.pass_num / r.amount_screw