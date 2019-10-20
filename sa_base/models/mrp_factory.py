# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError
import logging

_logger = logging.getLogger(__name__)


class MrpFactory(models.Model):
    _name = 'mrp.factory'
    _description = "MRP Factory"
    _log_access = False
    _order = 'code_related'
    _inherits = {'resource.resource': 'resource_id'}

    code_related = fields.Char(related='resource_id.code', string="Factory Code", store=True)
    name_related = fields.Char(related='resource_id.name', string="Factory Name", store=True)
    country_id = fields.Many2one('res.country', string='Nationality (Country)')

    resource_id = fields.Many2one('resource.resource', string='Resource',
                                  ondelete='cascade', required=True, auto_join=True)

