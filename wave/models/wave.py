# -*- coding: utf-8 -*-

from odoo import models, fields, api
from pyecharts import Bar
import pyecharts


class Wave(models.TransientModel):
    '''
        抽象模型. 不需要持久化
    '''
    _name = 'wave.wave'
    _description = 'Waveform scope wizard'

    # . In its computed method do::

    name = fields.Char()
    value = fields.Integer()
    description = fields.Text()

    echart = fields.Text(string='EChart', store=False)

    def _get_echart_data(self):
        bar = Bar("我的第一个图表", "这里是副标题")
        bar.add("服装", ["衬衫", "羊毛衫", "雪纺衫", "裤子", "高跟鞋", "袜子"], [5, 20, 36, 10, 75, 90])
        bar.print_echarts_options()
        pyecharts.configure(force_js_embed=True)
        return bar.render_embed()

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        result = super(Wave, self).read(fields, load=load)
        if 'echart' in fields and load == '_classic_read':
            result[0].update({'echart': self._get_echart_data()})
        return result

    @api.multi
    def button_query(self):
        pass
