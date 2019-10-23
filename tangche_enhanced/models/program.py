# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError


class ControllerProgram(models.Model):
    _inherit = 'controller.program'

    @api.constrains('code')
    def _constraint_program_code(self):
        try:
            for program in self:
                program.code = int(program.code)
        except Exception:
            raise ValidationError(_('Program: {0} Code Must Be Integer!!!'.format(program.name)))
