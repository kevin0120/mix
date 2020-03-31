# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError
import math, json
import urllib
import logging
from requests import Request as Requests, ConnectionError, RequestException

_logger = logging.getLogger(__name__)

from odoo.addons.sa_base.models.mrp_worksegment import DELETE_ALL_MASTER_WROKORDERS_API


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    external_url = fields.Text('Work Center External URL', compute='_compute_external_url')

    hmi_id = fields.Many2one('maintenance.equipment', string='Human Machine Interface(HMI)', copy=False,
                             domain=lambda self: [('category_id', '=', self.env.ref('sa_base.equipment_hmi').id)])
    masterpc_id = fields.Many2one('maintenance.equipment', string='Work Center Controller(MasterPC)', copy=False,
                                  domain=lambda self: [
                                      ('category_id', '=', self.env.ref('sa_base.equipment_MasterPC').id)])
    io_id = fields.Many2one('maintenance.equipment', string='Remote IO', copy=False,
                            domain=lambda self: [('category_id', '=', self.env.ref('sa_base.equipment_IO').id)])

    rfid_id = fields.Many2one('maintenance.equipment', string='Radio Frequency Identification(RFID)', copy=False,
                              domain=lambda self: [('category_id', '=', self.env.ref('sa_base.equipment_RFID').id)])

    controller_ids = fields.Many2many('maintenance.equipment', 'controller_center_rel', 'workcenter_id', 'controller_id',
                                      string='Tightening Controllers', copy=False)

    tool_ids = fields.Many2many('maintenance.equipment', 'tool_workcenter_rel', 'workcenter_id', 'tool_id',
                               string='Tightening Tools(Tightening Gun)', copy=False)

    controller_ids_domain = fields.Char(
        compute="_compute_controller_ids_domain",
        readonly=True,
        store=False,
    )

    tool_ids_domain = fields.Char(
        compute="_compute_gun_ids_domain",
        readonly=True,
        store=False,
    )

    _sql_constraints = [('code_hmi', 'unique(hmi_id)', 'Only one HMI is allowed'),
                        ('code_rfid', 'unique(rfid_id)', 'Only one RFID is allowed'),
                        ('code_io', 'unique(io_id)', 'Only one Remote IO is allowed')]

    @api.constrains('controller_ids', 'gun_ids')
    def _constraint_equipments(self):
        self.ensure_one()
        workcenter_ids = self.env['mrp.workcenter'].sudo().search([('id', '!=', self.id)])
        for workcenter in workcenter_ids:
            # org_list = workcenter.controller_ids.ids
            # new_list = self.controller_ids.ids
            # new = len(new_list) if new_list else 0
            # org = len(org_list) if org_list else 0
            # org_list.extend(new_list)
            # if len(set(org_list)) != new + org:
            #     raise ValidationError('控制器设置重复')
            org_list = workcenter.gun_ids.ids
            new_list = self.gun_ids.ids
            new = len(new_list) if new_list else 0
            org = len(org_list) if org_list else 0
            org_list.extend(new_list)
            if len(set(org_list)) != new + org:
                raise ValidationError('拧紧枪设置重复')

    @api.multi
    def button_sync_operations(self):
        operation_obj_sudo = self.env['mrp.routing.workcenter'].sudo()
        for center in self:
            master = center.masterpc_id
            if not master:
                continue
            connections = master.connection_ids.filtered(
                lambda r: r.protocol == 'http') if master.connection_ids else None
            if not connections:
                continue
            url = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, DELETE_ALL_MASTER_WROKORDERS_API) for connect in
                   connections][0]
            center._delete_workcenter_all_opertaions(url)
            operations = operation_obj_sudo.search([('workcenter_id', '=', center.id)])
            for operation in operations:
                operation.button_send_mrp_routing_workcenter()

    @api.multi
    def _compute_external_url(self):
        base_url = self.env['ir.config_parameter'].sudo().get_param('web.base.url')
        for rec in self:
            rec.external_url = urllib.quote(
                u'{0}/web#id={1}&view_type=form&model=mrp.workcenter'.format(base_url, rec.id))

    @api.multi
    @api.depends('masterpc_id')
    def _compute_controller_ids_domain(self):
        category_id = self.env.ref('sa_base.equipment_screw_controller').id
        for rec in self:
            rec.controller_ids_domain = json.dumps(
                [('id', 'in', rec.masterpc_id.child_ids.ids), ('category_id', '=', category_id)])
            rec.controller_ids = [(5,)]  # 去除所有的枪 重新设置

    @api.multi
    @api.depends('controller_ids')
    def _compute_gun_ids_domain(self):
        category_id = self.env.ref('sa_base.equipment_Gun').id
        for rec in self:
            child_ids = rec.controller_ids.mapped('child_ids')
            rec.tool_ids_domain = json.dumps([('id', 'in', child_ids.ids), ('category_id', '=', category_id)])
            rec.tool_ids = [(5,)]  # 去除所有的枪 重新设置