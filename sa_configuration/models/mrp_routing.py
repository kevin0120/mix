# -*- coding: utf-8 -*-

from odoo import models, fields, api, _


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'

    workcenter_id = fields.Many2one('mrp.workcenter', copy=False)

    # program_id = fields.Many2one('controller.program', string='程序号')

    group_id = fields.Many2one('mrp.routing.group', string='Routing Group')
    max_redo_times = fields.Integer('Operation Max Redo Times', default=1)  # 此项重试业务逻辑在HMI中实现

    worksheet_img = fields.Binary('worksheet_img')

    max_op_time = fields.Integer('Max Operation time(second)', default=30)

    _sql_constraints = [('routing_group_wc_uniq', 'unique(routing_id,group_id, workcenter_id)',
                         'Per Routing only has one unique Routing group per Work Center!')]

    @api.multi
    def name_get(self):
        return [(operation.id, u"[{0}]{1}@{2}".format(operation.name, operation.group_id.code, operation.workcenter_id.name)) for operation in self]  # 强制可视化时候名称显示的是code


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

    @api.onchange('code', 'strategy')
    def _onchange_code_style(self):
        for program in self:
            program.name = u"{0}({1})".format(program.code, program.strategy)

    @api.model
    def create(self, vals):
        if 'name' not in vals:
            vals['name'] = u"{0}({1})".format(vals['code'], vals['strategy'])
        return super(ControllerProgram, self).create(vals)


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
        action = self.env.ref('sa_configuration.sa_mrp_routing_workcenter_action').read()[0]
        # workcenter_id = self.env.ref('sa_configuration.cunrong_default_workcenter').id
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
