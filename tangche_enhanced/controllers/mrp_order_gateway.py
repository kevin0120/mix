# -*- coding: utf-8 -*-
from odoo import api, SUPERUSER_ID
from odoo.http import request, Response, Controller, route
from odoo.exceptions import ValidationError
from odoo.addons.common_sa_utils.http import sa_success_resp, sa_fail_response
import os
from odoo.addons.spc.models.push_workorder import MASTER_WROKORDERS_API
import logging
import pprint
import itertools
import requests
from tenacity import retry, wait_exponential, retry_if_exception_type, stop_after_delay, RetryError
import json

MES_BASE_URL = os.environ.get('ENV_SA_TS002_MES_BASE_URL', 'http://127.0.0.1:1888')

schema = "http"

ORDER_REQUIRED_FIELDS = ['code', 'product_code', 'track_code', 'workcenter', 'worksheet', 'product', 'operation',
                         'steps']

headers = {'Content-Type': 'application/json'}

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


def package_tightening_points(tightening_points):
    ret = []
    for tp in tightening_points:
        pset = tp.program_id.code
        if isinstance(tp.program_id.code, str):
            pset = int(tp.program_id.code)
        tool_id = tp.tightening_tool_ids[0].tool_id if tp.tightening_tool_ids else None
        if not tool_id:
            _logger.error('Can Not Found Tool:{0}'.format(tp.name or tp.code))
        val = {
            'tightening_tool': tool_id.serial_no if tool_id else False,
            'x': tp.x_offset,
            'y': tp.y_offset,
            'pset': pset,
            'sequence': tp.sequence,
            'group_sequence': tp.group_sequence if tp.group_id else tp.sequence,
            'is_key': tp.is_key,  # 是否为关键拧紧点
            'max_redo_times': tp.max_redo_times,
            'key_num': tp.group_id.key_num if tp.group_id else 1,
        }
        ret.append(val)
    return ret


def convert_ts002_order(env, vals):
    ret = vals
    mes_work_steps = vals.get('steps')
    if not mes_work_steps:
        raise ValidationError("Can Not Get Work Step From External System")
    work_steps = list(itertools.chain.from_iterable(mes_work_steps))  # flat array make sure
    vals['steps'] = work_steps
    tightening_steps = filter(lambda step: step.get('test_type') == 'tightening', work_steps)
    if not tightening_steps:
        return vals
    for ts in tightening_steps:
        tc = ts.get('code')
        ws = env['sa.quality.point'].search(['|', ('ref', '=', tc), ('name', '=', tc)])
        if not ws:
            _logger.error("Can Not Found Tightening Step By Code:{0}".format(tc))
            continue
        rws = ws
        if len(ws) != 1:
            _logger.error("Tightening Step By Code:{0} Is Not Unique".format(tc))
            rws = ws[0]
        expect_tightening_total = ts.get('tightening_total', 0)
        if not expect_tightening_total:
            raise ValidationError('Can Not Found Tightening Total Within The Tightening Step: {0}'.format(tc))
        if len(rws.operation_point_ids) != expect_tightening_total:
            _logger.error(
                "Tightening Count Is Not Equal Within Tightening System By Code:{0}. "
                "Except:{1}, Real:{2}".format(tc, expect_tightening_total, len(rws.operation_point_ids)))

        if not rws.worksheet_img:
            raise ValidationError('Can Not Found Tightening Image Within The Tightening Step: {0}'.format(tc))
        ts.update({'tightening_image_by_step_code': tc})
        val = package_tightening_points(rws.operation_point_ids)
        ts.update({'tightening_points': val})  # 将拧紧点的包包裹进去
    return ret


def package_workcenter_location_data(workcenter_id, val):
    ret = []
    entry = val['workcenter']
    if not isinstance(entry, dict):
        entry = {'code': workcenter_id.code}
    else:
        entry.update({'code': workcenter_id.code})
    workcenter_locations = workcenter_id.sa_workcenter_loc_ids
    if not workcenter_locations:
        _logger.error('Can Not Found Any Location In Work Center')
        return
    for loc in workcenter_locations:
        data = {
            'product_code': loc.product_id.default_code if loc.product_id else False,
            'io_output': loc.io_output,
            'io_input': loc.io_input
        }
        ret.append(data)
    entry.update({'locations': ret})


