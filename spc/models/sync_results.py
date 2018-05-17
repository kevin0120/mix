# -*- coding: utf-8 -*-

import datetime
import logging
import werkzeug.urls
import requests
from dateutil import parser

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.models import AbstractModel
from odoo.tools.translate import _
from odoo.tools import config
from odoo.tools import misc

import json
from dateutil.relativedelta import relativedelta

_logger = logging.getLogger(__name__)


MASTER_RESULT_API = '/api/v1/results'



class ResultSync(AbstractModel):
    _name = "result.sync"

    def _get_masterpc_results(self, url):
        headers = {'Content-Type': 'application/json'}
        payloads = {'has_upload': False, 'result': ['ok', 'nok']}
        ret = requests.get(url=url, params=payloads, headers=headers)
        if ret.status_code != 200:
            _logger.debug('Sync Result fail')
            return
        ret = ret.json()
        for result in ret:
            op_result = self.env['operation.result'].sudo().browse(result['id'])
            if not op_result:
                _logger.debug('Sync Result can not found result id: %d' % result['id'])
                continue
            rid = result.pop('id')
            if 'cur_objects' in result:
                result.update({
                    'cur_objects': json.dumps(result['cur_objects'])
                })

            if 'control_date' in result:
                _t = parser.parse(result['control_date']) if result['control_date'] else None
                if _t:
                    result.update({
                        'control_date': fields.Datetime.to_string((_t - _t.utcoffset()))
                    })
            ret = op_result.sudo().write(result)
            if not ret:
                _logger.debug(u'更新结果 写入结果失败 result id: %d' % result['id'])
                continue
            data = {'has_upload': True}
            ret = requests.patch(url=url + '/{0}'.format(rid), data=json.dumps(data), headers=headers)
            if ret.status_code != 200:
                _logger.debug(u'更新MasterPC hasupload标志位失败')
            return True

    @api.multi
    def result_sync(self):
        domain = [('protocol', '=', 'http'), ('equipment_id.category_name', '=', 'MasterPC')]
        connections = self.env['maintenance.equipment.connection'].sudo().search(domain)
        urls = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_RESULT_API) for connect in connections]
        for url in urls:
            self._get_masterpc_results(url)


