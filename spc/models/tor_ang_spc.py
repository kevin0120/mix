# -*- coding: utf-8 -*-
from scipy.stats import norm, dweibull, weibull_max, weibull_min, invweibull,exponweib
from odoo import models, fields, api
from odoo.exceptions import UserError,ValidationError
from pyecharts import Overlap, Bar, Line, Grid
import pyecharts
from pandas import DataFrame
import numpy as np
from odoo.tools import float_round


DEFAULT_LIMIT = 5000
MIN_LIMIT = 1000

class TorAngSPCReport(models.TransientModel):
    _name = "ta.spc.wizard"
    _description = "Result Statistics Process Control"

    query_date_from = fields.Datetime(string='Query Date From')
    query_date_to = fields.Datetime(string='Query Date to', default=fields.Datetime.now())
    vehicle_id = fields.Many2one('product.product', string='Vehicle Type', domain=[('sa_type', '=', 'vehicle')])
    screw_id = fields.Many2one('product.product', string='Screw Type', domain=[('sa_type', '=', 'screw')])
    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line')
    limit = fields.Integer('Query Limit', default=DEFAULT_LIMIT)
    spc_target = fields.Selection([('torque', 'Torque'), ('angle', 'Angle')], string='统计对象', default='torque')

    normal_dist = fields.Text(string='Normal Distribution', store=False)
    weibull_dist = fields.Text(string='Weibull Distribution', store=False)
    weibull_dist_method = fields.Selection([('double', 'Double Weibull'),('inverted', 'Inverted'), ('exponential', 'Exponential'),('min', 'Min Weibull'), ('max', 'Max Weibull')], string='韦伯分布统计方法', default='min')
    scatter = fields.Text(string='Scatter', store=False)


    @api.constrains('limit')
    def _constraint_limit(self):
        if self.limit < MIN_LIMIT:
            raise UserError(u'查询数量不得小于最小查询数量:{0}'.format(MIN_LIMIT))

    def _get_weibull_dist(self, qty, mean=None, std=None, scale=1.0, shape=5.0):

        x_line = np.arange(mean - std * 4.0, mean + std * 5.0, 1 * std)

        if self.weibull_dist_method == 'double':
            _data = dweibull(shape, loc=mean, scale=std)
            y_line = _data.pdf(x_line) * qty

        if self.weibull_dist_method == 'inverted':
            _data = invweibull(shape, loc=mean, scale=std)
            y_line = _data.pdf(x_line) * qty

        if self.weibull_dist_method == 'exponential':
            _data = exponweib(scale, shape, loc=mean, scale=std)
            y_line = _data.pdf(x_line) * qty

        if self.weibull_dist_method == 'min':
            _data = weibull_min(shape, loc=mean, scale=std)
            y_line = _data.pdf(x_line) * qty

        if self.weibull_dist_method == 'max':
            _data = weibull_max(shape, loc=mean, scale=std)
            y_line = _data.pdf(x_line) * qty

        line = Line(width=1280, height=600)
        line.add(u'{0}'.format(self.spc_target), x_line, y_line, xaxis_name=u'{0}'.format(self.spc_target),
                 yaxis_name=u'数量(Quantity)',
                 line_color='rgba(0 ,255 ,127,0.5)', legend_pos='center',
                 is_smooth=True, line_width=2,
                 tooltip_tragger='axis', is_fill=True, area_color='#20B2AA', area_opacity=0.4)
        pyecharts.configure(force_js_embed=True)
        return line.render_embed()

    def _get_normal_dist(self, mean=None, std=None):
        norm_data = norm(mean, std)
        x_bar = np.arange(mean - std * 3.0, mean + std * 4.0, 1 * std)
        y_bar = norm_data.pdf(x_bar)

        x_line = np.arange(mean - std * 3.0, mean + std * 4.0, 1 * std)
        y_line = norm_data.pdf(x_line)

        bar = Bar(title=u"Normal Distribution({0})".format(self.spc_target), title_pos="50%", width=960, height=1440)
        bar.add(u'{0}'.format(self.spc_target), x_bar, y_bar, legend_orient="vertical", legend_top="45%", legend_pos='50%',
                xaxis_name=u'{0}'.format(self.spc_target), mark_line=["min", "max"], yaxis_name_gap=100,
                label_color=['rgba(255,106,106,0.5)'],
                yaxis_name=u'概率(Probability)')
        line = Line(width=960, height=1440)
        line.add(u'{0}'.format(self.spc_target).format(self.spc_target), x_line, y_line, xaxis_name=u'{0}'.format(self.spc_target),
                 yaxis_name=u'概率(Probability)',
                 line_color='rgba(0 ,255 ,127,0.5)', is_legend_show=False,
                 is_smooth=True, line_width=2,
                 is_datazoom_show=True, datazoom_type='both',
                 tooltip_tragger='axis',is_fill=True, area_color='#20B2AA', area_opacity=0.4)

        grid = Grid(width=1920, height=1440)
        grid.add(bar, grid_bottom="60%", grid_left="60%")
        grid.add(line, grid_bottom="60%", grid_right="60%")
        pyecharts.configure(force_js_embed=True)
        return grid.render_embed()

    def _get_scatter(self,  data):
        qty = len(data)
        x_line = np.arange(1, qty + 1, 1)
        y_line = data
        line = Line(width=1280, height=800)
        line.add(u'{0}'.format(self.spc_target), x_line, y_line, xaxis_name=u'Sequence',
                 yaxis_name=u'{0}'.format(self.spc_target), mark_line=["min", "average", "max"],
                 line_color='rgba(0 ,255 ,127,0.5)', legend_pos='center',
                 is_smooth=True, line_width=2, is_more_utils=True,
                 tooltip_tragger='axis')
        pyecharts.configure(force_js_embed=True)
        return line.render_embed()

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        data = DataFrame()
        mean = 0.0
        std = 0.0
        result = super(TorAngSPCReport, self).read(fields, load=load)
        if 'normal_dist' in fields or 'weibull_dist' in fields or'scatter' in fields and load == '_classic_read':
            data, length = self._get_data()
            if data.empty:
                self.env.user.notify_warning(u'查询获取结果:{0},请重新定义查询参数或等待新结果数据'.format(length))
                return result
            mean = np.mean(data)
            std = np.std(data)
        if 'normal_dist' in fields and not data.empty:
            result[0].update({'normal_dist': self._get_normal_dist(mean=mean,std=std)})
        if 'weibull_dist' in fields and not data.empty:
            scale_parameter = self.env['ir.config_parameter'].sudo().get_param('weibull.scale', default=1.0)
            shape_parameter = self.env['ir.config_parameter'].sudo().get_param('weibull.shape', default=5.0)
            result[0].update({'weibull_dist': self._get_weibull_dist(len(data), mean=mean, std=std,
                                                                     scale=scale_parameter, shape=shape_parameter)})
        if 'scatter' in fields and not data.empty:
            result[0].update({'scatter': self._get_scatter(data)})

        return result

    def _get_data(self):
        domain = [('measure_result', 'in', ['ok', 'nok'])]
        fields = ['measure_torque']
        order = 'measure_torque desc'
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
            order = 'measure_degree desc'
        _data = self.env['operation.result'].sudo().search_read(domain=domain, fields=fields, limit=self.limit)
        length = len(_data)
        if length < self.limit and length < 100:
            return DataFrame(),length
        df = DataFrame.from_dict(_data)
        df = df['measure_degree'] if self.spc_target == 'angle' else df['measure_torque']
        return df, length

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