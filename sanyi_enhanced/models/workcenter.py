# -*- coding: utf-8 -*-

import logging
from odoo import models, fields, api
from odoo.exceptions import ValidationError
import requests as Requests
from requests import ConnectionError, RequestException
import json
from odoo.addons.sa_base.models.mrp_worksegment import DELETE_ALL_MASTER_WROKORDERS_API
from odoo.addons.sa_base.models.mrp_bom import MASTER_ROUTING_API

WORKCENTERS_BASE_INFO_API = '/rush/v1/mrp.workcenter/base_info'

_logger = logging.getLogger(__name__)


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    @api.one
    def _delete_workcenter_all(self, url):
        try:
            ret = Requests.delete(url, headers={'Content-Type': 'application/json'}, timeout=5)
            if ret.status_code == 200:
                self.env.user.notify_info(u'删除工位信息成功')
                return True
        except ConnectionError as e:
            self.env.user.notify_warning(u'删除工位信息失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'删除工位信息失败, 错误原因:{0}'.format(str(e)))
        except RequestException as e:
            self.env.user.notify_warning(u'删除工位信息失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'删除工位信息失败, 错误原因:{0}'.format(str(e)))
        return False

    @api.one
    def package_equipments_info(self):
        ret = []
        tools = self.equipment_ids.filtered(lambda tool: tool.category_name in ['tightening_wrench', 'tightening_gun'])
        for tool in tools:
            connection = tool.connection_ids.filtered(
                lambda r: r.protocol == 'rawtcp') if tool.connection_ids else None
            conn = connection[0] if connection else None
            user = tool.technician_user_id
            val = {
                'serial_no': tool.serial_no,
                'connection': 'tcp://{0}:{1}'.format(conn.ip, conn.port) if conn else None,
                'technician_user': {'name': user.login, 'uuid': user.uuid} if user else None,
            }
            ret.append(val)
        return ret

    @api.one
    def _get_workcenter_data(self):

        tools = self.package_equipments_info()
        data = {
            'name': self.name,
            'code': self.code,
            'tools': tools[0] if len(tools) > 0 else None
        }
        return data

    @api.one
    def _do_sync_workcenter(self, url):
        wc = self._get_workcenter_data()
        data = wc[0] if len(wc) > 0 else None
        try:
            ret = Requests.post(url, headers={'Content-Type': 'application/json'}, data=json.dumps(data), timeout=2)
            if ret.status_code == 200:
                self.env.user.notify_info(u'同步工位信息成功')
                return True
        except ConnectionError as e:
            self.env.user.notify_warning(u'同步工位信息失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'同步工位信息失败, 错误原因:{0}'.format(str(e)))
        except RequestException as e:
            self.env.user.notify_warning(u'同步工位信息失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'同步工位信息失败, 错误原因:{0}'.format(str(e)))
        return False

    @api.multi
    def button_sync_workcenter(self):
        try:
            for center in self:
                connects = center.get_workcenter_masterpc_http_connect()
                if not len(connects):
                    _logger.error("Can Not Found MasterPC Connect Info For Work Center:{0}".format(center.name))
                    continue
                connect = connects[0]
                if not connect:
                    _logger.error(
                        "Sync Operation Fail. Can Not Found Connect Info For WorkCenter: {0}".format(center.name))
                    continue
                delete_all_endpoint = 'http://{0}:{1}{2}'.format(connect.ip, connect.port,
                                                                 WORKCENTERS_BASE_INFO_API)
                center._delete_workcenter_all(delete_all_endpoint)
                center._do_sync_workcenter(delete_all_endpoint)

            return True
        except Exception as e:
            _logger.error("button_sync_operations Error", e)
            raise e

    @api.multi
    def button_sync_operations(self):
        """
        只同步拧紧相关的工步信息到工位控制器，为了后续的方便使用
        """
        try:
            operation_obj_sudo = self.env['mrp.routing.workcenter'].sudo()
            for center in self:
                connects = center.get_workcenter_masterpc_http_connect()
                if not len(connects):
                    _logger.error("Can Not Found MasterPC Connect Info For Work Center:{0}".format(center.name))
                    continue
                connect = connects[0]
                if not connect:
                    _logger.error(
                        "Sync Operation Fail. Can Not Found Connect Info For WorkCenter: {0}".format(center.name))
                    continue
                # delete_all_endpoint = 'http://{0}:{1}{2}'.format(connect.ip, connect.port,
                #                                                  DELETE_ALL_MASTER_WROKORDERS_API)
                #center._delete_workcenter_all_opertaions(delete_all_endpoint)
                if self.type == "rework":
                    # 返修工位
                    self.send_related_operations()
                else:
                    operations = operation_obj_sudo.search([('workcenter_id', '=', center.id)])
                    for operation in operations:
                        operation.button_send_mrp_routing_workcenter()
            return True
        except Exception as e:
            _logger.error("button_sync_operations Error", e)
            raise e

    @api.multi
    def send_related_operations(self):
        vals = []
        workcenters = self.env['mrp.workcenter'].search([('qc_workcenter_id', '=', self.id)])
        operation_obj_sudo = self.env['mrp.routing.workcenter'].sudo()
        for wc in workcenters:
            operations = operation_obj_sudo.search([('workcenter_id', '=', wc.id)])
            for operation in operations:
                op = operation.package_operaions()
                val = op[0] if len(op) > 0 else None
                if not val is None:
                    vals.extend(val)

        connects = self.get_workcenter_masterpc_http_connect()
        if not len(connects):
            _logger.error(
                "Can Not Found MasterPC Connect Info For Work Center:{0}".format(self.name))
            return
        connect = connects[0]
        url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_ROUTING_API)
        try:
            ret = Requests.put(url, data=json.dumps(vals), headers={'Content-Type': 'application/json'},
                               timeout=3)
            if ret.status_code == 200:
                self.env.user.notify_info(u'下发工艺成功')
        except ConnectionError as e:
            self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))
        except RequestException as e:
            self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))