# -*- coding: utf-8 -*-
from odoo import api, SUPERUSER_ID
from odoo.http import request, Response, Controller, route
from odoo.exceptions import ValidationError
from odoo.addons.spc import MASTER_WROKORDERS_API
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


class TangcheMrpOrderGateway(Controller):
    @route('/ts002/workorders', type='json', methods=['PUT', 'OPTIONS'], auth='none', cors='*', csrf=False)
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

            body = json.dumps({'msg': "Tightening System Create Work Order Success"})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=201, headers=headers)

        except Exception as e:
            body = json.dumps({'msg': e.message})
            headers = [('Content-Type', 'application/json'), ('Content-Length', len(body))]
            return Response(body, status=404, headers=headers)
