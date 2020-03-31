# -*- coding: utf-8 -*-
# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import api, models,fields
from odoo.tools.safe_eval import safe_eval, test_python_expr

class MailComposeMessage(models.TransientModel):
    _inherit = 'mail.compose.message'

    @api.multi
    def send_mail(self, auto_commit=False):
        if self._context.get('default_model') == 'approval.approval' and self._context.get('default_res_id') and self._context.get('mark_flow_status'):
            approval = self.env['approval.approval'].browse([self._context['default_res_id']])
            related_model = None
            if approval.res_model and approval.res_id:
                related_model = self.env[approval.res_model].browse([approval.res_id])
            if self._context.get('mark_flow_status') == 'rejected':
                if not approval.category_id.back_pre_step:
                    approval.state = 'done'
                    approval.flow_id.state = 'rejected'
                    approval.close_date = fields.Datetime.now()
                else:
                    approval.change_flow_id_to_pre()
                getattr(related_model, 'callback_action_refuse')() if hasattr(related_model, 'callback_action_refuse') else None    # 回调机制
            elif self._context.get('mark_flow_status') == 'approved':
                approval.flow_id.state = 'approved'
                approval.change_flow_id_by_next()
                if approval.state == 'done':
                    getattr(related_model, 'callback_action_approve')() if hasattr(related_model, 'callback_action_approve') else None  # 回调机制
            if approval.state == 'done' and approval.category_id.finish_action: # 无论是拒绝还是全部同意都会执行finish_action
                safe_eval(approval.category_id.finish_action.strip(), approval._get_eval_context(), mode="exec", nocopy=True)
            self = self.with_context(mail_post_autofollow=True)
            
        return super(MailComposeMessage,self).send_mail(auto_commit=auto_commit)