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


class SPCSetting(models.TransientModel):
    _name = 'spc.config.settings'
    _inherit = 'res.config.settings'

    scale_parameter = fields.Float('Weibull Distribution Scale Parameter', default=1.0)

    shape_parameter = fields.Float('Weibull Distribution Shape Parameter', default=5.0)

    def get_default_all(self, fields):
        scale_parameter = self.env["ir.config_parameter"].get_param("weibull.scale", default=1.0)

        shape_parameter = self.env["ir.config_parameter"].get_param("weibull.shape", default=5.0)

        return dict(
            scale_parameter=scale_parameter,
            shape_parameter=shape_parameter
        )

    def set_scale_parameter(self):
        self.env['ir.config_parameter'].set_param("weibull.scale",
                                                  self.scale_parameter or 1.0,
                                                  groups=['spc.group_spc_user'])

    def set_shape_parameter(self):
        self.env['ir.config_parameter'].set_param("weibull.shape",
                                                  self.shape_parameter or 5.0,
                                                  groups=['spc.group_spc_user'])
