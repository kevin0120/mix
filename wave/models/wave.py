# -*- coding: utf-8 -*-

from odoo import models, fields, api
from pyecharts import Bar
import pyecharts
from minio import Minio
from urlparse import urlsplit
import json

minioClient = None


class Wave(models.TransientModel):
    '''
        瞬态模型. 不需要持久化
    '''
    _name = 'wave.wave'
    _description = 'Waveform scope wizard'

    # . In its computed method do::

    query_date_from = fields.Datetime(string='Query Date From')
    query_date_to = fields.Datetime(string='Query Date to', default=fields.Datetime.now())
    vehicle_id = fields.Many2one('product.product', string='Vehicle Type', domain=[('sa_type','=','vehicle')])
    screw_id = fields.Many2one('product.product', string='Screw Type', domain=[('sa_type','=','screw')])
    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line')
    segment_id = fields.Many2one('mrp.worksegament', string='Work Segment')
    knr_code = fields.Char(string='KNR')
    vin_code = fields.Char(string='VIN')
    limit = fields.Integer('Query Limit', default=80)

    wave = fields.Text(string='Waves', store=False)

    def _recreate_minio_client(self):
        global minioClient
        minio_url = self.env['ir.config_parameter'].get_param('minio.url')
        minio_bucket = self.env["ir.config_parameter"].get_param("minio.bucket")
        minio_access_key = self.env["ir.config_parameter"].get_param("minio.access_key")
        minio_secret_key = self.env["ir.config_parameter"].get_param("minio.secret_key")
        secruity = False if urlsplit(minio_url).scheme == 'http' else True
        if not minioClient:
            minioClient = Minio(minio_url.split('://')[-1], access_key=minio_access_key, secret_key=minio_secret_key, secure=secruity)
        elif minioClient._endpoint_url != minio_url or minioClient._access_key != minio_access_key or minioClient._secret_key != minio_secret_key:
            minioClient = Minio(minio_url.split('://')[-1], access_key=minio_access_key, secret_key=minio_secret_key, secure=secruity)
        if not minioClient.bucket_exists(minio_bucket):
            minioClient.make_bucket(minio_bucket)
        return minioClient, minio_bucket

    def _get_echart_data(self,datas):
        bar = Bar("我的第一个图表", "这里是副标题")
        bar.add("服装", ["衬衫", "羊毛衫", "雪纺衫", "裤子", "高跟鞋", "袜子"], [5, 20, 36, 10, 75, 90])
        bar.print_echarts_options()
        pyecharts.configure(force_js_embed=True)
        return bar.render_embed()

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        result = super(Wave, self).read(fields, load=load)
        if 'wave' in fields and load == '_classic_read':
            datas = self._get_data()
            result[0].update({'wave': self._get_echart_data(datas)})
        return result

    def _get_data(self):
        client, bucket = self._recreate_minio_client()
        ret = []
        for i in range(2):
            stream_data = client.get_object(bucket, 'cur.json')
            data = json.loads(stream_data.data.decode('utf-8'))
        return data

    @api.multi
    def button_query(self):
        pass
