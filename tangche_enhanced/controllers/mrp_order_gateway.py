# -*- coding: utf-8 -*-
from odoo import api, SUPERUSER_ID, fields
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

# ORDER_REQUIRED_FIELDS = ['WIPORDERNO', 'PRODUCTNO']

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


def str_time_to_rfc3339(s_time):
    sp = s_time.split(' ')
    return sp[0] + 'T' + sp[1] + 'Z'


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
    for ltp in tightening_points:
        if ltp._name == 'mrp.wo.consu.line':
            tp = request.env['operation.point'].sudo().search([('qcp_id', '=', ltp.point_id.id)], limit=1)
        else:
            tp = ltp
        pset = tp.program_id.code
        if isinstance(tp.program_id.code, str):
            pset = int(tp.program_id.code)
        if ltp._name == 'mrp.wo.consu.line':
            tool_id = ltp.tool_id if ltp.tool_id else None
        else:
            tool_id = tp.tool_id if tp.tool_id else None
        if not tool_id:
            _logger.error('Can Not Found Tool:{0}'.format(tp.name or tp.code))
        val = {
            'tightening_tool': tool_id.serial_no if tool_id else False,
            'x': tp.x_offset,
            'y': tp.y_offset,
            'pset': pset,
            'sequence': tp.sequence,
            'group_sequence': tp.group_sequence,
            'is_key': tp.is_key,  # 是否为关键拧紧点
            'nut_no': tp.name,  # 螺栓编号
            'max_redo_times': tp.max_redo_times,
            'key_num': tp.group_id.key_num if tp.group_id else 1,
        }
        ret.append(val)
    return ret


#
# def package_finished_product(env, vals):
#     if 'code' not in vals:
#         raise ValidationError("Can Not Get Finished Product Code From External System")
#     code = vals.get('code')
#     product_id = env['product.product'].search(['|', ('default_code', '=', code), ('name', '=', code)])
#     if not product_id:
#         raise ValidationError("Can Not Get Finished Product By Code:{0}".format(code))

def pack_step_payload(env, consum_lines):
    payloads = []

    type_tightening_id = env.ref('quality.test_type_tightening').id
    type_tightening_point_id = env.ref('quality.test_type_tightening_point').id
    for idx, step in enumerate(consum_lines.filtered(lambda t: t.test_type_id.id != type_tightening_point_id)):
        ts = {
            "code": step.name,
            "desc": step.note or '',
            "failure_msg": step.failure_message or '',
            "sequence": idx + 1,
            "skippable": step.can_do_skip,
            "undoable": step.can_do_redo,
            "test_type": step.test_type_id.technical_name,
            "consume_product": step.component_id.code,
            "text": step.reason or '',
            "tolerance_min": step.tolerance_min,
            "tolerance_max": step.tolerance_max,
            "target": step.norm,
        }
        ts.update({'tightening_image_by_step_code': step.name or step.ref})
        if step.test_type_id.id == type_tightening_id:
            val = package_tightening_points(step.operation_point_ids)
            ts.update({'tightening_points': val})  # 将拧紧点的包包裹进去
            ts.update({'tightening_total': len(step.operation_point_ids)})  # 将拧紧点的包包裹进去
        payloads.append(ts)

    return payloads


