# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError


class QualityPoint(models.Model):
    _inherit = "sa.quality.point"

    @api.constrains('sa_operation_ids')
    def _ts002_constrain_operation_ids(self):
        for step in self:
            if len(step.sa_operation_ids) > 1:
                raise ValidationError(u'唐山客车 每个工步最多定义一个作业')
