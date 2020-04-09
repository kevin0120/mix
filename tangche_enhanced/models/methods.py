# -*- coding: utf-8 -*-

from odoo import models, api
from odoo.addons.tangche_enhanced.controllers.mes import MES_ENDPOINT_MAP
from functools import wraps
import logging
import pprint
from odoo.tools import exception_to_unicode
from odoo.exceptions import UserError
from requests import Response
import urllib

_logger = logging.getLogger(__name__)

# 3秒 timeout
DEFAULT_TIMEOUT = 3


# 同步用户， 根据工号判断用户是否存在， 无则创建有则更新
def _mes_sync_user(env, value):
    if not env or not value:
        return False

    try:
        # TODO: 处理用户班次
        user_shift = value.get('group', '')

        # 根据用户工号定位用户
        user_id = env['res.users'].search([('uuid', '=', value.get('EmployeeNumber', ''))], limit=1)
        user_data = {
            'name': value.get('Name', ''),
            'login': value.get('Name', ''),
        }
        if not user_id:
            # 创建用户
            user_data['uuid'] = value.get('EmployeeNumber', '')
            env['res.users'].create(user_data)
        else:
            # 更新用户
            user_id.write(user_data)

        return True

    except Exception as e:
        _logger.error('_mes_sync_user Failed: {0}'.format(e.message))
        return False


# 同步设备， 根据设备序列号判断设备是否存在， 无则创建有则更新
def _mes_sync_device(env, value):
    if not env or not value:
        return False

    try:
        dev_sn = value.get('EquipmentNo', '')
        dev_ip = value.get('EquipmentIP', '')
        dev_type = value.get('EquipmentType', '')

        device_id = env['maintenance.equipment'].search([('serial_no', '=', dev_sn)], limit=1)
        device_type_id = env['maintenance.equipment.category'].search([('technical_name', '=', dev_type)],
                                                                      limit=1)
        if not device_type_id:
            raise Exception("Device Type Not Found: {}".format(value.get('EquipmentType', '')))

        device_data = {
            'name': value.get('EquipmentName', ''),
            'category_id': device_type_id.id,
        }

        parent_device_id = env['maintenance.equipment'].search([('serial_no', '=', value.get('ParentNo', ''))], limit=1)
        if parent_device_id:
            device_data['parent_id'] = parent_device_id.id

        workcenter_id = env['mrp.workcenter'].search([('code', '=', value.get('Location', ''))], limit=1)
        if workcenter_id:
            device_data['workcenter_id'] = workcenter_id.id

        if not device_id:
            # 创建设备
            device_data['serial_no'] = dev_sn
            device_id = env['maintenance.equipment'].create(device_data)
        else:
            # 更新设备
            device_id.write(device_data)

        if dev_ip is not str(''):
            conn_vals = str.split(str(dev_ip), ':')
            conn_data = {
                'protocol': 'modbustcp' if dev_type == 'IO' else 'rawtcp',
                'name': 'Connection',
                'equipment_id': device_id.id,
                'ip': conn_vals[0],
                'port': int(conn_vals[1])
            }

            # 更新设备ip
            conn_id = env['maintenance.equipment.connection'].search([('equipment_id', '=', device_id.id)], limit=1)
            if not conn_id:
                # 创建设备连接
                env['maintenance.equipment.connection'].create(conn_data)
            else:
                # 更新设别连接
                conn_id.write(conn_data)

    except Exception as e:
        _logger.error('_mes_sync_device Failed: {0}'.format(e.message))
        return False


# 同步产线, 根据产线编号判断产线是否存在， 无则创建有则更新
def _mes_sync_line(env, value):
    if not env or not value:
        return False

    try:
        line_code = value.get('ProLineCode', '')
        line_name = value.get('ProLineName', '')

        line_data = {
            'name': line_name
        }

        line_id = env['mrp.assemblyline'].search([('code', '=', line_code)], limit=1)
        if not line_id:
            # 创建产线
            line_data['code'] = line_code
            env['mrp.assemblyline'].create(line_data)
        else:
            # 更新产线
            line_id.write(line_data)

    except Exception as e:
        _logger.error('_mes_sync_line Failed: {0}'.format(e.message))
        return False


# 同步工段(唐车工位)，根据工段编号判断工段是否存在， 无则创建有则更新
def _mes_sync_worksection(env, value):
    if not env or not value:
        return False

    try:
        section_code = value.get('WorkstationCode', '')
        section_name = value.get('WorkstationName', '')
        section_data = {
            'name': section_name,
        }

        line_code = value.get('ProLineCode', '')
        line_id = env['mrp.assemblyline'].search([('code', '=', line_code)], limit=1)
        if line_id:
            section_data['workassembly_id'] = line_id.id

        section_id = env['mrp.worksection'].search([('code', '=', section_code)], limit=1)
        if not section_id:
            # 创建工段
            section_data['code'] = section_code
            env['mrp.worksection'].create(section_data)
        else:
            # 更新工段
            section_id.write(section_data)

    except Exception as e:
        _logger.error('_mes_sync_worksection Failed: {0}'.format(e.message))
        return False


