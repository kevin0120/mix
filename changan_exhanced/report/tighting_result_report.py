# -*- coding: utf-8 -*-

from datetime import timedelta
from dateutil.relativedelta import relativedelta
from odoo import api, fields, models, _
from odoo.exceptions import UserError
from itertools import groupby

class ReportTightResultReport(models.AbstractModel):
    _inherit = 'report.spc.report_result_tighting'

    @api.model
    def render_html(self, docids, data=None):
        report_pages = [[]]
        if not data.get('form'):
            raise UserError(_("Form content is missing, this report cannot be printed."))

        Report = self.env['report']
        date_from = fields.Datetime.to_string(fields.Datetime.from_string(data['form']['date_from']))
        date_to = fields.Datetime.to_string(fields.Datetime.from_string(data['form']['date_to']))

        holidays_report = Report._get_report_from_name('spc.report_result_tighting')
        results = self.env['operation.result'].search([('control_date', '>=', date_from), ('control_date', '<=', date_to ) ])
        #
        # if not results:
        #     return True

        for vin, lines in groupby(results, lambda r: r.vin):
            if report_pages[-1] and report_pages[-1][-1]['name']:
                report_pages.append([])
                # Append category to current report page
            report_pages[-1].append({
                'name': vin or 'Uncategorized',
                'lines': list(lines)
            })

        docargs = {
            'doc_ids': self.ids,
            'doc_model': holidays_report.model,
            'docs': report_pages,
            'get_day': fields.Datetime.now(),
        }
        return Report.render('spc.report_result_tighting', docargs)