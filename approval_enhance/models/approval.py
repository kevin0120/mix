# -*- coding: utf-8 -*-

from odoo import api, fields, models, SUPERUSER_ID, _
from odoo.exceptions import UserError, AccessError, ValidationError
from odoo.tools import DEFAULT_SERVER_DATE_FORMAT, DEFAULT_SERVER_DATETIME_FORMAT
from odoo.workflow.helpers import Record
from odoo.workflow.workitem import Environment
from odoo.tools.safe_eval import safe_eval, test_python_expr
from datetime import date

class ApprovalCategory(models.Model):
    _inherit = 'approval.category'
    
    model_id = fields.Many2one('ir.model', string=u'应用于单据')
    res_model = fields.Char(u'模型名称', related='model_id.model')
    finish_action = fields.Text(u'动作',help=u"审批结束后单据执行的动作")
    edit_group_id = fields.Many2one('res.groups',u'可修改审批中(后)单据组')
    back_pre_step = fields.Boolean(string=u'审批拒绝回退到上一步')
    condition = fields.Text(u'单据执行该审批条件')
  
    @api.constrains('finish_action')
    def _check_python_finish_action(self):
        for category in self.filtered('finish_action'):
            msg = test_python_expr(expr=category.finish_action.strip(), mode="exec")
            if msg:
                raise ValidationError(msg)
                
