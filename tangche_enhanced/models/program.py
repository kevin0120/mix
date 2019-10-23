# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError


class ControllerProgram(models.Model):
    _inherit = 'controller.program'

    @api.onchange('code')
    def _onchange_program_code(self):
        try:
            for program in self:
                program.code = int(program.code).__str__()
        except Exception:
            raise ValidationError(_('Program: {0} Code: {1} Must Be Integer!!!'.format(
                program.name, program.code)))