def get_masterpc_order_url_and_package(env, vals):
    workcenter_code = vals.get('workcenter')
    if not workcenter_code:
        raise ValidationError(
            'Can Not Found Work Center:{0} From The WorkOrder Request Values!'.format(workcenter_code))
    workcenter_id = env['mrp.workcenter'].search([('code', '=', workcenter_code)], limit=1)
    if not workcenter_id:
        raise ValidationError('Can Not Found Work Center:{0} From The WorkOrder!'.format(workcenter_code))
    master_pc = env['maintenance.equipment'].search(
        [('workcenter_id', '=', workcenter_id.id), ('category_name', '=', 'MasterPC')], limit=1)
    if not master_pc:
        raise ValidationError('Can Not Found Work Center:{0} From The WorkOrder!'.format(workcenter_code))
    connections = master_pc.connection_ids.filtered(
        lambda r: r.protocol == 'http') if master_pc.connection_ids else None
    if not connections:
        raise ValidationError('Can Not Found Connection Info For Work Center:{0}!'.format(workcenter_code))
    connect = connections[0]
    url = schema + '://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API)
    return workcenter_id, url


@retry(wait=wait_exponential(multiplier=1, min=1, max=3), stop=stop_after_delay(5), retry=retry_if_exception_type())
def post_order_2_masterpc(master_url, data):
    resp = requests.post(master_url, data=json.dumps(data), headers=headers)
    if resp.status_code != 201:
        msg = 'TS002 Post WorkOrder To MasterPC Fail'
        return sa_fail_response(msg=msg)
    msg = "Tightening System Create Work Order Success"
    return sa_success_resp(status_code=201, msg=msg)


@retry(wait=wait_exponential(multiplier=1, min=1, max=3), stop=stop_after_delay(5), retry=retry_if_exception_type())
def query_order_from_mes(order_code, workcenter_code):
    if not order_code:
        raise ValidationError('query_order_from_mes Query WorkOrder Must Include Order Code')
    url = '{0}/workorders'.format(MES_BASE_URL)
    params = {'code': order_code}
    if workcenter_code:
        params.update({'workcenter': workcenter_code})
    resp = requests.get(url, params=params)
    if resp.status_code != 200:
        msg = 'TS002 Query Work Order:{0} From MES Fail'.format(order_code)
        return sa_fail_response(msg=msg, extra=resp.json())
    return resp.json()


def _convert_orders_info(env, values):
    validate_ts002_order_req_vals(values)  # make sure field is all in request body
    validate_ts002_req_val_by_entry(values)
    payload = convert_ts002_order(env, values)
    workcenter_id, master_url = get_masterpc_order_url_and_package(env, payload)
    package_workcenter_location_data(workcenter_id, payload)
    _logger.debug("TS002 Get Order: {0}".format(pprint.pformat(payload)))

    resp = post_order_2_masterpc(master_url, payload)
    if resp:
        return resp

    msg = "Tightening System Create Work Order Success"
    return sa_success_resp(status_code=201, msg=msg)


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
            return _convert_orders_info(env, vals)
        except RetryError as err:
            msg = str(err)
            return sa_fail_response(msg=msg)
        except Exception as e:
            msg = str(e)
            return sa_fail_response(msg=msg)

    @route('/api/v1/ts002/workorders', type='json', methods=['GET', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def tangcheQueryOrderGateway(self, **kw):
        try:
            if 'code' not in kw:
                raise ValidationError('Query WorkOrder Must Include Order Code')
            resp = query_order_from_mes(kw.get('code'), kw.get('workcenter'))
            if isinstance(resp, Response):
                # 返回值如果是response，代表异常，否则返回的是字典对象
                return resp
            env = api.Environment(request.cr, SUPERUSER_ID, request.context)
            vals = resp
            return _convert_orders_info(env, vals)

        except RetryError as err:
            msg = str(err)
            return sa_fail_response(msg=msg)
        except Exception as e:
            msg = str(e)
            return sa_fail_response(msg=msg)
