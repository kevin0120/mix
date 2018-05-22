# -*- coding: utf-8 -*-

import logging

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.models import AbstractModel
import requests as Requests
from requests import ConnectionError, RequestException
import json

_logger = logging.getLogger(__name__)


AIIS_RESULT_API = '/api/v1/results'


class PushResult(AbstractModel):
    _name = "result.push"

    def _push_results(self, aiis_urls, results):
        for url in aiis_urls:
            for result in results:
                data = {
                    'agna': result.production_id.equipment_name,
                    'fname': result.production_id.factory_name,
                    'year': result.production_id.year,
                    'pin': result.production_id.pin,
                    'pin_check': result.production_id.pin_check_code,
                    'assembly_line': result.production_id.assembly_line_id.code,
                    'lnr': result.production_id.lnr,
                    'nut_no': result.consu_product_id.screw_type_code,
                    'date': result.control_date,
                    'result': result.measure_result.upper(),
                    'MI': result.measure_torque,
                    'WI': result.measure_degree
                }
                try:
                    ret = Requests.put(url, data=json.dumps(data), headers={'Content-Type': 'application/json'})
                    if ret.status_code == 200:
                        result.write({'sent': True})  ### 更新发送结果
                except ConnectionError:
                    break  # 此url的结果都不需要上传
                except RequestException as e:
                    print(e)
        return True

    @api.multi
    def result_push(self):
        domain = [('sent', '=', False), ('measure_result', 'in', ['ok', 'nok'])]
        results = self.env['operation.result'].sudo().search(domain)
        if not results:
            return True
        _aiis_urls = self.env['ir.config_parameter'].sudo().get_param('aiis.urls')
        if not _aiis_urls:
            return True
        aiis_urls = _aiis_urls.split(',')
        ret = map(lambda url: self._push_results(url, results), aiis_urls)
        return True


