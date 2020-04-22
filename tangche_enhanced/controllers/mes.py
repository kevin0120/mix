# -*- coding: utf-8 -*-
from odoo import http, SUPERUSER_ID, api
from os import getenv
from odoo.http import request, Response
import requests as Requests
from api_data import api_data
import json

MES_BASE_URL = getenv('ENV_SA_TS002_MES_BASE_URL', 'http://127.0.0.1:8080/api/tcmes')


def get_mes_endpoint(endpoint):
    return '{}{}'.format(MES_BASE_URL, endpoint)


MES_ENDPOINT_MAP = {
    'user_sync': {'url': get_mes_endpoint('/user'), 'method': Requests.get},  # 用户信息同步
    'device_sync': {'url': get_mes_endpoint('/device'), 'method': Requests.get},
    'productline_sync': {'url': get_mes_endpoint('/productline'), 'method': Requests.get},
    'worksegment_sync': {'url': get_mes_endpoint('/station'), 'method': Requests.get},
    'workcenter_sync': {'url': get_mes_endpoint('/workbench'), 'method': Requests.get},
}


class MESAPI(http.Controller):
    pass
