# -*- coding: utf-8 -*-

import datetime
import logging
import werkzeug.urls
import requests
import itertools
from ast import literal_eval

from odoo import api, release, SUPERUSER_ID, fields
from odoo.exceptions import UserError
from odoo.models import AbstractModel
from odoo.tools.translate import _
from odoo.tools import config
from odoo.tools import misc

from datetime import datetime, date
from dateutil.relativedelta import relativedelta

_logger = logging.getLogger(__name__)


MATERPC_RESULT_API = '/api/v1/results'


def _get_masterpc_results(url):
    payloads = {'has_upload': False, 'result': ['ok', 'nok']}
    return requests.get(url=url, params=payloads).json()


class ResultSync(AbstractModel):
    _name = "result.sync"

    @api.multi
    def result_sync(self):
        domain = [('protocol', '=', 'http'),('equipment_id.category_name', '=', 'MasterPC')]
        connections = self.env['maintenance.equipment.connection'].sudo().search(domain)
        urls = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MATERPC_RESULT_API) for connect in connections]
        _ret = map(_get_masterpc_results, urls)
        ret = list(itertools.chain.from_iterable(_ret)) # 平展整个list
        for result in ret:
            pass


