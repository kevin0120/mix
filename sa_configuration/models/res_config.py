# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import api, fields, models


class SAConfiguration(models.TransientModel):
    _name = 'sa.config.settings'
    _inherit = 'res.config.settings'

    generate_result_sequence = fields.Selection([
        (0, "Set Sequence by Operation(default)"),
        (1, "Set Sequence by Per Vehicle")
        ], string="Result Sequences")

    @api.multi
    def set_default_generate_result_sequence(self):
        check = self.env.user.has_group('base.group_system')
        Values = check and self.env['ir.values'].sudo() or self.env['ir.values']
        for config in self:
            Values.set_default('sa.config.settings', 'generate_result_sequence', config.generate_result_sequence)