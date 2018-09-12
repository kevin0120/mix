# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from datetime import timedelta
from dateutil.relativedelta import relativedelta
from odoo import api, fields, models, _
from odoo.exceptions import UserError
from itertools import groupby

class TightResultReport(models.TransientModel):
    _name = 'tight.result.report'
    _description = u'拧紧结果报告'

    def _get_default_date_from(self):
        year = fields.Date.from_string(fields.Date.today()).strftime('%Y')
        return '{}-01-01'.format(year)

    def _get_default_date_to(self):
        date = fields.Date.from_string(fields.Date.today())
        return date.strftime('%Y') + '-' + date.strftime('%m') + '-' + date.strftime('%d')

    date_from = fields.Date(string='Start Date', required=True, default=_get_default_date_from)
    date_to = fields.Date(string='End Date', required=True, default=_get_default_date_to)

    @api.multi
    def print_report(self):
        """
         To get the date and print the report
         @return: return report
        """
        self.ensure_one()
        data = {'ids': self.env.context.get('active_ids', [])}
        res = self.read()
        res = res and res[0] or {}
        data.update({'form': res})
        return self.env['report'].get_action(self, 'spc.report_result_tighting', data=data)


class ReportTightResultReport(models.AbstractModel):
    _name = 'report.spc.report_result_tighting'

    @api.model
    def render_html(self, docids, data=None):
        report_pages = [[]]
        if not data.get('form'):
            raise UserError(_("Form content is missing, this report cannot be printed."))

        Report = self.env['report']
        date_from = fields.Datetime.to_string(fields.Datetime.from_string(data['form']['date_from']))
        date_to = fields.Datetime.to_string(fields.Datetime.from_string(data['form']['date_to']))

        holidays_report = Report._get_report_from_name('spc.report_result_tighting')
        results = self.env['operation.result'].search([ ('control_date', '>=', date_from), ('control_date', '<=', date_to ) ])

        for category, lines in groupby(results, lambda r: r.production_id):
            if report_pages[-1] and report_pages[-1][-1]['name']:
                report_pages.append([])
                # Append category to current report page
            report_pages[-1].append({
                'name': category and category.vin or 'Uncategorized',
                'lines': list(lines)
            })

        docargs = {
            'doc_ids': self.ids,
            'doc_model': holidays_report.model,
            'docs': report_pages,
            'get_day': fields.Datetime.now(),
        }
        return Report.render('spc.report_result_tighting', docargs)