# -*- coding: utf-8 -*-

from odoo import api, fields, models, SUPERUSER_ID, _
from odoo.exceptions import UserError, AccessError, ValidationError
from odoo.tools import DEFAULT_SERVER_DATE_FORMAT, DEFAULT_SERVER_DATETIME_FORMAT
import uuid

class ApprovalCategory(models.Model):
    _name = 'approval.category'
    _description = 'HR Approval Category'

    @api.one
    @api.depends('approval_ids')
    def _compute_fold(self):
        self.fold = False if self.todo_approval_count else True

    def _compute_todo_approval_count(self):
        for record in self:
            if record.approval_ids:
                todo_request_ids = record.approval_ids.search(
                    [('category_id', '=', record.id),'|',('flow_id.employee_id.user_id','=',self.env.uid),'&',('flow_id.employee_id.agent_enable','=',True),('flow_id.employee_id.agent_user_id.user_id','=',self.env.uid),('state','=','start')])
                record.todo_approval_count = len(todo_request_ids)
            else:
                record.todo_approval_count = 0

    def _compute_favorite_approval_request_count(self):
        # 用户发起的 并未完成的。
        for record in self:
            if record.approval_ids:
                my_favorite_request_ids = record.approval_ids.search(
                    [('category_id', '=', record.id), ('owner_user_id.user_id', '=', self.env.uid), ('state', '!=', 'done')])
                record.favorite_request_count = len(my_favorite_request_ids)
            else:
                record.favorite_request_count = 0

    @api.multi
    def _compute_approval_request_needaction_count(self):
        approval_data = self.env['approval.approval'].read_group([
                ('category_id', 'in', self.ids),
                ('message_needaction', '=', True)
            ], ['category_id'], ['category_id'])
        mapped_data = dict([(m['category_id'][0], m['category_id_count']) for m in approval_data])
        for category in self:
            category.approval_request_needaction_count = mapped_data.get(category.id, 0)



    name = fields.Char('Category Name', required=True, translate=True)
    sequence_id = fields.Many2one('ir.sequence', string='Entry Sequence',
                                  help="This field contains the information related to the numbering of approval.",
                                  required=True, copy=False)
    archive = fields.Boolean(default=False,
                             help="Set archive to true to hide the approval category without deleting it.")
    color = fields.Integer('Color Index')
    description = fields.Html('Description')
    flow_line_ids = fields.One2many('approval.flow.line', 'category_id', string='Approval flow', copy=False)
    approval_ids = fields.One2many('approval.approval', 'category_id', string='Approvals', copy=False)
    todo_approval_count = fields.Integer(string="Approvals", compute='_compute_todo_approval_count')
    favorite_request_count = fields.Integer(string="Approval Requests", compute='_compute_favorite_approval_request_count')

    approval_request_needaction_count = fields.Integer(compute='_compute_approval_request_needaction_count', string="Requests")
    active = fields.Boolean(default=True,
                             help="Set archive to true to hide the approval category without deleting it.")

    _sql_constraints = [
        ('name_uniq', 'unique (name)',
         _("the approval category name must be UNIQUE!")),
    ]