# 同步工位(唐车台位)，根据工位编号判断工位是否存在， 无则创建有则更新
def _mes_sync_workcenter(env, value):
    if not env or not value:
        return False

    try:
        workcenter_code = value.get('WorkbenchCode', '')
        workcenter_name = value.get('WorkbenchName', '')
        workcenter_data = {
            'name': workcenter_name,
        }

        section_code = value.get('WorkstationCode', '')
        section_id = env['mrp.worksection'].search([('code', '=', section_code)], limit=1)
        if section_id:
            workcenter_data['worksegment_id'] = section_id.id

        workcenter_id = env['mrp.workcenter'].search([('code', '=', workcenter_code)], limit=1)
        if not workcenter_id:
            # 创建工位
            workcenter_data['code'] = workcenter_code
            env['mrp.workcenter'].create(workcenter_data)
        else:
            # 更新工位
            workcenter_id.write(workcenter_data)

    except Exception as e:
        _logger.error('_mes_sync_workcenter Failed: {0}'.format(e.message))
        return False


def mes_method_wrapper(key):
    endpoint = MES_ENDPOINT_MAP.get(key, None)
    if not endpoint:
        msg = u'{}为定义此方法'.format(key)
        _logger.error(msg)
        raise UserError(msg)

    def method_wrapper(f):
        @wraps(f)
        def wrapper(*args, **kwargs):
            resp = None
            try:
                data = f(*args, **kwargs)
                url = endpoint.get('url')
                method = endpoint.get('method')
                if not url or not method:
                    _logger.error(u"数据缺失")
                    return resp
                headers = {'Content-Type': 'application/json'}
                if data is not None:
                    url += '?{}'.format(urllib.urlencode(data))
                resp = method(url, headers=headers, timeout=DEFAULT_TIMEOUT)
                status_code = getattr(resp, 'status_code')
                if status_code > 400:
                    raise UserError('{}'.format(resp.content))
                rd = resp.json()
                _logger.debug("status_code:{}, data:{}".format(status_code, pprint.pformat(rd, indent=4)))
                if rd.get('Status', 'E') in ['E', 'e']:
                    raise UserError('{}'.format(rd.get('Message', '')))

                return rd['data']
            except Exception as e:
                _logger.error("{} Error: {}".format(f.__name__, exception_to_unicode(e)))
                raise e

        return wrapper

    return method_wrapper


@mes_method_wrapper('user_sync')
def dosync_user_info(station_code):
    if not station_code:
        return None
    return {'station': station_code}


@mes_method_wrapper('device_sync')
def dosync_device_info(station_code):
    if not station_code:
        return None
    return {'station': station_code}


@mes_method_wrapper('productline_sync')
def dosync_productline_info():
    return None


@mes_method_wrapper('worksegment_sync')
def dosync_worksegment(productline_code):
    if not productline_code:
        return None
    return {'productlineCode': productline_code}


@mes_method_wrapper('workcenter_sync')
def dosync_workcenter():
    return None


class MESMETHODAPI(models.AbstractModel):
    _name = 'mes.api_mixin'

    # 同步用户
    @api.multi
    def sync_user_info(self, station_code=None):
        try:
            ret = dosync_user_info(station_code)
            for value in ret:
                _mes_sync_user(self.env, value)
        except Exception as e:
            raise e

    # 同步设备
    @api.multi
    def sync_device_info(self, station_code=None):
        try:
            ret = dosync_device_info(station_code)
            for value in ret:
                _mes_sync_device(self.env, value)
        except Exception as e:
            raise e

    # 同步产线
    @api.multi
    def sync_productline_info(self):
        try:
            ret = dosync_productline_info()
            for value in ret:
                _mes_sync_line(self.env, value)
        except Exception as e:
            raise e

    # 同步工段（唐车工位）
    @api.multi
    def sync_worksegment_info(self, productline_code=None):
        try:
            ret = dosync_worksegment(productline_code)
            for value in ret:
                _mes_sync_worksection(self.env, value)
        except Exception as e:
            raise e

    # 同步工位（唐车台位）
    @api.multi
    def sync_workcenter_info(self):
        try:
            ret = dosync_workcenter()
            for value in ret:
                _mes_sync_workcenter(self.env, value)
        except Exception as e:
            raise e