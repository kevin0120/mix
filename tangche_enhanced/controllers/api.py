# -*- coding: utf-8 -*-
from odoo import http, SUPERUSER_ID, api
from odoo.http import request, Response
from api_data import api_data
from odoo.addons.api.controllers.controllers import BaseApi
import json


class ts002_api_swagger(BaseApi):
    @http.route()
    def _api_doc(self):
        return json.dumps(api_data)
