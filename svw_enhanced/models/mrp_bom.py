# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError, UserError

import json

from odoo.addons.sa_base.models.mrp_worksegment import DELETE_ALL_MASTER_WROKORDERS_API


class MrpBomLine(models.Model):
    _inherit = 'mrp.bom.line'

    masterpc_id = fields.Many2one('maintenance.equipment', string='Work Center Controller(MasterPC)',
                                  related="operation_id.workcenter_id.masterpc_id")

    controller_id = fields.Many2one('maintenance.equipment', string='Tightening Controller', copy=False)

    tool_id = fields.Many2one('maintenance.equipment', string='Tightening Tool(Gun/Wrench)', copy=False)

    # _sql_constraints = [
    #     ('unique_operation_bom_id', 'unique(bom_id,operation_id)', 'Every Bom unique operation'),
    # ]

    controller_id_domain = fields.Char(
        compute="_compute_gun_id_domain",
        readonly=True,
        store=False,
    )

    tool_id_domain = fields.Char(
        compute="_compute_gun_id_domain",
        readonly=True,
        store=False,
    )

    @api.onchange('operation_id')
    def _onchange_operation(self):
        self.ensure_one()
        self.controller_id = False
        self.tool_id = False

    @api.onchange('operation_point_id')
    def _onchange_operation_point_id(self):
        self.ensure_one()
        self.product_id = self.operation_point_id.product_id

    @api.onchange('controller_id')
    def _onchange_controller(self):
        self.ensure_one()
        self.tool_id = False

    @api.multi
    @api.depends('operation_id.workcenter_id')
    def _compute_gun_id_domain(self):
        for rec in self:
            rec.tool_id_domain = json.dumps([('id', 'in', rec.workcenter_id.gun_ids.ids)])
            rec.controller_id_domain = json.dumps([('id', 'in', rec.workcenter_id.controller_ids.ids)])

    @api.multi
    def unlink(self):
        for line in self:
            master = line.workcenter_id.masterpc_id if line.workcenter_id else None
            if not master:
                raise UserError(u"未找到工位上的工位控制器")
            connections = master.connection_ids.filtered(
                lambda r: r.protocol == 'http') if master.connection_ids else None
            if not connections:
                raise UserError(u"未找到工位上的工位控制器的连接信息")
            url = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, DELETE_ALL_MASTER_WROKORDERS_API) for connect in
                   connections][0]
            ret = self._push_del_routing_workcenter(line=line, url=url)
            if not ret:
                self.env.user.notify_warning(u"未删除物料清单行")
        need_unlink_quality_points = self.mapped('work_step_id')
        if need_unlink_quality_points:
            need_unlink_quality_points.sudo().unlink()
        ret = super(MrpBomLine, self).unlink()
        return ret
