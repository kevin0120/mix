# -*- coding: utf-8 -*-

from odoo import models, fields, api
import pyecharts
from pyecharts import Bar


class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    def _compute_bokeh_chart(self):
        for rec in self:
            bar = Bar("我的第一个图表", "这里是副标题")
            bar.add("服装", ["衬衫", "羊毛衫", "雪纺衫", "裤子", "高跟鞋", "袜子"], [5, 20, 36, 10, 75, 90])
            bar.print_echarts_options()
            pyecharts.configure(force_js_embed=True)
            x = bar.render_embed()
            rec.echart = x

    routing_id = fields.Many2one(
        'mrp.routing', string='Operation', related='operation_id.routing_id', readonly=1)

    worksheet_img = fields.Binary(
        'Worksheet', related='routing_id.worksheet_img', readonly=True)

    echart = fields.Text(
        string='EChart',
        compute=_compute_bokeh_chart)