class ApprovalFlow(models.Model):
    _name = 'approval.flow'
    _description = 'HR Approval Flow'
    _order = 'sequence, id'

    def _valite_user_action_auth(self):
        if self.id != self.approval_id.flow_id.id:
            raise UserError(_('Now, You can not do this operation'))
        if self.employee_id.user_id.id != self.env.uid:
            raise UserError(_('You have no access to do this operation'))
        return

    @api.multi
    def approved(self):
        """
        confirm one or more order line, update order status and create new cashmove
        """
        self._valite_user_action_auth()
        return self.approval_id.action_send_email(flow_id_state='approved')

    @api.multi
    def do_reject(self):
        self.ensure_one()
        if self.id != self.approval_id.flow_id.id:
            raise UserError(_('Now, You can not do this operation'))
        if self.employee_id.user_id.id != self.env.uid:
            raise UserError(_('You have no access to do this operation'))
        return

    @api.multi
    def _send_mail_to_approval_flows(self, template_xmlid="approval.approval_template_needaction_mail", force_send=False):

        res = False

        colors = {
            'draft': 'grey',
            'start': 'green',
            'approved': '#FFFF00',
            'rejected': 'red'
        }
        rendering_context = dict(
            self._context,
            color=colors,
            action_id=self.env.ref('approval.act_approval_approval_all').id,
            dbname=self._cr.dbname,
            base_url=self.env['ir.config_parameter'].get_param('web.base.url', default='http://localhost:8069')
        )
        invitation_template = self.env.ref(template_xmlid).with_context(rendering_context)

        mails_to_send = self.env['mail.mail']
        for flow in self:
            if flow.employee_id or flow.employee_id.email:
                mail_id = invitation_template.send_mail(flow.id)

                vals = {}
                vals['model'] = None  # We don't want to have the mail in the tchatter while in queue!
                vals['res_id'] = False
                current_mail = self.env['mail.mail'].browse(mail_id)
                current_mail.mail_message_id.write(vals)
                mails_to_send |= current_mail

        if force_send or mails_to_send:
            res = mails_to_send.send()

        return res


    @api.multi
    def rejected(self):
        """
        cancel one or more order.line, update order status and unlink existing cashmoves
        """
        self.ensure_one()
        self._valite_user_action_auth()
        return self.approval_id.action_send_email(flow_id_state='rejected')

    def _compute_actived(self):
        for record in self:
            if record.id != record.approval_id.flow_id.id:
                pass
            else:
                record.actived = True

    def _compute_have_access(self):
        uid = self.env.uid
        for record in self:
            if record.employee_id.user_id.id == uid and record.approval_id.flow_id.id == record.id:
                record.have_operation_access = True

    name = fields.Char(related='employee_id.name', store=True)

    employee_id = fields.Many2one('hr.employee', string="Auditor", required=True, domain=[('user_id', '!=', False)])
    approval_id = fields.Many2one('approval.approval', string='Approval', ondelete='cascade')
    sequence = fields.Integer('Sequence', default=10)
    actived = fields.Boolean(string='In the flow', default=False, compute=_compute_actived)
    state = fields.Selection([('draft', _('Draft')),
                              ('start', _('Start')),
                              ('rejected', _('Rejected')),
                              ('approved', _('Approved'))],
                             string='Approval State', default="draft")

    have_operation_access = fields.Boolean(string="have_operation_access",default=False,compute=_compute_have_access)

    access_token = fields.Char(
        'Security Token', copy=False, default=lambda self: str(uuid.uuid4()),
        required=True)

    # _sql_constraints = [
        # ('employee_id_approval_id_uniq', 'unique (employee_id, approval_id)',
         # _("You cannot have two same auditor with the same approval request on the same object!")),
    # ]


class ApprovalFlowLine(models.Model):
    """ Model for case stages. This models the main stages of a Approval Request management flow. """

    _name = 'approval.flow.line'
    _description = 'HR Approval Flow'
    _order = 'sequence, id'

    name = fields.Char(related='employee_id.name', store=True)
    employee_id = fields.Many2one('hr.employee',  string="Auditor", required=True, domain=[('user_id', '!=', False)])
    category_id = fields.Many2one('approval.category', string='Category', ondelete='cascade')
    sequence = fields.Integer('Sequence', default=1)

    # _sql_constraints = [
        # ('employee_id_category_id_uniq', 'unique (employee_id, category_id)',
         # _("You cannot have two same auditor with the same approval category on the same object!")),
    # ]


