# -*- coding: utf-8 -*-

from odoo import models, api
from odoo.addons.tangche_enhanced.controllers.mes import MES_ENDPOINT_MAP
from functools import wraps
import logging
import pprint
from odoo.tools import exception_to_unicode
from odoo.exceptions import UserError

_logger = logging.getLogger(__name__)


def mes_method_wrapper(key):
    endpoint = MES_ENDPOINT_MAP.get(key, None)
    if not endpoint:
        msg = u'{}为定义此方法'.format(key)
        _logger.error(msg)
        raise UserError(msg)

    def method_wrapper(f):
        @wraps(f)
        def wrapper(*args, **kwargs):
            try:
                data = f(*args, **kwargs)
                url = endpoint.get('url')
                method = endpoint.get('method')
                if not url or not method or not data:
                    _logger.error(u"数据缺失")
                    return
                headers = {'Content-Type': 'application/json'}
                resp = method(url, json=data, headers=headers)
                status_code = getattr(resp, 'status_code')
                rd = resp.json()
                _logger.debug("status_code:{}, data:{}".format(status_code, pprint.pformat(rd, indent=4)))
                if status_code > 400:
                    raise UserError('{}'.format(resp.content))
                if rd.get('Status', 'E') in ['E', 'e']:
                    raise UserError('{}'.format(rd.get('Message', '')))
            except Exception as e:
                _logger.error("{} Error: {}".format(f.__name__, exception_to_unicode(e)))
                raise e
            return resp

        return wrapper

    return method_wrapper


@mes_method_wrapper('user_sync')
def dosync_user_info(station_code):
    return {'station': station_code}


@mes_method_wrapper('device_sync')
def dosync_device_info(station_code):
    return {'station': station_code}


class MESMETHODAPI(models.AbstractModel):
    _name = 'mes.api_mixin'

    @api.multi
    def sync_user_info(self, station_code):
        try:
            ret = dosync_user_info(station_code)
            # todo: 将数据进行处理
        except Exception as e:
            raise e

    @api.multi
    def sync_device_info(self, station_code):
        try:
            ret = dosync_device_info(station_code)
            # todo: 将数据进行处理
        except Exception as e:
            raise e
