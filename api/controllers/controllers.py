# -*- coding: utf-8 -*-
from odoo import http
from odoo.http import request, Response

from api_data import api_data,DEFAULT_LIMIT
import json

NORMAL_RESULT_FIELDS_READ = ['id', 'name', 'login', 'active', 'uuid']


class BaseApi(http.Controller):
    # 如果想使用此API,必须在配置文件中指定数据库方可使用
    @http.route('/api/v1/doc', type='http', auth='none', cors='*', csrf=False)
    def _api_doc(self):
        return json.dumps(api_data)

    @http.route('/api/v1/res.users', type='http', auth='none', cors='*', csrf=False)
    def _get_users_list_info(self, **query_params):
        _limit = DEFAULT_LIMIT
        if 'limit' in query_params:
            _limit = int(query_params['limit'])
        domain = [('id', '!=', 1)]
        if 'uuids' in query_params:
            uuids = query_params['uuids'].split(',')
            domain += [('uuid', 'in', uuids)]
        users = request.env['res.users'].sudo().search(domain, limit=_limit).read(fields=NORMAL_RESULT_FIELDS_READ)
        if len(users) != len(uuids):
            # 有某些uuid未找到
            pass
        for user in users:
            if 'active' in user:
                user.update({
                    'status': 'active' if user['active'] else 'archived'
                })
                user.pop('active')
        return Response(json.dumps(users), headers={'content-type': 'application/json'}, status=200)

    @http.route('/api/v1/res.users/<string:uuid>', type='http', auth='none', cors='*', csrf=False)
    def _get_user_info(self, uuid):
        user_id = request.env['res.users'].sudo().search([('uuid', '=', uuid)])[0]

        if not user_id:
            return Response(json.dumps({'msg': 'User not found'}), headers={'content-type': 'application/json'}, status=404)

        ret = user_id.sudo().read(fields=NORMAL_RESULT_FIELDS_READ)[0]
        if 'active' in ret:
            ret.update({
                'status': 'active' if ret['active'] else 'archived'
            })
            ret.pop('active')

        return Response(json.dumps(ret), headers={'content-type': 'application/json'}, status=200)

    @http.route('/api/v1/res.users/batch_archived', type='json', auth='none', cors='*', csrf=False)
    def _bach_patch_user_archived(self):
        uuids = request.jsonrequest
        user_ids = request.env['res.users'].sudo().search([('uuid', 'in', uuids)])

        if not user_ids:
            return Response(json.dumps({'msg': 'User not found'}), headers={'content-type': 'application/json'},
                            status=404)

        ret = user_ids.sudo().write({
            'active': False
        })
        if not ret:
            return Response(json.dumps({'msg': 'Batch Archived fail'}), headers={'content-type': 'application/json'},
                            status=405)
        ret = user_ids.sudo().read(fields=NORMAL_RESULT_FIELDS_READ)[0]
        if 'active' in ret:
            ret.update({
                'status': 'active' if ret['active'] else 'archived'
            })
            ret.pop('active')

        return Response(json.dumps(ret), headers={'content-type': 'application/json'}, status=200)

