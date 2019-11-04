# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.


import requests as Requests

from requests import ConnectionError, RequestException, exceptions

from odoo import api, fields, models, _
from odoo.exceptions import ValidationError

DELETE_ALL_MASTER_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter/all'


class MrpWorkAssembly(models.Model):
    _name = 'mrp.assemblyline'
    _description = 'Work Assembly Line'
    _order = "id"

    name = fields.Char('Assembly Line', copy=False)
    code = fields.Char('Reference', copy=False, required=True)
    worksegment_count = fields.Integer('Work Sections', compute='_compute_worksegments_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    worksegment_ids = fields.One2many('mrp.worksection', 'workassembly_id', 'Work Sections', copy=False)

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Assembly Line is allowed')]

    @api.multi
    @api.depends('worksegment_ids')
    def _compute_worksegments_count(self):
        for line in self:
            line.worksegment_count = len(line.worksegment_ids)

    @api.multi
    @api.depends('name', 'code')
    def name_get(self):
        res = []
        for line in self:
            name = u"[{0}] {1}".format(line.code, line.name)
            res.append((line.id, name))
        return res


class MrpWorkSegment(models.Model):
    _name = 'mrp.worksection'
    _description = 'Work Section'
    _order = "id"

    name = fields.Char('Work Section Name', copy=False)
    code = fields.Char('Work Section Reference', copy=False, required=True)
    workassembly_id = fields.Many2one('mrp.assemblyline', string='Work Assembly Line')
    workcenter_count = fields.Integer('Work Centers', compute='_compute_workcenters_count')

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    workcenter_ids = fields.One2many('mrp.workcenter', 'worksegment_id', 'Work Centers', copy=False)

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Section is allowed')]

    @api.multi
    @api.depends('workcenter_ids')
    def _compute_workcenters_count(self):
        for segment in self:
            segment.workcenter_count = len(segment.workcenter_ids)

    @api.multi
    @api.depends('name', 'code')
    def name_get(self):
        res = []
        for segment in self:
            name = u"[{0}] {1}".format(segment.code, segment.name)
            res.append((segment.id, name))
        return res


class MrpWorkCenterLoc(models.Model):
    _name = 'mrp.workcenter.loc'
    _description = 'Work Center Location For Picking Component'

    workcenter_id = fields.Many2one('mrp.workcenter', string='Mrp WorkCenter', required=True, copy=False)

    product_id = fields.Many2one('product.product', string='Component', copy=False)

    @api.multi
    @api.depends('workcenter_id', 'product_id')
    def name_get(self):
        res = []
        for line in self:
            name = u"[{0}]@{1}".format(line.product_id.name, line.workcenter_id.name)
            res.append((line.id, name))
        return res


class MrpWorkCenter(models.Model):
    _inherit = 'mrp.workcenter'

    type = fields.Selection([('operate', 'Operate'),
                             ('rework', 'Rework')], default='operate')

    sequence = fields.Integer(
        'Sequence', default=lambda self: int(self.env['ir.sequence'].next_by_code('mrp.workcenter.sequence')) or 1,
        required=True,
        help="Gives the sequence order when displaying a list of work centers.")

    user_ids = fields.Many2many('res.users', 'workcenter_users_rel', 'workcenter_id', 'user_id',
                                string='Responsible Persons',
                                default=lambda self: [(4, self.env.uid)])

    qc_workcenter_id = fields.Many2one('mrp.workcenter', string='Quality Check Work Center')

    worksegment_id = fields.Many2one('mrp.worksection', copy=True)

    sa_workcentergroup_ids = fields.Many2many('mrp.workcenter.group', 'mrp_workcenter_rel', 'workcenter_id', 'group_id',
                                              string="Workcenters Group", copy=False)

    sa_workcenter_loc_ids = fields.One2many('mrp.workcenter.loc', 'workcenter_id',
                                            string='Location For Place Component', copy=True)

    @api.one
    def get_workcenter_masterpc_http_connect(self):
        workcenter_id = self
        master_pc = self.env['maintenance.equipment'].search(
            [('workcenter_id', '=', workcenter_id.id), ('category_name', '=', 'MasterPC')], limit=1)
        if not master_pc:
            return None
        connections = master_pc.connection_ids.filtered(
            lambda r: r.protocol == 'http') if master_pc.connection_ids else None
        if not connections:
            raise None
        connect = connections[0]
        return connect

    @api.multi
    def button_sync_operations(self):
        pass

    @api.one
    def _delete_workcenter_all_opertaions(self, url):
        try:
            ret = Requests.delete(url, headers={'Content-Type': 'application/json'}, timeout=1)
            if ret.status_code == 200:
                self.env.user.notify_info(u'删除工艺成功')
                return True
        except ConnectionError as e:
            self.env.user.notify_warning(u'删除工艺失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'删除工艺失败, 错误原因:{0}'.format(str(e)))
        except RequestException as e:
            self.env.user.notify_warning(u'删除工艺失败, 错误原因:{0}'.format(str(e)))
            raise ValidationError(u'删除工艺失败, 错误原因:{0}'.format(str(e)))
        return False

    # @api.multi
    # def _update_create_workcenter_group_tool_by_tool(self):
    #     tool_category_ids = self.env.ref('sa_base.equipment_Gun') + self.env.ref('sa_base.equipment_Wrench')
    #     for workcenter_id in self:
    #         already_equipment_group_ids = self.env['mrp.workcenter.group.tool'].search(
    #             [('workcenter_id', '=', workcenter_id.id)])
    #         if not already_equipment_group_ids:
    #             continue
    #         already_equipment_ids = already_equipment_group_ids.mapped('tool_id')
    #         tool_ids = workcenter_id.equipment_ids.filtered(
    #             lambda r: r.category_id in tool_category_ids and r not in already_equipment_ids)
    #         for tool in tool_ids:
    #             for wg in workcenter_id.sa_workcentergroup_ids:
    #                 val = {
    #                     "workgroup_id": wg.id,
    #                     "workcenter_id": tool.workcenter_id.id,
    #                     "tool_id": tool.id,
    #                 }
    #                 self.env['mrp.workcenter.group.tool'].sudo().create(val)

    @api.multi
    def write(self, vals):
        ret = super(MrpWorkCenter, self).write(vals)
        return ret
