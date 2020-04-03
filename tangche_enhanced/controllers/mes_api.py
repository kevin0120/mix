# -*- coding: utf-8 -*-
from odoo.models import AbstractModel
from odoo import api
import logging
import requests
import mrp_order_gateway

headers = {'Content-Type': 'application/json'}

# 3秒 timeout
DEFAULT_TIMEOUT = 3

_logger = logging.getLogger(__name__)


# 同步用户， 根据工号判断用户是否存在， 无则创建有则更新
def _mes_sync_user(env, value):
    if not env or not value:
        return False

    try:
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
        tool_sn = value.get('EquipmentNo', '')
        device_id = env['maintenance.equipment'].search([('serial_no', '=', tool_sn)], limit=1)
        device_type_id = env['maintenance.equipment.category'].search([('name', '=', value.get('EquipmentType', ''))],
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
            device_data['serial_no'] = tool_sn
            env['maintenance.equipment'].create(device_data)
        else:
            # 更新设备
            device_id.write(device_data)

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
        section_name = value.get('WorkstationlName', '')
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


class MesApi(AbstractModel):
    _name = "mes.api"

    @api.multi
    def sync_users(self):
        url = mrp_order_gateway.MES_BASE_URL + '/TSIF01'
        try:
            resp = requests.get(url=url, headers=headers, timeout=DEFAULT_TIMEOUT)
            if resp.status_code != 200:
                return

            body = resp.json()
            for value in body:
                _mes_sync_user(self.env, value)
        except Exception as e:
            _logger.error('MesApi.sync_users Failed: {0}'.format(e.message))


    @api.multi
    def sync_devices(self):
        url = mrp_order_gateway.MES_BASE_URL + '/TSIF02'
        try:
            resp = requests.get(url=url, headers=headers, timeout=DEFAULT_TIMEOUT)
            if resp.status_code != 200:
                return

            body = resp.json()
            for value in body:
                _mes_sync_device(self.env, value)
        except Exception as e:
            _logger.error('MesApi.sync_devices Failed: {0}'.format(e.message))


    @api.multi
    def sync_lines(self):
        url = mrp_order_gateway.MES_BASE_URL + '/TSIF03'
        try:
            resp = requests.get(url=url, headers=headers, timeout=DEFAULT_TIMEOUT)
            if resp.status_code != 200:
                return

            body = resp.json()
            for value in body:
                _mes_sync_line(self.env, value)
        except Exception as e:
            _logger.error('MesApi.sync_lines Failed: {0}'.format(e.message))


    @api.multi
    def sync_worksections(self):
        url = mrp_order_gateway.MES_BASE_URL + '/TSIF04'
        try:
            resp = requests.get(url=url, headers=headers, timeout=DEFAULT_TIMEOUT)
            if resp.status_code != 200:
                return

            body = resp.json()
            for value in body:
                _mes_sync_worksection(self.env, value)
        except Exception as e:
            _logger.error('MesApi.sync_worksections Failed: {0}'.format(e.message))


    @api.multi
    def sync_workcenters(self):
        url = mrp_order_gateway.MES_BASE_URL + '/TSIF05'
        try:
            resp = requests.get(url=url, headers=headers, timeout=DEFAULT_TIMEOUT)
            if resp.status_code != 200:
                return

            body = resp.json()
            for value in body:
                _mes_sync_workcenter(self.env, value)
        except Exception as e:
            _logger.error('MesApi.sync_workcenters Failed: {0}'.format(e.message))
