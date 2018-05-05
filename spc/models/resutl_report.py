# -*- coding: utf-8 -*-

from odoo import tools
from odoo import models, fields, api


class OperationResultReport(models.Model):
    _name = "operation.result.report"
    _description = "Result Statistics"
    _auto = False
    _rec_name = 'date'

    @api.model_cr
    def init(self):
        tools.drop_view_if_exists(self._cr, self._table)
        self._cr.execute('''
        create or replace view asset_asset_report as (
          SELECT 
        )''')

        results = self.env.cr.fetchall()
        for result in results:
            print result