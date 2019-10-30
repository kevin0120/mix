# -*- coding: utf-8 -*-
from odoo import api, SUPERUSER_ID
from odoo.http import request, Response, Controller, route
from odoo.exceptions import ValidationError
from odoo.addons.common_sa_utils.http import sa_http_session, sa_success_resp, sa_fail_response

from odoo.addons.spc.models.push_workorder import MASTER_WROKORDERS_API
import logging
import pprint
import json

schema = "http"

ORDER_REQUIRED_FIELDS = ['code', 'product_code', 'track_code', 'workcenter', 'worksheet', 'product', 'operation',
                         'steps']

HAVE_SOME_REQUIRED_FIELDS = {
    'worksheet': ['url', 'name', 'revision'],
    'environments': ['text', 'test_type', 'code', 'sequence'],
    'steps': ['code', 'test_type', 'sequence', 'failure_msg'],
    'operation': ['code', 'resources'],
    'product': ['url', 'code'],
    'components': ['code', 'is_key'],
}

_logger = logging.getLogger(__name__)


def validate_ts002_order_req_vals(vals):
    for field in ORDER_REQUIRED_FIELDS:
        if field not in vals:
            _logger.error('Field: {0} not in vals: {1}'.format(field, pprint.pformat(vals, indent=4)))
            raise ValidationError('Field: {0} Is Required, But Not in Request'.format(field))


def _validate_ts002_req_val_by_entry(entry, vals):
    if not vals:
        _logger.debug("Value Of Entry {0} Is Empty".format(entry))
        return
    if entry not in HAVE_SOME_REQUIRED_FIELDS.keys():
        return
    requireds = HAVE_SOME_REQUIRED_FIELDS.get(entry)
    for field in requireds:
        if isinstance(vals, dict):
            data = vals.get(field)
            if not data:
                raise ValidationError('Field: {0} Is Required In {1} , But Now Not in Request'.format(field, entry))
        if isinstance(vals, list):
            for d in vals:
                if not isinstance(d, dict):
                    continue
                keys = d.keys()
                if field not in keys:
                    raise ValidationError(
                        'Field: {0} Is Required In {1}, Value: {2}, But Now Not in Request'.format(
                            field, entry, pprint.pformat(d, indent=4)))


def validate_ts002_req_val_by_entry(vals):
    for entry, val in vals.items():
        _validate_ts002_req_val_by_entry(entry, val)


def convert_ts002_order(env, vals):
    ret = {}

    return ret


def get_masterpc_order_url(env, vals):
    workcenter_code = vals.get('workcenter')
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
    url = schema + '://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API)
    return workcenter_id, url


class TangcheMrpOrderGateway(Controller):
    @route('/api/v1/ts002/workorders', type='json', methods=['POST', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def tangcheOrderGateway(self, **kw):
        """
        唐车工单转换函数，将从MES接收的工单进行转换
        :param kw:
        :return:
        """
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        vals = request.jsonrequest
        try:
            validate_ts002_order_req_vals(vals)  # make sure field is all in request body

            validate_ts002_req_val_by_entry(vals)

            payload = convert_ts002_order(env, vals)
            _logger.debug("TS002 Get Order: {0}".format(pprint.pformat(payload)))

            session = sa_http_session()

            workcenter_id, master_url = get_masterpc_order_url(env, vals)
            resp = session.post(master_url, data=payload)
            if resp.status_code != 201:
                msg = 'TS002 Post WorkOrder To MasterPC Fail'
                return sa_fail_response(msg=msg)

            msg = "Tightening System Create Work Order Success"
            return sa_success_resp(status_code=201, msg=msg)
        except Exception as e:
            msg = e.message
            return sa_fail_response(msg=msg)
