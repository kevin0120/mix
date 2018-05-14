# -*- coding: utf-8 -*-

from odoo import models, fields, exceptions, _


class AiisSettings(models.TransientModel):
    _name = 'aiis.config.settings'
    _inherit = 'res.config.settings'

    aiis_url = fields.Char('Aiis URL', default='http://127.0.0.1:8080')

    def get_default_all(self, fields):
        aiis_url = self.env["ir.config_parameter"].get_param("aiis.url", default='http://127.0.0.1:9000')

        return dict(
            aiis_url=aiis_url,
        )

    # minio_url
    def set_minio_url(self):
        self.env['ir.config_parameter'].set_param("aiis.url",
                                                  self.aiis_url or 'http://127.0.0.1:9000',
                                                  groups=['base.group_system'])