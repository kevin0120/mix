# -*- coding: utf-8 -*-

from odoo import models, fields, api, _, SUPERUSER_ID

from odoo.exceptions import UserError




class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    workcenter_id = fields.Many2one('mrp.workcenter', copy=False)

    socket = fields.Char(string='Socket No')

    op_job_id = fields.Many2one('controller.job', string='Job')

    operation_point_ids = fields.One2many('operation.point', 'operation_id', string='Operation Points')

    operation_point_group_ids = fields.One2many('operation.point.group', 'operation_id',
                                                string='Operation Points Group(multi-spindle)')

    group_id = fields.Many2one('mrp.routing.group', string='Routing Group')
    max_redo_times = fields.Integer('Operation Max Redo Times', default=3)  # 此项重试业务逻辑在HMI中实现

    worksheet_img = fields.Binary('worksheet_img')

    sync_download_time = fields.Datetime(string=u'同步下发时间')

    max_op_time = fields.Integer('Max Operation time(second)', default=60)

    _sql_constraints = [('routing_group_wc_uniq', 'unique(routing_id,group_id, workcenter_id)',
                         'Per Routing only has one unique Routing group per Work Center!')]


    @api.multi
    def button_resequence(self):
        self.ensure_one()
        has_sort_point_list = self.env['operation.point']
        group_idx = 0
        need_add = False
        for idx, point_group in enumerate(self.operation_point_group_ids):
            point_group.write({'sequence': idx + 1})
            for point in point_group.operation_point_ids:
                need_add = True
                point.write({'group_sequence': group_idx + 1})
                has_sort_point_list += point
            if need_add:
                need_add = False
                group_idx += 1
        not_sort_list = self.operation_point_ids - has_sort_point_list
        for idx, point in enumerate(not_sort_list):
            point.write({'group_sequence': group_idx + idx + 1})
        for idx, point in enumerate(self.operation_point_ids.sorted(key=lambda r: r.group_sequence)):
            point.write({'sequence': idx + 1})

    @api.multi
    def name_get(self):
        return [(operation.id, u"[{0}]{1}@{2}/{3}".format(operation.name, operation.group_id.code, operation.workcenter_id.name, operation.routing_id.name)) for operation in self]  # 强制可视化时候名称显示的是code


class MrpPR(models.Model):
    _name = 'mrp.routing.group'
    _description = 'Manufacutre Routing Group'

    code = fields.Char(string='Routing Group Rerence', required=True)

    name = fields.Char(string='Routing Group')

    _sql_constraints = [('name_uniq', 'unique(name)',
                         'Routing Group name must be unique!'),
                        ('code_uniq', 'unique(code)',
                         'Routing Group code must be unique!')
                        ]

    @api.onchange('code')
    def _routing_group_code_change(self):
        for routing in self:
            routing.name = routing.code

    @api.onchange('name')
    def _routing_group_name_change(self):
        for routing in self:
            routing.code = routing.name

    @api.multi
    def name_get(self):
        return [(group.id, group.code) for group in self]  # 强制可视化时候名称显示的是code


class ControllerJob(models.Model):
    _name = 'controller.job'
    _description = 'Controller Job'

    name = fields.Char('Job Name')
    code = fields.Char('Job Code', required=True, help=u'Job')

    active = fields.Boolean('Active', default=True)

    description = fields.Html('Description')

    @api.multi
    def unlink(self):
        if self.env.uid != SUPERUSER_ID:
            raise UserError(_(u"Only SuperUser can delete Job ID"))
        return super(ControllerJob, self).unlink()


class ControllerProgram(models.Model):
    _name = 'controller.program'
    _description = 'Controller Program'

    name = fields.Char('Program Name')
    code = fields.Char('Program Code', required=True, help=u'程序号')
    strategy = fields.Selection([('AD', 'Torque tightening'),
                                 ('AW', 'Angle tightening'),
                                 ('ADW', 'Torque/Angle tightening'),
                                 ('LN', 'Loosening'),
                                 ('AN', 'Number of Pulses tightening'),
                                 ('AT', 'Time tightening')], required=True)

    active = fields.Boolean('Active', default=True)

    @api.onchange('code', 'strategy')
    def _onchange_code_style(self):
        for program in self:
            program.name = u"{0}({1})".format(program.code, program.strategy)

    @api.model
    def create(self, vals):
        if 'name' not in vals:
            vals['name'] = u"{0}({1})".format(vals['code'], vals['strategy'])
        return super(ControllerProgram, self).create(vals)

    @api.multi
    def unlink(self):
        if self.env.uid != SUPERUSER_ID:
            raise UserError(_(u"Only SuperUser can delete program"))
        return super(ControllerProgram, self).unlink()


class MrpRouting(models.Model):
    """ Specifies routings of work centers """
    _inherit = 'mrp.routing'

    code = fields.Char('Reference', copy=False)

    operation_count = fields.Integer(string='Operations', compute='_compute_operation_count')


    @api.onchange('code')
    def _routing_code_change(self):
        for routing in self:
            routing.name = routing.code

    @api.onchange('name')
    def _routing_code_change(self):
        for routing in self:
            routing.code = routing.name

    @api.depends('operation_ids')
    def _compute_operation_count(self):
        for routing in self:
            routing.operation_count = len(routing.operation_ids)

    @api.multi
    def action_sa_view_operation(self):
        action = self.env.ref('sa_base.sa_mrp_routing_workcenter_action').read()[0]
        # workcenter_id = self.env.ref('sa_base.cunrong_default_workcenter').id
        ids = self.ids
        # bom specific to this variant or global to template
        action['context'] = {
            'default_routing_id': ids[0],
            # 'default_workcenter_id': self.workcenter_id.id
        }
        action['domain'] = [('routing_id', 'in', [self.ids])]
        return action

    @api.multi
    def name_get(self):
        return [(routing.id, routing.code) for routing in self]  # 强制可视化时候名称显示的是code