def convert_ts002_order(env, vals):
    try:
        code = vals.get('requestInfo').get('MOMWIPORDER').get('WIPORDERNO')
        mom_productno = vals.get('requestInfo').get('MOMWIPORDER').get('PRODUCTNO')
        worksection = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('WORKCENTER')
        date_planned_start = vals.get('requestInfo').get('MOMWIPORDER').get('SCHEDULEDSTARTDATE')
        date_planned_complete = vals.get('requestInfo').get('MOMWIPORDER').get('SCHEDULEDCOMPLETIONDATE')
        worksheet = {
            "name": vals.get('requestInfo').get('MOMWIPORDER').get('WIDOCS').get('WIDOC').get('DESCRIPT'),
            "revision": vals.get('requestInfo').get('MOMWIPORDER').get('WIDOCS').get('WIDOC').get('DOCVR'),
            "url": vals.get('requestInfo').get('MOMWIPORDER').get('WIDOCS').get('WIDOC').get('DOCURL'),
        }
        modeldoc = vals.get('requestInfo').get('MOMWIPORDER').get('MODELDOCS').get('MODELDOC')
        products = list()
        for mo in modeldoc:
            product = {
                'code': mo.get('PRODUCTNO'),
                'url': mo.get('URLLOCATION'),
            }
            products.append(product)

        SYSTEMTYPE = vals.get('requestInfo').get('SYSTEMTYPE')
        WIPORDERTYPE = vals.get('requestInfo').get('MOMWIPORDER').get('WIPORDERTYPE')
        LOCATION = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('LOCATION')
        SKILL = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('SKILLS').get('SKILL')


        OPRSEQUENCENO = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('OPRSEQUENCENO')
        MOMDISPOSITIONS = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('MOMDISPOSITIONS').get(
            'MOMDISPOSITION')
        MOMCONFIG = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('MOMCONFIG')
        RESOURCEGROUP = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('RESOURCEGROUP')
        STARTEMPLOYEE = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('STARTEMPLOYEE')
        RESOURCENAMES = vals.get('requestInfo').get('MOMWIPORDER').get('MOMWIPORDEROPR').get('RESOURCENAMES')
        PARENTWIPORDERNO = vals.get('requestInfo').get('MOMWIPORDER').get('PARENTWIPORDERNO')
        PARENTWIPORDERTYPE = vals.get('requestInfo').get('MOMWIPORDER').get('PARENTWIPORDERTYPE')
    except RetryError:
        msg = 'TS002  WorkOrder Payload is not qualifed!!'
        return [], msg
    except Exception:
        msg = 'TS002  WorkOrder Payload is not qualifed!!'
        return [], msg

    ret = list()
    ws = env['mrp.worksection'].search([('code', '=', worksection)])
    for ts in ws.workcenter_ids:
        pro = env['product.product'].search([('default_code', '=', mom_productno)])
        bom = env['mrp.bom'].search([('product_id', '=', pro.id)])

        # mrw1 = env['mrp.routing.workcenter'].search(
        #     [('workcenter_id.id', '=', ts.id)]).filtered(lambda r: bom.routing_id in r.sa_routing_ids)

        mrw = bom.routing_id.sa_operation_ids.filtered(lambda r: r.workcenter_id.id == ts.id)
        _steps = pack_step_payload(env, mrw.sa_step_ids)
        vals = {
            'code': code,
            'track_no': mom_productno,
            'product_code': mom_productno,
            'workcenter': ts.code,
            'date_planned_start': str_time_to_rfc3339(
                date_planned_start) if date_planned_start else str_time_to_rfc3339(fields.Datetime.now()),
            'date_planned_complete': str_time_to_rfc3339(
                date_planned_complete) if date_planned_complete else str_time_to_rfc3339(fields.Datetime.now()),
            'worksheet': worksheet,
            'products': products,
            'operation': {
                'code': mrw.name or mrw.ref,
                'desc': '',
                "resources": {
                    "users": [],
                    "equipments": [],
                },
            },
            'components': [],
            'environments': [],
            "SYSTEMTYPE": SYSTEMTYPE,
            "WIPORDERTYPE": WIPORDERTYPE,
            "MOMDISPOSITIONS": MOMDISPOSITIONS,
            "MOMCONFIG": MOMCONFIG,
            "RESOURCEGROUP": RESOURCEGROUP,
            "STARTEMPLOYEE": STARTEMPLOYEE,
            "RESOURCENAMES": RESOURCENAMES,
            "PARENTWIPORDERNO": PARENTWIPORDERNO,
            "PARENTWIPORDERTYPE": PARENTWIPORDERTYPE,
            "LOCATION": LOCATION,
            "OPRSEQUENCENO": OPRSEQUENCENO,
            "SKILL": SKILL,

            'steps': _steps,
        }
        ret.append(vals)
        # fixme: 当前demo数据还没补全，提前返回
        # return ret, False
    return ret, False