class approval(models.Model):
    _name = 'approval.approval'
    _inherit = ['mail.thread', 'ir.needaction_mixin']
    _description = 'HR Approval'

    @api.depends('stage_ids', 'state')
    def action_start(self):
        if self.state != 'draft':
            raise UserError('Cannot Active the approval Start!')
        self.state = 'start'
        for stage in self.stage_ids:
            stage.state = 'start'
        self._default_stage()  # 设定默认 flow_id
        self.stage_ids._send_mail_to_approval_flows()

        '''回调机制'''
        if not self.res_model or not self.res_id:
            return
        model = self.env[self.res_model].search([('id', '=', self.res_id)])
        getattr(model, 'callback_action_start')() if hasattr(model, 'callback_action_start') else None

    @api.depends('flow_id', 'state')
    def action_canceled(self):
        self.state = 'cancel'
        self.flow_id = None

        '''回调机制'''
        if not self.res_model or not self.res_id:
            return
        model = self.env[self.res_model].search([('id', '=', self.res_id)])
        getattr(model, 'callback_action_canceled')() if hasattr(model, 'callback_action_canceled') else None

    @api.depends('flow_id', 'state', 'stage_ids')
    def action_reset(self):
        self.state = 'draft'
        for stage_id in self.stage_ids:
            stage_id.write({'state': 'draft'})
        self.flow_id = None

        '''回调机制'''
        if not self.res_model or not self.res_id:
            return
        model = self.env[self.res_model].search([('id', '=', self.res_id)])
        getattr(model, 'callback_action_reset')() if hasattr(model, 'callback_action_reset') else model

    def _default_stage(self):
        if self.category_id:
            ret = self.stage_ids.sorted("sequence")[0]
            self.flow_id = ret

    @api.model
    def _read_group_stage_ids(self, stages, domain, order):
        """ Read group customization in order to display all the stages in the
            kanban view, even if they are empty
        """
        stage_ids = stages._search([], order=order, access_rights_uid=SUPERUSER_ID)
        return stages.browse(stage_ids)

    @api.multi
    def action_follow(self):
        """ Wrapper because message_subscribe_users take a user_ids=None
            that receive the context without the wrapper.
        """
        return self.message_subscribe_users()

    @api.multi
    def action_unfollow(self):
        """ Wrapper because message_unsubscribe_users take a user_ids=None
            that receive the context without the wrapper.
        """
        return self.message_unsubscribe_users()

    @api.multi
    def _message_auto_subscribe_notify(self, partner_ids):
        # Do not notify user it has been marked as follower of its employee.
        return

    @api.multi
    def _track_template(self, tracking):
        res = super(approval, self)._track_template(tracking)
        test_task = self[0]
        changes, tracking_value_ids = tracking[test_task.id]
        if 'stage_id' in changes and test_task.stage_id.mail_template_id:
            res['stage_id'] = (test_task.stage_id.mail_template_id, {'composition_mode': 'mass_mail'})
        return res

    @api.multi
    def _track_subtype(self, init_values):
        self.ensure_one()
        if 'flow_id' in init_values and self.state == 'cancel':
            return 'approval.mt_approval_canceled'
        elif 'flow_id' in init_values and self.state == 'start':
            return 'approval.mt_approval_start'
        elif 'state' in init_values and self.state == 'draft':
            return 'approval.mt_approval_created'
        elif 'state' in init_values and self.state == 'done':
            return 'approval.mt_approval_done'
        return super(approval, self)._track_subtype(init_values)


    @api.onchange('category_id')
    def _onchange_category(self):
        if not self.category_id:
            return
        else:
            self.name = self.category_id.sequence_id.next_by_id()

    @api.multi
    def write(self, values):
        result = super(approval, self).write(values)
        if 'category_id' not in values.keys() and 'stage_ids' not in values.keys():
            return result
        if 'category_id' in values.keys():
            self.with_context(dont_notify=True).create_approval_flows(force_unlink=True)
        if 'stage_ids' in values.keys():
            self.create_approval_flows()
        return result

    @api.multi
    def create_approval_flows(self, force_unlink=False):
        for approval in self:
            alreay_flow_employee_ids = approval.stage_ids.mapped('employee_id')
            flow_ids = self.env['approval.flow']
            category_employee_ids = approval.category_id.mapped('flow_line_ids').mapped('employee_id')
            for employee in category_employee_ids.filtered(lambda employee_id: employee_id not in alreay_flow_employee_ids):
                values = {
                    'approval_id': approval.id,
                    'employee_id': employee.id,
                    'state': 'draft'
                }

                flow_id = self.env['approval.flow'].create(values)
                flow_ids |= flow_id

            if flow_ids:
                user_ids = [stage_id.employee_id.user_id.id for stage_id in flow_ids]
                approval.message_subscribe_users(user_ids=user_ids)

            all_partners = approval.message_partner_ids
            all_flow_partners = approval.stage_ids.mapped('employee_id').mapped('user_id').mapped('partner_id')
            to_notify_partner = filter(lambda partner_id: partner_id not in all_partners, all_flow_partners)

            meeting_partners = self.env['res.partner']

            for partner in to_notify_partner:
                if not partner:
                    continue
                meeting_partners |= partner
            approval.message_subscribe(partner_ids=meeting_partners.ids)

            # We remove old attendees who are not in partner_ids now.
            if not force_unlink:
                return True

            employee_to_remove = alreay_flow_employee_ids.filtered(lambda employee_id: employee_id not in category_employee_ids)
            #
            if employee_to_remove:
                flow_to_remove = approval.stage_ids.search([('employee_id','in',employee_to_remove.ids)])
                flow_to_remove.unlink()

        return True

    @api.model
    def create(self, vals):
        context = dict(self.env.context, mail_create_nolog=True, mail_create_nosubscribe=True) ### 添加nolog 因为state 将会触发发送信息
        request = super(approval, self.with_context(context)).create(vals)
        request.state = 'draft'
        request.res_model = context.get('default_res_model') or None    # 获取关联模型名称的方法
        request.create_approval_flows(force_unlink=False)
        return request

    @api.multi
    def unlink(self):
        uid = self.env.uid
        if uid != SUPERUSER_ID:
            raise AccessError(_('You have no access to delete approval record'))
        for rec in self:
            if rec.state == 'start':
                raise UserError(_('你不可以删除一张正在执行中的审批单'))
        return super(approval, self).unlink()

    @api.multi
    def _compute_pending_employee_id(self):
        for recode in self:
            recode.pending_employee_id = recode.flow_id.employee_id

    @api.one
    @api.depends('flow_id')
    def change_flow_id_by_next(self):
        ret = self.stage_ids.filtered(lambda r: r.state == 'start').sorted('sequence')
        if ret:
            self.flow_id = ret.ids[0]
        else:  # 找不到结果说明所有值都为approved
            self.state = 'done'
            self.close_date = fields.Datetime.now()

    @api.model
    def _default_own_user_id(self):
        uid = self.env.uid
        employee_id = self.env['hr.employee'].search([('user_id', '=', uid)], limit=1)
        if not employee_id:
            raise AccessError('the current User is not set the related Employee')
        return employee_id

    def _compute_have_access(self):
        uid = self.env.uid
        for record in self:
            if uid == SUPERUSER_ID:
                record.have_operation_access = True
                continue
            if record.owner_user_id.user_id.id == uid:
                record.have_operation_access = True

    def action_send_email(self, flow_id_state):
        """ Open a window to compose an email, with the template - 'event_badge'
            message loaded by default
        """
        if flow_id_state == 'rejected':
            template = self.env.ref('approval.approval_reject_template_mail',False)
        elif flow_id_state == 'approved':
            template = self.env.ref('approval.approval_approve_template_mail', False)
        compose_form = self.env.ref('mail.email_compose_message_wizard_form',False)
        ctx = dict(
            self.env.context,
            default_model='approval.approval',
            default_res_id=self.ids[0],
            default_use_template=bool(template),
            default_template_id=template and template.id or False,
            mark_flow_status = flow_id_state,
            default_composition_mode='comment')

        return {
            'name': _('Approval Reject Mail'),
            'type': 'ir.actions.act_window',
            'view_type': 'form',
            'view_mode': 'form',
            'res_model': 'mail.compose.message',
            'views': [(compose_form.id, 'form')],
            'view_id': compose_form.id,
            'target': 'new',
            'context': ctx,
        }

    @api.model
    def set_related_res_model_id(self, res_id):
        """ 对外提供的设置审批关联表单id方法，在create之后调用
        """
        self.res_id = res_id if isinstance(res_id, int) else 0

    name = fields.Char('Name', required=True)
    owner_user_id = fields.Many2one('hr.employee', 'Approval Maker', readonly=True, required=True, default=_default_own_user_id)

    category_id = fields.Many2one('approval.category', string='Category', required=True)

    stage_ids = fields.One2many('approval.flow', 'approval_id', copy=False, group_expand='_read_group_stage_ids')
    flow_id = fields.Many2one('approval.flow', string='Auditor', track_visibility='onchange') # if Flow_id is None mean the approval is rejected or on Draft state
    create_date = fields.Datetime('Create Date', readonly=True, default=fields.Datetime.now)
    pending_employee_id = fields.Many2one('hr.employee', 'Current Auditor', compute='_compute_pending_employee_id')
    state = fields.Selection(   [('draft', _('Draft')),
                                  ('start', _('Start')),
                                  ('cancel', _('Cancel')),
                                  ('done', _('Done'))],
                                string='Approval State', store=True, track_visibility='always', default='draft')
    active = fields.Boolean(default=True,
                             help="Set archive to true to hide the approval request without deleting it.")

    close_date = fields.Datetime('Close Date', readonly=True)

    last_update_employee_id = fields.Many2one('hr.employee', 'Last update Employee', readonly=True)

    have_operation_access = fields.Boolean(string="have_operation_access", default=False, compute='_compute_have_access')

    description = fields.Html()

    '''审批关联表单相关字段'''
    res_model = fields.Char(string='related model name', default=None)

    res_id = fields.Integer(string='related model id', default=0)
