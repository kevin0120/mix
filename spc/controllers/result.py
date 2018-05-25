# -*- coding: utf-8 -*-
from odoo import http,fields,api, SUPERUSER_ID
import json
from odoo.http import request,Response
from dateutil import parser
import requests as Requests
from requests import ConnectionError,RequestException

DEFAULT_LIMIT = 80

NORMAL_RESULT_FIELDS_READ = ['workorder_id', 'id', 'product_id', 'consu_product_id', 'op_time', 'measure_result', 'workcenter_id']


def _post_aiis_result_package(aiis_urls, results):
    if not aiis_urls:
        return False
    for url in aiis_urls:
        for result in results:
            data = {
                'agna': result.production_id.equipment_name,
                'fname': result.production_id.factory_name,
                'year': result.production_id.year,
                'pin': result.production_id.pin,
                'pin_check': result.production_id.pin_check_code,
                'assembly_line': result.production_id.assembly_line_id.code,
                'lnr': result.production_id.lnr,
                'nut_no': result.consu_product_id.screw_type_code,
                'date': result.control_date,
                'result': result.measure_result.upper(),
                'MI': result.measure_torque,
                'WI': result.measure_degree
            }
            try:
                ret = Requests.put(url, data=json.dumps(data), headers={'Content-Type': 'application/json'})
                if ret.status_code == 200:
                    result.write({'sent': True})  ### 更新发送结果
            except ConnectionError:
                break  # 退出循环,进入下个aiis发送节点
            except RequestException as e:
                print(e)
    return True


