# -*- coding: utf-8 -*-
from scipy.stats import norm
from odoo import models, fields, api
from pyecharts import Overlap, Bar, Line, Grid
import pyecharts
from pandas import DataFrame


class TorAngSPCReport(models.TransientModel):
    _name = "ta.spc.wizard"
    _description = "Result Statistics Process Control"

    query_date_from = fields.Datetime(string='Query Date From')
    query_date_to = fields.Datetime(string='Query Date to', default=fields.Datetime.now())
    vehicle_id = fields.Many2one('product.product', string='Vehicle Type', domain=[('sa_type', '=', 'vehicle')])
    screw_id = fields.Many2one('product.product', string='Screw Type', domain=[('sa_type', '=', 'screw')])
    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line')
    limit = fields.Integer('Query Limit', default=5000)
    spc_target = fields.Selection([('torque', 'Torque'), ('angle', 'Angle')], string='统计对象', default='torque')

    normal_dist = fields.Text(string='Normal Distribution', store=False)
    weibull_dist = fields.Text(string='Weibull Distribution', store=False)
    scatter = fields.Text(string='Scatter', store=False)

    def _get_normal_dist(self,data):
        grid = Grid(width=1920, height=1080)
        attr = ["{}月".format(i) for i in range(1, 13)]
        v1 = [2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3]
        v2 = [2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3]
        v3 = [2.0, 2.2, 3.3, 4.5, 6.3, 10.2, 20.3, 23.4, 23.0, 16.5, 12.0, 6.2]

        bar = Bar(title=u"Normal Distribution({0})".format(self.spc_target), title_pos="50%")
        bar.add("蒸发量", attr, v1)
        bar.add("降水量", attr, v2, yaxis_formatter=" ml", yaxis_max=250,
                legend_pos="85%", legend_orient="vertical", legend_top="45%")
        line = Line()
        line.add("平均温度", attr, v3, yaxis_formatter=" °C", is_smooth=True )
        overlap = Overlap(width=1920, height=1080)
        overlap.add(bar)
        overlap.add(line, is_add_yaxis=True, yaxis_index=1)

        grid.add(overlap)
        pyecharts.configure(force_js_embed=True)
        return grid.render_embed()
        # bar = Bar(u'Normal Distribution', title_pos='center', width=1920,height=1080)
        # style = Style()
        # pie_style = style.add(label_pos="center", is_label_show=True, label_text_color=None)
        #
        # bar.add("", ["剧情", "xx"], [25, 75], center=[10, 30],
        #         radius=[18, 24], legend_pos="left", legend_orient='vertical', **pie_style)
        # bar.add("", ["奇幻", "yy"], [24, 76], center=[30, 30],
        #         radius=[18, 24],  legend_pos="right", legend_orient='vertical', **pie_style)
        # pyecharts.configure(force_js_embed=True)
        # return bar.render_embed()

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        result = super(TorAngSPCReport, self).read(fields, load=load)
        if 'normal_dist' in fields and load == '_classic_read':
            data = self._get_data()
            result[0].update({'normal_dist': self._get_normal_dist(data)})
        return result

    def _get_data(self):
        domain = [('measure_result', 'in', ['ok', 'nok'])]
        fields = ['measure_torque']
        if self.query_date_from:
            domain += [('control_date', '>=', self.query_date_from)]
        if self.query_date_to:
            domain += [('control_date', '<=', self.query_date_to)]
        if self.vehicle_id:
            domain += [('product_id', '=', self.vehicle_id.id)]
        if self.screw_id:
            domain += [('consu_product_id', '=', self.screw_id.id)]
        if self.assembly_line_id:
            domain += [('assembly_line_id', '=', self.assembly_line_id.id)]
        if self.spc_target == 'torque':
            fields = ['measure_torque']
        if self.spc_target == 'angle':
            fields = ['measure_degree']
        return self.env['operation.result'].sudo().search_read(domain=domain, fields=fields, limit=self.limit)

    # @api.multi
    # def button_query_vehicle(self):
    #     result = {'count': 0,
    #               'ok': 0,
    #               'lacking': 0,
    #               'nok': 0,
    #               'used': 0}
    #     knr_code = '%' + self.knr_code + '%' if self.knr_code else '%'
    #     query = """
    #                       SELECT b.knr as knr, count(*) as count, o.measure_result as result ,o.lacking as lack
    #                       FROM operation_result o
    #                       FULL JOIN mrp_production b ON (b.id = o.production_id)
    #                       WHERE b.knr LIKE '%s'
    #                       AND o.control_date >= '%s'
    #                       AND o.control_date <= '%s'
    #                       group by o.measure_result, o.lacking, b.knr
    #                     """ % (knr_code, self.query_date_from, self.query_date_to)
    #     self.env.cr.execute(query, ())
    #     data = [row for row in self.env.cr.dictfetchall()]
    #     df = DataFrame.from_dict(data)
    #     print(df)
    #     _df = df.groupby('knr')
    #     print(_df)
    #
    #     result['used'] = result['count'] - result['lacking']
    #     return result

    @api.multi
    def button_query(self):
        pass