# -*- coding: utf-8 -*-
from odoo import http, SUPERUSER_ID, api
from os import getenv
from odoo.http import request, Response
import requests as Requests
from api_data import api_data
import json

MES_BASE_URL = getenv('MES_BASE_URL', 'http://172.26.214.80:8080/api/tcmes')


def get_mes_endpoint(endpoint):
    return '{}{}'.format(MES_BASE_URL, endpoint)


MES_ENDPOINT_MAP = {
    'user_sync': {'url': get_mes_endpoint('/user'), 'method': Requests.post},  # 用户信息同步
    'device_sync': {'url': get_mes_endpoint('/device'), 'method': Requests.post},
}


class MESAPI(http.Controller):
    pass
