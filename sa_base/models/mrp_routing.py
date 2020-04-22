# -*- coding: utf-8 -*-

from odoo import models, fields, api, _, SUPERUSER_ID

from odoo.exceptions import UserError, ValidationError

import logging
import json

_logger = logging.getLogger(__name__)


class OperationStepRel(models.Model):
    _name = 'work.step.operation.rel'

    _log_access = False

    _order = "sequence asc, id desc"

    sequence = fields.Integer('Sequence', default=0)

    operation_id = fields.Many2one('mrp.routing.workcenter', required=True)

    step_id = fields.Many2one('sa.quality.point')

    @api.multi
    def name_get(self):
        res = []
        for point in self:
            res.append((point.id, _('#%s[%s]') % (point.sequence, point.operation_id.name)))
        return res


class MrpRoutingWorkcenter(models.Model):
    _inherit = 'mrp.routing.workcenter'
    """重写routing_id的定义"""
    routing_id = fields.Many2one(
        'mrp.routing', 'Parent Routing',
        index=True, ondelete='set null', required=False,
        help="The routing contains all the Work Centers used and for how long. This will create work orders afterwards"
             "which alters the execution of the manufacturing order. ")

    sa_routing_ids = fields.Many2many('mrp.routing', 'routing_operation_rel', 'operation_id', 'routing_id',
                                      string="Routes", copy=False)

    workcenter_group_id = fields.Many2one('mrp.workcenter.group', string='Operation Work Center Group', copy=False,
                                          required=True)
    workcenter_ids = fields.Many2many('mrp.workcenter', related='workcenter_group_id.sa_workcenter_ids',
                                      string='Work Center Set', copy=False, readonly=True)

    sa_step_ids = fields.One2many('work.step.operation.rel', 'operation_id', string="Steps", copy=False)

    prefer_workcenter_id_domain = fields.Char(compute='_compute_workcenter_group', readonly=True, store=False, )

    workcenter_id = fields.Many2one('mrp.workcenter', string='Prefer Work Center', copy=False,
                                    required=False)

    socket = fields.Char(string='Bolt Socket No')
    op_job_id = fields.Many2one('controller.job', string='Job', track_visibility="onchange")
    # op_job_id = fields.Many2one('controller.job', string='Job')

    # operation_point_ids = fields.One2many('operation.point', 'operation_id', string='Operation Points')
    #
    # operation_point_group_ids = fields.One2many('operation.point.group', 'operation_id',
    #                                             string='Operation Points Group(multi-spindle)')

    group_id = fields.Many2one('mrp.routing.group', string='Routing Group')
    max_redo_times = fields.Integer('Operation Max Redo Times', default=3,
                                    track_visibility="onchange")  # 此项重试业务逻辑在HMI中实现

    sync_download_time = fields.Datetime(string=u'同步下发时间')

    max_op_time = fields.Integer('Max Operation time(second)', default=60, track_visibility="onchange")

    routing_tracking_count = fields.Integer(compute='_compute_routing_tracking_count', string="Routing Modification")

    _sql_constraints = [('routing_group_wc_uniq', 'unique(routing_id,group_id, workcenter_id)',
                         'Per Routing only has one unique Routing group per Work Center!')]

    step_count = fields.Integer(string='Steps', compute='_compute_step_count')

    active = fields.Boolean(default=True)

    @api.multi
    @api.depends('workcenter_group_id')
    def _compute_workcenter_group(self):
        for routing in self:
            workcenter_ids = routing.workcenter_group_id and routing.workcenter_group_id.sa_workcenter_ids
            if not workcenter_ids:
                continue
            routing.workcenter_id = workcenter_ids[0]  # 重新设置优先工位
            routing.prefer_workcenter_id_domain = json.dumps([('id', 'in', workcenter_ids.ids)])

    @api.depends('sa_step_ids')
    def _compute_step_count(self):
        for routing in self:
            dd = routing.sa_step_ids.mapped('step_id').filtered(lambda r: r and r.test_type != 'tightening_point')
            routing.step_count = len(dd)

    @api.multi
    def action_sa_show_steps(self):
        self.ensure_one()
        action = self.env.ref('sa_base.sa_mrp_step_action').read()[0]
        ids = self.ids
        picking_type_id = self.env['stock.picking.type'].search([('code', '=', 'mrp_operation')], limit=1).id
        test_type_id = self.env.ref('quality.test_type_text').id
        ctx = dict(self._context, default_picking_type_id=picking_type_id, default_company_id=self.company_id.id,
                   default_operation_id=self.id,
                   default_test_type_id=test_type_id)
        action.update({'context': ctx})
        action['domain'] = [('id', 'in', self.sa_step_ids.mapped('step_id').ids),
                            ('test_type', '!=', 'tightening_point')]
        action['name'] = _("Work Steps")
        return action

    @api.multi
    def action_sa_view_routing_tracking(self):
        self.ensure_one()
        action = self.env.ref('mail.action_view_mail_message').read()[0]
        # # workcenter_id = self.env.ref('sa_base.cunrong_default_workcenter').id
        # ids = self.ids
        # # bom specific to this variant or global to template
        action['context'] = {
            'search_default_model': 'mrp.routing.workcenter',
            'search_default_res_id': self.id
        }
        # action['domain'] = [('routing_id', 'in', [self.ids])]
        return action

    @api.multi
    def _track_subtype(self, init_values):
        # if 'op_job_id' in init_values:
        #     return 'account_voucher.mt_voucher_state_change'
        return super(MrpRoutingWorkcenter, self)._track_subtype(init_values)

    @api.multi
    def _compute_routing_tracking_count(self):
        for routing in self:
            routing.routing_tracking_count = len(routing.message_ids)

    @api.one
    def _push_mrp_routing_workcenter(self, url):
        return True

    @api.multi
    def button_send_mrp_routing_workcenter(self):
        return True

    @api.multi
    def get_operation_points(self):
        if not self:
            return []
        self.ensure_one()
        vals = []
        for point in self.operation_point_ids:
            vals.append({
                'sequence': point.sequence,
                'x_offset': point.x_offset,
                'y_offset': point.y_offset
            })
        return vals

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
        for idx, point in enumerate(not_sort_list.sorted(key=lambda r: r.sequence)):
            point.write({'group_sequence': group_idx + idx + 1})
        for idx, point in enumerate(self.operation_point_ids.sorted(key=lambda r: r.group_sequence)):
            point.write({'sequence': idx + 1})

    @api.multi
    def unlink(self):
        if self.env.uid != SUPERUSER_ID:
            raise ValidationError(u'不允许删除作业')

    @api.multi
    def name_get(self):
        return [(operation.id,
                 u"[{0}]{1}@{2}/{3}".format(operation.name, operation.group_id.code, operation.workcenter_id.name,
                                            operation.routing_id.name)) for operation in self]  # 强制可视化时候名称显示的是code


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

    control_mode = fields.Selection([('pset', 'Parameter Set'), ('job', 'Assembly Process')],
                                    default='pset', string='Control Mode For Tightening')

    # @api.onchange('code', 'strategy')
    # def _onchange_code_style(self):
    #     for program in self:
    #         program.name = u"{0}({1})".format(program.code, program.strategy)

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

    '''重写operation_ids的定义'''
    sa_operation_ids = fields.Many2many('mrp.routing.workcenter', 'routing_operation_rel', 'routing_id', 'operation_id',
                                        string="Operations", copy=True)
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
            routing.operation_count = len(routing.sa_operation_ids)

    @api.multi
    def action_sa_view_operation(self):
        action = self.env.ref('sa_base.sa_mrp_routing_workcenter_action').read()[0]
        # workcenter_id = self.env.ref('sa_base.cunrong_default_workcenter').id
        ids = self.ids
        # bom specific to this variant or global to template
        action['context'] = {
            'default_sa_routing_ids': [(4, ids[0], None)]
            # 'default_workcenter_id': self.workcenter_id.id
        }
        action['domain'] = [('sa_routing_ids', 'in', self.ids)]
        return action

    @api.multi
    def name_get(self):
        return [(routing.id, u'[{0}]{1}'.format(routing.code, routing.name)) for routing in self]  # 强制可视化时候名称显示的是code