class ModelApproval(models.AbstractModel):
    _name = "model.approval"
    
    @api.one
    @api.depends('approval_id','approval_id.state','approval_flow_ids','approval_flow_ids.state')
    def _compute_approval_state(self):
        if not isinstance(self.id, models.NewId) and not self.is_approval_enabled:
            self.approval_state = 'approved'
        else:
            if not self.approval_id or self.approval_id.state in ('draft','cancel'):
                self.approval_state = 'none'
            elif all(flow.state == 'approved' for flow in self.approval_flow_ids):
                self.approval_state = 'approved'
            elif any(flow.state == 'rejected' for flow in self.approval_flow_ids):
                self.approval_state = 'rejected'
            elif any(flow.state == 'start' for flow in self.approval_flow_ids):
                self.approval_state = 'in process'
            else:
                self.approval_state = 'none'
            
    
    def _set_approval_state(self):
        pass
            
    approval_id = fields.Many2one('approval.approval', string=u'当前审批', copy=False)
    approval_id_state = fields.Selection(related='approval_id.state',string=u'当前审批状态', track_visibility='none')
    approval_flow_ids = fields.One2many('approval.flow','res_id',string=u'审批进度', domain=lambda self: [('res_model', '=', self._name)])
    approval_state = fields.Selection([
        ('none', u'未审批'),
        ('in process', u'审批中'),
        ('approved', u'已审批'),
        ('rejected', u'已拒绝')], string=u'审批状态', default='none', required=True,store=True, compute=_compute_approval_state, inverse='_set_approval_state', copy=False, track_visibility='onchange')
    have_approval_access = fields.Boolean(related='approval_id.flow_id.have_operation_access')
    is_approval_enabled = fields.Boolean(u'启用审批', compute="_compute_is_approval_enabled", search="_search_is_approval_enabled")
    
    def _get_eval_context(self):
        eval_context = self.env['ir.actions.act_window']._get_eval_context()
        eval_context.update({
            # orm
            'env': self.env,
            'model': self._name,
            # record
            'record': self,
            # deprecated use record, records
            'object': self,
            'obj': self,            
        })
        
        return eval_context  
        
    def approval_expr_eval_expr(self, lines, env):
        """
        Evaluate each line of ``lines`` with the ``Environment`` environment, returning
        the value of the last line.
        """
        assert lines, 'You used a NULL action in a workflow, use dummy node instead.'
        result = False
        for line in lines.split('\n'):
            line = line.strip()
            if not line:
                continue
            if line == 'True':
                result = True
            elif line == 'False':
                result = False
            else:
                result = safe_eval(line, env, nocopy=True)

        return result
        
    def _compute_is_approval_enabled(self):
        approval_category_ids = self.env['approval.category'].search([('res_model','=',self._name)])
        result = False
        if approval_category_ids:
            if any(not approval_category_id.condition for approval_category_id in approval_category_ids):                
                result = True
                
        for record in self:
            if result == True:
                record.is_approval_enabled = True
            else:
                for approval_category_id in approval_category_ids:
                    if record.approval_expr_eval_expr(approval_category_id.condition, record._get_eval_context()):
                        record.is_approval_enabled = True
                        break
                        
    def _search_is_approval_enabled(self, operator, value):
        ids = []
        model_ids = self.search([])
        for p in model_ids:
            if p.is_approval_enabled == value:
                ids.append(p.id)
        return [('id', 'in', ids)]
    
                        
    @api.multi
    def action_submit_approval(self):
        self.ensure_one()
        approval_category_ids = self.env['approval.category'].search([('res_model','=',self._name)])
        if len(approval_category_ids)>1: #条件
            approval_category_ids = approval_category_ids.filtered(lambda ac_id: ac_id.condition and self.approval_expr_eval_expr(ac_id.condition, self._get_eval_context()))
        approval_category_id = approval_category_ids and approval_category_ids[0]
        if not approval_category_id:
            if self.env.context.get('raise_error',False):
                raise UserError(_(u'单据审批流程未定义！请联系系统管理员！'))
            else:
                return
        if not self.approval_id:            
            approval_id = self.env['approval.approval'].with_context(default_res_model=self._name).create({'category_id':approval_category_id.id,'name':approval_category_id.sequence_id.next_by_id(),'res_model':self._name,'res_id':self.id})
            approval_id.action_start()
            self.approval_id = approval_id
        else:
            self.approval_id.category_id = approval_category_id
            self.approval_id.create_approval_flows(force_unlink=True)
            if self.approval_id.state == 'cancel':
                self.approval_id.action_reset()
            self.approval_id.action_start()

    #查看审批单
    @api.multi
    def action_view_approvals(self):
        self.ensure_one()
        return  {
            'name': _(u'审批单'),
            'type': 'ir.actions.act_window',
            'view_type': 'form',
            'view_mode': 'tree,form',
            'res_model': 'approval.approval',
            'view_id': False,
            'domain': [('id', 'in', self.env['approval.approval'].search([('res_model','=',self._name),('res_id','=',self.id)]).ids)],
        }
    
    @api.multi
    def button_approved(self):
        rec = self
        initial_values = {rec.id: {'state': rec._compute_state_with_approval()}}
        if rec.approval_id:
            ret = rec.approval_id.flow_id.approved()
            rec.message_track(rec.fields_get(['state']), initial_values)
        return ret or True
    
    

    def compute_show_edit(self):
        ret = True
        if self.approval_state != 'none':
            approval_category_id = self.approval_id.category_id
            if not approval_category_id:#直接导入的可用记录
                approval_category_id = self.env['approval.category'].search([('res_model','=',self._name)],limit=1)

            edit_group = approval_category_id.edit_group_id
            #wired,直接使用关联字段have_approval_access，结果不同
            if (not self.approval_id.flow_id.have_operation_access and (not edit_group or self.env.user not in edit_group.users)):
                ret = False
        return ret

    
class ApprovalFlowLine(models.Model):
    """ Model for case stages. This models the main stages of a Approval Request management flow. """

    _inherit = 'approval.flow.line'

    name = fields.Char(related=False)
    employee_id = fields.Many2one('hr.employee', required=False)
    condition_employee = fields.Char(u'审批人表达式',help=u"审批人由业务单据相关字段决定")
    condition = fields.Text(u'审批条件')
    action = fields.Char(u'动作',help=u"审批后执行的动作")
    model_id = fields.Many2one('ir.model', related=u'category_id.model_id')
    

