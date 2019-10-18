# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError, UserError

import requests as Requests

from requests import ConnectionError, RequestException, exceptions

import json

MASTER_DEL_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter.delete'

MASTER_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter'


class MrpBom(models.Model):
    _inherit = 'mrp.bom'
    operation_ids = fields.Many2many('mrp.routing.workcenter', related='routing_id.sa_operation_ids',
                                     copy=False, readonly=True)

    has_operations = fields.Boolean(compute='_compute_has_operations')

    @api.multi
    @api.depends('operation_ids')
    def _compute_has_operations(self):
        for bom_id in self:
            bom_id.has_operations = True if len(bom_id.operation_ids) else False

    #
    # @api.multi
    # def button_add_operation(self):
    #     self.ensure_one()
    #     compose_form = self.env.ref('sa_base.mrp_bom_operation_wizard_from', False)
    #     ctx = dict(
    #         self.env.context,
    #         default_routing_id=self.routing_id.id if self.routing_id else False
    #         )
    #
    #     return {
    #         'name': _('MRP Operation Setting'),
    #         'type': 'ir.actions.act_window',
    #         'view_type': 'form',
    #         'view_mode': 'form',
    #         'res_model': 'mrp.routing.wc.form',
    #         'views': [(compose_form.id, 'form')],
    #         'view_id': compose_form.id,
    #         'target': 'new',
    #         'context': ctx,
    #     }

    @api.multi
    def button_resequence(self):
        self.ensure_one()
        for idx, bom_line_id in enumerate(self.bom_line_ids):
            bom_line_id.write({'sequence': idx + 1})

    @api.onchange('routing_id', 'product_id')
    def _onchange_routing_id(self):
        self.code = u'[{0}]{1}'.format(self.routing_id.name, self.product_id.name)
        self.bom_line_ids = [(5,)]  # 删除所有BOM行

    @api.constrains('product_id', 'product_tmpl_id')
    def _product_tmpl_product_constraint(self):
        if self.product_id.product_tmpl_id.id != self.product_tmpl_id.id:
            raise ValidationError(_(u'The product template "%s" is invalid on product with name "%s"') % (
                self.product_tmpl_id.name, self.product_id.name))

    @api.constrains('product_id', 'routing_id', 'active')
    def _constraint_active_product_routing(self):
        if not self.active:
            return
            ###只有激活状态才检查
        count = self.env['mrp.bom'].search_count(
            [('id', '!=', self.id), ('product_id', '=', self.product_id.id), ('routing_id', '=', self.routing_id.id),
             ('active', '=', True)])
        if count:
            raise ValidationError(
                _(u'The product had a related routing config "%s" been actived!') % (self.product_id.name))

    @api.constrains('product_tmpl_id', 'routing_id', 'active')
    def _constraint_active_product_tmpl_routing(self):
        if not self.active:
            return
            ###只有激活状态才检查
        count = self.env['mrp.bom'].search_count(
            [('id', '!=', self.id), ('product_tmpl_id', '=', self.product_tmpl_id.id),
             ('routing_id', '=', self.routing_id.id), ('active', '=', True)])
        if count:
            raise ValidationError(
                _(u'The product Template had a related routing config "%s" been actived!') % (
                    self.product_tmpl_id.name))

    @api.model
    def create(self, vals):
        auto_operation_inherit = self.env['ir.values'].get_default('sa.config.settings', 'auto_operation_inherit')
        if auto_operation_inherit and 'routing_id' in vals:
            routing_id = self.env['mrp.routing'].browse(vals['routing_id'])
            operation_ids = routing_id.operation_ids
            vals.update({'operation_ids': [(6, None, operation_ids.ids)]})
        ret = super(MrpBom, self).create(vals)
        return ret

    @api.multi
    def write(self, vals):
        ret = super(MrpBom, self).write(vals)
        return ret

    @api.multi
    def unlink(self):
        raise ValidationError(u'不允许删除物料清单')

    @api.one
    def button_send_mrp_routing_workcenter(self):
        if not self.operation_ids:
            return True
        for operation in self.operation_ids:
            master = operation.workcenter_id.masterpc_id if operation.workcenter_id else None
            if not master:
                continue
            connections = master.connection_ids.filtered(
                lambda r: r.protocol == 'http') if master.connection_ids else None
            if not connections:
                continue
            url = \
                ['http://{0}:{1}{2}'.format(connect.ip, connect.port, MASTER_WROKORDERS_API) for connect in
                 connections][0]

            operation._push_mrp_routing_workcenter(url)
        return True


class MrpBomLine(models.Model):
    _inherit = 'mrp.bom.line'

    @api.depends('work_step_id')
    def _compute_product_qty(self):
        operation_point_sudo = self.env['operation.point'].sudo()
        for line in self:
            operation_point_id = operation_point_sudo.search([('qcp_id', '=', line.work_step_id.id)])
            if not operation_point_id:
                line.product_qty = 1.0
            line.product_qty = operation_point_id.product_qty

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    work_step_id = fields.Many2one('sa.quality.point', ondelete='cascade', required=True)

    product_qty = fields.Float('Product Quantity', default=1.0, required=True, compute=_compute_product_qty, store=True)

    group_id = fields.Many2one('mrp.routing.group', related="operation_id.group_id", string='Routing Group')

    program_id = fields.Many2one('controller.program', related="work_step_id.program_id", readonly=True,
                                 string='Tightening Program(Pset/Assembly Process)')

    workcenter_ids = fields.Many2many('mrp.workcenter', related="operation_id.workcenter_ids", string='Work Centers',
                                      readonly=True)

    workcenter_id = fields.Many2one('mrp.workcenter', string='Prefer Work Center', copy=False,
                                    related='operation_id.workcenter_id', required=False)

    @api.model
    def create(self, vals):
        line = super(MrpBomLine, self).create(vals)
        return line

    @api.multi
    def write(self, vals):
        res = super(MrpBomLine, self).write(vals)
        if 'product_qty' not in vals:
            return res
        for line in self:
            if not line.work_step_id:
                continue
            line.work_step_id.write({'product_qty': line.product_qty})
        return res

    @api.multi
    def _push_del_routing_workcenter(self, line, url):
        val = [{
            'product_type': line.bom_id.product_id.default_code,
            "id": line.operation_id.id,
        }]
        try:
            ret = Requests.put(url, data=json.dumps(val), headers={'Content-Type': 'application/json'}, timeout=1)
            if ret.status_code == 204:
                self.env.user.notify_info(u'删除工艺成功')
                return True
        except ConnectionError as e:
            self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))
            return False
        except RequestException as e:
            self.env.user.notify_warning(u'下发工艺失败, 错误原因:{0}'.format(e.message))
            return False
