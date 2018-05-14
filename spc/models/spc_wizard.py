# -*- coding: utf-8 -*-

from odoo import tools
from odoo import models, fields, api
from pyecharts import Line, Bar
import pyecharts


class OperationResultReport(models.TransientModel):
    _name = "operation.result.wizard"
    _description = "Result Statistics Process Control"

    query_date_from = fields.Datetime(string='Query Date From')
    query_date_to = fields.Datetime(string='Query Date to', default=fields.Datetime.now())
    vehicle_id = fields.Many2one('product.product', string='Vehicle Type', domain=[('sa_type', '=', 'vehicle')])
    screw_id = fields.Many2one('product.product', string='Screw Type', domain=[('sa_type', '=', 'screw')])
    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line')
    segment_id = fields.Many2one('mrp.worksegament', string='Work Segment')
    knr_code = fields.Char(string='KNR')
    vin_code = fields.Char(string='VIN')
    limit = fields.Integer('Query Limit', default=80)

    wave = fields.Text(string='Waves', store=False)

    def _get_echart_data(self):
        bar = Bar("我的第一个图表", "这里是副标题")
        bar.use_theme('dark')
        bar.add("服装", ["衬衫", "羊毛衫", "雪纺衫", "裤子", "高跟鞋", "袜子"], [5, 20, 36, 10, 75, 90])
        pyecharts.configure(force_js_embed=True)
        return bar.render_embed()

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        result = super(OperationResultReport, self).read(fields, load=load)
        if 'wave' in fields and load == '_classic_read':
            result[0].update({'wave': self._get_echart_data()})
        return result

    def _get_data(self):
        domain = []
        if self.query_date_from:
            domain += [('control_date', '>=', self.query_date_from)]
        if self.query_date_to:
            domain += [('control_date', '<=', self.query_date_to)]
        if self.vehicle_id:
            domain += [('product_id', '=', self.vehicle_id.id)]
        if self.screw_id:
            domain += [('consu_product_id', '=', self.vehicle_id.id)]
        ret = self.env['operation.result'].sudo().search(domain, limit=self.limit)

    @api.multi
    def button_query(self):
        pass