class approval(models.Model):
    _inherit = 'approval.approval'

    res_model = fields.Char(u'模型名称', index=True)
    res_id = fields.Integer(u'记录ID', index=True)
    origin = fields.Char(u'源单据', compute='_compute_res_name', store=True)
    employee_id = fields.Many2one('hr.employee', related='flow_id.employee_id', string=u'审批人')
    
    @api.depends('res_model', 'res_id')
    def _compute_res_name(self):
        for record in self:
            if record.res_model and record.res_id:
                res_record = self.env[record.res_model].browse(record.res_id)
                record.origin = res_record.display_name
    
    def _get_eval_context(self):
        eval_context = self.env['ir.actions.act_window']._get_eval_context()
        model = self.env[self.res_model] if self.res_model else self.env[self._name]
        rid = self.res_id if self.res_id else self.id
        record = self.res_model and model.browse(rid) or self
        eval_context.update({
            # orm
            'env': self.env,
            'model': model,
            # record
            'record': record,
            # deprecated use record, records
            'object': record,
            'obj': record,            
        })
        
        return eval_context  
        
    def approval_expr_eval_expr(self, lines, env):
        """
        Evaluate each line of ``lines`` with the ``Environment`` environment, returning
        the value of the last line.
        """
        assert lines, 'You used a NULL action in a workflow, use dummy node instead.'
        result = False
        for line in lines.split('\n'):
            line = line.strip()
            if not line:
                continue
            if line == 'True':
                result = True
            elif line == 'False':
                result = False
            else:
                result = safe_eval(line, env, nocopy=True)

        return result    
        
    @api.multi
    def open_to_approval_object(self):
        self.ensure_one()
        return  {
            'name': _(u'业务单'),
            'type': 'ir.actions.act_window',
            'view_type': 'form',
            'view_mode': 'form',
            'res_model': self.res_model,
            'res_id': self.res_id,
        
        }
    
    @api.multi
    def create_approval_flows(self, force_unlink=False):
        for approval in self:
            already_flowed_line_ids = approval.stage_ids.mapped('line_id')
            flow_ids = self.env['approval.flow']
            flow_to_remove = self.env['approval.flow']

            eval_context = approval._get_eval_context()
            for category_flow in approval.category_id.flow_line_ids:                
                if category_flow not in already_flowed_line_ids:                    
                    if (not category_flow.condition or approval.approval_expr_eval_expr(category_flow.condition, eval_context)):
                        employee_id = category_flow.employee_id or category_flow.condition_employee and approval.approval_expr_eval_expr(category_flow.condition_employee, eval_context)
                        if not employee_id:
                            raise UserError(_(u'审批类型%s中的审批节点%s配置不正确，未找到审批人！'%(approval.category_id.name,category_flow.name)))
                        values = {
                            'approval_id': approval.id,
                            'employee_id': employee_id.id,
                            'state': 'draft',
                            'sequence':category_flow.sequence,
                            'name':category_flow.name,
                            'line_id':category_flow.id,
                        }
                
                        flow_id = self.env['approval.flow'].create(values)
                        flow_ids |= flow_id
                else:
                    flow_id = approval.stage_ids.filtered(lambda stage_id: stage_id.line_id == category_flow)
                    if(category_flow.condition and not approval.approval_expr_eval_expr(category_flow.condition, eval_context)):
                        flow_to_remove |= flow_id                    
                    else:
                        employee_id = category_flow.employee_id or category_flow.condition_employee and approval.approval_expr_eval_expr(category_flow.condition_employee, eval_context)
                        if not employee_id:
                            raise UserError(_(u'审批类型%s中的审批节点%s配置不正确，未找到审批人！'%(approval.category_id.name,category_flow.name)))
                        if flow_id.employee_id != employee_id:
                            flow_id.employee_id = employee_id
                    
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

            flow_to_remove |= approval.stage_ids.filtered(lambda flow: not flow.line_id or flow.line_id not in approval.category_id.flow_line_ids)
            #
            if flow_to_remove:
                flow_to_remove.unlink()
                
            if not approval.stage_ids:
                raise UserError(_(u'审批类型%s中的无审批节点生效，请联系系统管理员！'%(approval.category_id.name)))

        return True
    
    def _can_be_agent_by_cur_user(self):
        ret = False
        if self.owner_user_id.agent_enable and self.owner_user_id.agent_user_id.user_id.id  == self.env.uid:
            agent_stime = self.owner_user_id.agent_start_time
            agent_etime = self.owner_user_id.agent_stop_time
            cur_day = fields.Date.today()
            if (not agent_stime or cur_day >= agent_stime) and (not agent_etime or cur_day <= agent_etime):
                ret = True
        return ret
        
    def _compute_have_access(self):
        uid = self.env.uid
        for record in self:
            if uid == SUPERUSER_ID:                
                record.have_operation_access = True
                continue
            if (record.owner_user_id.user_id.id == uid or record._can_be_agent_by_cur_user()):
                record.have_operation_access = True


    def action_start(self):            
        super(approval,self).action_start()
        self.mapped('stage_ids').write({'approval_date':False,'real_emplyee_id':False})
        

    def action_canceled(self):
        super(approval,self).action_canceled()
        for stage_id in self.stage_ids:
            stage_id.write({'state': 'cancel'})           
            

    def change_flow_id_to_pre(self):
        if not self.category_id.back_pre_step:
            return
        
        if self.flow_id == self.stage_ids[0]:#回退到未审批        
            self.action_reset()
        else:
            self.flow_id.state = 'start'
            pre_flow = self.stage_ids.filtered(lambda r: r.state == 'approved')
            self.flow_id = pre_flow[-1:]
            self.flow_id.state = 'start'

                
