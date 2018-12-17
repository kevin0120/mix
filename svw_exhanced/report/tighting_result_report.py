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
        # results = self.env['operation.result'].search([ ('control_date', '>=', date_from), ('control_date', '<=', date_to ) ])

        query = "select tb.vin, json_agg(json_build_array(tb.control_date, tb.code, tb.measure_result)) " \
                "from (select ors.vin as vin, ors.control_date as control_date, ors.measure_result as measure_result, cps.code as code from operation_result ors join controller_program cps on cps.id = ors.program_id where ors.control_date >= '%s' and ors.control_date <= '%s' order by ors.control_date) tb " \
                "group by tb.vin" % (
            date_from,
            date_to
        )

        cr = self._cr
        cr.execute(query)
        results = cr.fetchall()

        #
        # if not results:
        #     return True

        # for category, lines in groupby(results, lambda r: r[1]):
        #     if report_pages[-1] and report_pages[-1][-1]['name']:
        #         report_pages.append([])
        #         # Append category to current report page
        #     report_pages[-1].append({
        #         'name': lines[0][0] and lines[0][0] or 'Uncategorized',
        #         'lines': list(lines)
        #     })

        docargs = {
            'doc_ids': self.ids,
            'doc_model': holidays_report.model,
            'docs': results,
            'get_day': fields.Datetime.now(),
        }

        return Report.render('spc.report_result_tighting', docargs)