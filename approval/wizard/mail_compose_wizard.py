# -*- coding: utf-8 -*-
# -*- coding: utf-8 -*-
# Part of Odoo. See LICENSE file for full copyright and licensing details.

from odoo import api, models,fields


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
                approval.flow_id.state = 'rejected'
                approval.flow_id = None
                approval.state = 'cancel'
                approval.close_date = fields.Datetime.now()
                getattr(related_model, 'callback_action_refuse')() if hasattr(related_model, 'callback_action_refuse') else None    # 回调机制
            elif self._context.get('mark_flow_status') == 'approved':
                approval.flow_id.state = 'approved'
                approval.change_flow_id_by_next()
                if approval.state == 'done':
                    getattr(related_model, 'callback_action_approve')() if hasattr(related_model, 'callback_action_approve') else None  # 回调机制
            self = self.with_context(mail_post_autofollow=True)
        return super(MailComposeMessage, self).send_mail(auto_commit=auto_commit)