class SPC(http.Controller):
    @http.route(['/api/v1/operation.results/<int:result_id>'], methods=['PUT','OPTIONS'], type='json', auth='none', cors='*', csrf=False)
    def _update_results(self, result_id):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        operation_result_id = env['operation.result'].browse(result_id)
        if not operation_result_id.exists():
            body = json.dumps({'msg': "result {0} not existed".format(result_id)})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=404, headers=headers)
            return response
        else:
            vals = request.jsonrequest
            if 'cur_objects' in vals:
                vals.update({
                    'cur_objects': json.dumps(vals['cur_objects'])
                })

            if 'control_date' in vals:
                _t = parser.parse(vals['control_date']) if vals['control_date'] else None
                if _t:
                    vals.update({
                        'control_date': fields.Datetime.to_string((_t - _t.utcoffset()))
                    })
            ret = operation_result_id.write(vals)
            if not ret:
                body = json.dumps({'msg': "update result %d fail" % result_id})
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                response = Response(body, status=405, headers=headers)
                return response
            if operation_result_id.measure_result == 'ok' or (operation_result_id.measure_result == 'nok' and operation_result_id.op_time >= operation_result_id.point_id.times):
                _aiis_urls = env['ir.config_parameter'].sudo().get_param('aiis.urls')
                if _aiis_urls:
                    aiis_urls = _aiis_urls.split(',')
                    ret = _post_aiis_result_package(aiis_urls, operation_result_id)
            # val = {
            #     'workorder_id': operation_result_id.workorder_id.id,
            #     'id':  operation_result_id.id,
            #     'product_id': operation_result_id.product_id.id,
            #     'consu_product_id': operation_result_id.consu_product_id.id,
            #     'op_time': operation_result_id.op_time,
            #     'measure_result': operation_result_id.measure_result,
            #     'workcenter_id': operation_result_id.workcenter_id.id
            # }
            # body = json.dumps(val)
            headers = [('Content-Type', 'application/json')]
            response = Response(status=204, headers=headers)
            return response

    @http.route(['/api/v1/operation.results'], methods=['PUT', 'OPTIONS'], type='json', auth='none', cors='*', csrf=False)
    def _batch_update_results(self):
        datas = request.jsonrequest
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        result_ids = [val['id'] for val in datas]
        operation_result_ids = env['operation.result'].search([('id', 'in', result_ids)])
        if not operation_result_ids:
            body = json.dumps({'msg': "result not existed"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=404, headers=headers)
            return response
        diff = list(set(result_ids) - set(operation_result_ids.ids))
        if len(diff) > 0:
            body = json.dumps({'msg': "result {0} not existed".format(diff)})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=404, headers=headers)
            return response
        for val in datas:
            if 'cur_objects' in val:
                val.update({
                    'cur_objects': json.dumps(val['cur_objects'])
                })

            if 'control_date' in val:
                _t = parser.parse(val['control_date']) if val['control_date'] else None
                if _t:
                    val.update({
                        'control_date': fields.Datetime.to_string((_t - _t.utcoffset()))
                    })
            # result_id = val.pop('id')
            # need_update_result = operation_result_ids.filtered(lambda r: r.id == result_id)
        ret = operation_result_ids.bulk_write(datas)
        if not ret:
            body = json.dumps({'msg': "update result fail"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=405, headers=headers)
            return response

        headers = [('Content-Type', 'application/json')]
        response = Response(status=204, headers=headers)
        return response

    @http.route(['/api/v1/operation.results', '/api/v1/operation.results/<int:result_id>'], type='http', auth='none', methods=['get'], cors='*', csrf=False)
    def _get_result_lists(self, result_id=None, **kw):
        domain = []
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        if result_id:
            quality_checks = env['operation.result'].search([('id', '=', result_id)])
        else:
            if 'date_from' in kw.keys():
                _t = parser.parse(kw['date_from'])
                domain += [('control_date', '>=', fields.Datetime.to_string((_t - _t.utcoffset())) )]
            if 'date_to' in kw.keys():
                _t = parser.parse(kw['date_to'])
                domain += [('control_date', '<=', fields.Datetime.to_string((_t - _t.utcoffset())) )]
            if 'limit' in kw.keys():
                limit = int(kw['limit'])
            else:
                limit = DEFAULT_LIMIT
            quality_checks = env['operation.result'].search(domain, limit=limit)
        if not quality_checks:
            body = json.dumps({'msg': "result not existed"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
        vals = [{
                'workorder_id': operation_result_id.workorder_id.id,
                'id':  operation_result_id.id,
                'product_id': operation_result_id.product_id.id,
                'consu_product_id': operation_result_id.consu_product_id.id,
                'op_time': operation_result_id.op_time,
                'measure_result': operation_result_id.measure_result,
                'workcenter_id': operation_result_id.workcenter_id.id
            } for operation_result_id in quality_checks]
        ret = vals[0] if result_id else vals
        body = json.dumps(ret)
        return Response(body, headers=[('Content-Type', 'application/json'),('Content-Length', len(body))], status=200)

    @http.route('/api/v1/operation.results/<int:result_id>/curves_add', type='json', methods=['PATCH', 'OPTIONS'], auth='none', cors='*', csrf=False)
    def _append_curves(self, result_id):
        env = api.Environment(request.cr, SUPERUSER_ID, request.context)
        operation_result_id = env['operation.result'].browse(result_id)
        if not operation_result_id:
            body = json.dumps({'msg': "result %d not existed" % result_id})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            response = Response(body, status=404, headers=headers)
            return response
        else:
            vals = request.jsonrequest
            write_values = dict()
            _vals = json.loads(operation_result_id.cur_objects if operation_result_id.cur_objects else json.dumps([]))
            _vals.append(vals)
            write_values['cur_objects'] = json.dumps(_vals)
            ret = operation_result_id.write(write_values)
            if ret:
                val = {
                    'workorder_id': operation_result_id.workorder_id.id,
                    'id': operation_result_id.id,
                    'product_id': operation_result_id.product_id.id,
                    'consu_product_id': operation_result_id.consu_product_id.id,
                    'op_time': operation_result_id.op_time,
                    'measure_result': operation_result_id.measure_result,
                    'workcenter_id': operation_result_id.workcenter_id.id
                }
                body = json.dumps(val)
                headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
                response = Response(body, status=200, headers=headers)
                return response
