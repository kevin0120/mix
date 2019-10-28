# -*- coding: utf-8 -*-
from odoo import api, SUPERUSER_ID
from odoo.http import request, Response, Controller, route
from odoo.exceptions import ValidationError
from odoo.addons.common_sa_utils.http import sa_http_session, sa_success_resp, sa_fail_response

from odoo.addons.spc.models.push_workorder import MASTER_WROKORDERS_API
import logging
import pprint
import json

ORDER_REQUIRED_FIELDS = []

_logger = logging.getLogger(__name__)


def validate_tangche_order_req_vals(vals):
    for field in ORDER_REQUIRED_FIELDS:
        if field not in vals:
            _logger.error('Field: {0} not in vals: {1}'.format(field, pprint.pformat(vals, indent=4)))
            raise ValidationError('Field: {0} Is Required, But Not in Request'.format(field))


def convert_tangche_order(env, vals):
    ret = {}

    return ret


def get_masterpc_order_url(env, vals):
    workcenter_code = vals.get('code')
    if not workcenter_code:
        raise ValidationError(
            'Can Not Found Work Center:{0} From The WorkOrder Request Values!'.format(workcenter_code))
    workcenter_id = env['mrp.workcenter'].search([('code', '=', workcenter_code)], limit=1)
    if not workcenter_id:
        raise ValidationError('Can Not Found Work Center:{0} From The WorkOrder!'.format(workcenter_code))
    master_pc = env['maintenance.equipment'].search('workcenter_id', '=', workcenter_id, limit=1)
    if not master_pc:
        raise ValidationError('Can Not Found Work Center:{0} From The WorkOrder!'.format(workcenter_code))
    connections = master_pc.connection_ids.filtered(
        lambda r: r.protocol == 'http') if master_pc.connection_ids else None
    if not connections:
        raise ValidationError('Can Not Found Connection Info For Work Center:{0}!'.format(workcenter_code))
    connect = connections[0]
    url = 'http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API)
    return url


class TangcheMrpOrderGateway(Controller):
    @route('/ts002/workorders', type='json', methods=['POST', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def tangcheOrderGateway(self, **kw):
        """
        唐车工单转换函数，将从MES接收的工单进行转换
        :param kw:
        :return:
        """
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        vals = request.jsonrequest
        try:
            validate_tangche_order_req_vals(vals)  # make sure field is all in request body

            payload = convert_tangche_order(env, vals)
            _logger.debug("TS002 Get Order: {0}".format(pprint.pformat(payload)))

            session = sa_http_session()

            master_url = get_masterpc_order_url(env, vals)
            resp = session.post(master_url, data=payload)
            if resp.status_code != 201:
                msg = 'TS002 Post WorkOrder To MasterPC Fail'
                return sa_fail_response(msg=msg)

            msg = "Tightening System Create Work Order Success"
            return sa_success_resp(status_code=201, msg=msg)

        except Exception as e:
            body = json.dumps({'msg': e.message})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