# def convert_ts002_order(env, vals):
#     ret = vals
#     mes_work_steps = vals.get('steps')
#     if not mes_work_steps:
#         raise ValidationError("Can Not Get Work Step From External System")
#     work_steps = list(itertools.chain.from_iterable(mes_work_steps))  # flat array make sure
#     vals['steps'] = work_steps
#     tightening_steps = filter(lambda step: step.get('test_type') == 'tightening', work_steps)
#     if not tightening_steps:
#         return vals
#     for ts in tightening_steps:
#         tc = ts.get('code')
#         ws = env['sa.quality.point'].search(['|', ('ref', '=', tc), ('name', '=', tc)])
#         if not ws:
#             _logger.error("Can Not Found Tightening Step By Code:{0}".format(tc))
#             continue
#         rws = ws
#         if len(ws) != 1:
#             _logger.error("Tightening Step By Code:{0} Is Not Unique".format(tc))
#             rws = ws[0]
#         expect_tightening_total = ts.get('tightening_total', 0)
#         if not expect_tightening_total:
#             raise ValidationError('Can Not Found Tightening Total Within The Tightening Step: {0}'.format(tc))
#         if len(rws.operation_point_ids) != expect_tightening_total:
#             _logger.error(
#                 "Tightening Count Is Not Equal Within Tightening System By Code:{0}. "
#                 "Except:{1}, Real:{2}".format(tc, expect_tightening_total, len(rws.operation_point_ids)))
#
#         if not rws.worksheet_img:
#             raise ValidationError('Can Not Found Tightening Image Within The Tightening Step: {0}'.format(tc))
#         ts.update({'tightening_image_by_step_code': tc})
#         val = package_tightening_points(rws.operation_point_ids)
#         ts.update({'tightening_points': val})  # 将拧紧点的包包裹进去
#     return ret

def package_workcenter_location_data(workcenter_id, val):
    ret = []
    components = []
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
            'equipment_sn': loc.equipment_id.serial_no if loc.equipment_id else False,
            'io_output': loc.io_output,
            'io_input': loc.io_input
        }
        component = {
            'is_key': True,
            'code': loc.product_id.default_code if loc.product_id else False,
        }
        ret.append(data)
        components.append(component)
    entry.update({'locations': ret})
    val['workcenter'] = entry
    val['components'] = components


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
    try:
        resp = requests.post(master_url, data=json.dumps(data), headers=headers)
        if resp.status_code != 201:
            return False
        return True
    except RetryError:
        return False
    except Exception:
        return False


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
    # validate_ts002_order_req_vals(values)  # make sure field is all in request body
    # validate_ts002_req_val_by_entry(values)
    payloads, err = convert_ts002_order(env, values)
    if err:
        return sa_fail_response(msg=err)
    result_list = list()
    for payload in payloads:
        workcenter_id, master_url = get_masterpc_order_url_and_package(env, payload)
        package_workcenter_location_data(workcenter_id, payload)
        _logger.debug("TS002 Get Order: {0}".format(pprint.pformat(payload)))

        resp = post_order_2_masterpc(master_url, [payload])
        result_list.append(resp)
    if all(result_list):
        msg = 'TS002 Post WorkOrder To MasterPC Fail'
        return sa_fail_response(msg=msg)
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

    @route('/api/v1/ts002/workorders', type='http', methods=['GET', 'OPTIONS'], auth='none', cors='*', csrf=False)
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

    @route('/api/v1/ts002/mrp.workcenter/<string:ref>/OperationDownload', type='http', methods=['PUT', 'OPTIONS'],
           auth='none', cors='*', csrf=False)
    def tangcheWorkCenterOperationSync(self, ref, **kw):
        try:
            if not ref:
                raise ValidationError('Work Center Code Must Be Include')
            env = api.Environment(request.cr, SUPERUSER_ID, request.context)
            work_center_sudo = env['mrp.workcenter'].sudo()
            workcenter_id = work_center_sudo.search(['|', ('code', '=', ref), ('name', '=', ref)])
            if not workcenter_id:
                raise ValidationError('Can Not Found Work Center By Ref:{0}'.format(ref))
            if len(workcenter_id) > 1:
                raise ValidationError('Work Center By Ref:{0} Is Not Unique!'.format(ref))
            workcenter_id.button_sync_operations()
            msg = "Sync Operation Success"
            return sa_success_resp(status_code=200, msg=msg)
        except RetryError as err:
            msg = str(err)
            return sa_fail_response(msg=msg)
        except Exception as e:
            msg = getattr(e, 'message') or getattr(e, 'name')
            return sa_fail_response(msg=msg)
