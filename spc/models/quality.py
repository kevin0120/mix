# -*- coding: utf-8 -*-

from odoo import models, fields, api


class QualityPoint(models.Model):
    _inherit = "sa.quality.point"

    @api.multi
    def action_see_spc_control(self):
        action = self.env.ref('spc.quality_check_action_spc').read()[0]
        action.update({
            'context': {
                'search_default_qcp_id': self.id
            }
        })
        return action

    @api.multi
    def action_see_result(self):
        action = self.env.ref('spc.operation_result_action_main').read()[0]
        action.update({
            'context': {
                'search_default_qcp_id': self.id
            }
        })
        return action


class QualityCheck(models.Model):
    _inherit = "sa.quality.check"

    @api.depends('result_ids')
    def _get_operation_results(self):
        for check in self:
            check.result_count = len(check.result_ids)

    result_ids = fields.One2many('operation.result', 'quality_check_id', string='Operation Results')

    result_count = fields.Integer(string='Operation Result Count', compute=_get_operation_results)
