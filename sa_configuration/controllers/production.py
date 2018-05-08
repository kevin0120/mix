# -*- coding: utf-8 -*-
from odoo import http
import json
from odoo.http import request,Response


class SaConfiguration(http.Controller):

    @http.route('/api/v1/mrp.productions', type='json', methods=['POST','OPTIONS'], auth='none', cors='*', csrf=False)
    def assemble_mo_create(self):
        vals = request.jsonrequest
        vin = vals['vin']
        mo_name = u'{0}--V001--{1}-{2}-{3}={4}'.format(
            vals['equipment_name'],vals['factory_name'],vals['year'],vals['pin'],vals['pin_check_code'])

        count = request.env['mrp.production'].sudo().search_count(
            [('name', '=', mo_name)])
        if count > 0:
            # MO已存在
            Response.status = "400 Bad Request"
            return {"msg": "MO name " + mo_name + " already exists"}

        count = request.env['mrp.production'].sudo().search_count(
            [('vin', '=', vin)])
        if count > 0:
            # MO已存在
            Response.status = "400 Bad Request"
            return {"msg":"MO vin " + vin + " already exists"}

        vechile_code = vals['model']
        vals.pop('model')
        records = request.env['product.product'].sudo().search(
            [('vehicle_type_code', 'ilike', vechile_code)], limit=1)

        if not records.exists():
            # 找不到对应车型
            Response.status = "400 Bad Request"
            return {"msg":"vechile model " + vechile_code + " not found"}

        product_id = records[0]

        assemble_line = vals['assembly_line']
        vals.pop('assembly_line')
        records = request.env['mrp.assemblyline'].sudo().search(
            ['|', ('name', 'ilike', assemble_line), ('code', 'ilike', assemble_line)], limit=1)

        if not records.exists():
            # 找不到对应装配线
            Response.status = "400 Bad Request"
            return {"msg": "Assembly line " + assemble_line + " not found"}

        assembly_line_id = records[0]

        vals.update({'name': mo_name})
        vals.update({'product_id': product_id.id,
                     'bom_id': product_id.active_bom_id.id,
                     'product_tmpl_id': product_id.product_tmpl_id.id,
                     'product_uom_id': product_id.active_bom_id.product_uom_id.id,
                     'routing_id': product_id.active_bom_id.product_uom_id.id,
                     'assembly_line_id': assembly_line_id.id})

        prs = vals['prs']
        vals.pop('prs')
        vals.update(
            {'production_routings': json.dumps(prs)}
        )
        production = request.env['mrp.production'].sudo().create(vals)
        production.sudo().plan_by_prs()  ### 模拟点击安排,自动生成工单

        if production:
            # 创建MO成功
            Response.status = "201 Created"
            return {"msg":""}
        else:
            # 创建MO失败
            Response.status = "400 Bad Request"
            return {"msg": "create MO failed"}