class ApprovalFlow(models.Model):
    _inherit = 'approval.flow'
    
    def _can_be_agent_by_cur_user(self):
        ret = False
        if self.employee_id.agent_enable and self.employee_id.agent_user_id.user_id.id  == self.env.uid:
            agent_stime = self.employee_id.agent_start_time
            agent_etime = self.employee_id.agent_stop_time
            cur_day = fields.Date.today()
            if (not agent_stime or cur_day >= agent_stime) and (not agent_etime or cur_day <= agent_etime):
                ret = True
        return ret
    
    def _compute_have_access(self):
        uid = self.env.uid
        for record in self:
            if (record.employee_id.user_id.id == uid or record._can_be_agent_by_cur_user()) and record.approval_id.flow_id.id == record.id:
                record.have_operation_access = True
    
    name = fields.Char(related=False)
    res_model = fields.Char(u'模型名称', related='approval_id.res_model')
    res_id = fields.Integer(u'记录ID', related='approval_id.res_id')
    line_id = fields.Many2one('approval.flow.line',string='模板行')
    approval_date = fields.Datetime(u'审批日期',readonly=True)
    have_operation_access = fields.Boolean(string="have_operation_access",compute=_compute_have_access)
    state = fields.Selection(selection_add=[('cancel', u'取消')])
    #action = fields.Char(u'动作',help=u"审批后执行的动作")
    real_emplyee_id = fields.Many2one('hr.employee', string=u'执行审批人', readonly=True)
    
    def _valite_user_action_auth(self):
        if self.id != self.approval_id.flow_id.id:
            raise UserError('Now, You can not do this operation')
        if self.employee_id.user_id.id != self.env.uid and not self._can_be_agent_by_cur_user():
            raise UserError('You have no access to do this operation')
        return
        
    @api.multi
    def approved(self):
        res = super(ApprovalFlow,self).approved()
        self.approval_date = fields.Datetime.now()
        self.real_emplyee_id = self.env.user.employee_ids[0]
        return res
        
    @api.multi
    def rejected(self):
        res = super(ApprovalFlow,self).rejected()
        self.approval_date = fields.Datetime.now()
        self.real_emplyee_id = self.env.user.employee_ids[0]
        return res
        
    
    @api.multi
    def write(self, vals):
        if 'employee_id' in vals:
            old_employee_id = self.mapped('employee_id')
        super(ApprovalFlow,self).write(vals)
        if 'employee_id' in vals:
            new_employee_id = self.env['hr.employee'].browse(vals['employee_id'])
            self.mapped('approval_id').message_subscribe_users(user_ids=[new_employee_id.user_id.id])
            self.mapped('approval_id').message_post(_(u'审批节点%s的审批人由%s变更为%s'%(self[0].name,old_employee_id.name, new_employee_id.name)))
            
        return True
                
class ApprovalTemplate(models.Model):
    _name = 'approval.template'
    _description = 'Approval Template'
    
    res_model = fields.Char(
        'Related Document Model', required=True, index=True, help='Model of the followed resource')
    name = fields.Char('Category Name', required=True, translate=True)   
    
    