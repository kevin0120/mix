# -*- coding: utf-8 -*-

from odoo import fields, models, api, _
from odoo.exceptions import UserError

class ModelApprovalWizard(models.TransientModel):
    _name = 'model.approval.wizard'
    _decription = u'选择审批流程'
    
    approval_category = fields.Many2one('approval.category', string=u'审批流程', required=True, domain=lambda self:[('res_model', '=', self.env.context.get('active_model','none'))])
    
    
    @api.multi
    def action_confirm(self):
        self.ensure_one()
        record_id = False
        res_model = self.env.context.get('active_model',False)
        res_id = self.env.context.get('active_id',False)
        
        if res_model and res_id:
            record_id = self.env[res_model].browse(res_id)
            
        if not record_id or not hasattr(record_id, 'action_submit_approval'):
            raise UserError(_('无效操作或单据不支持审批！'))
        
        record_id.action_submit_approval()
        return {'type': 'ir.actions.act_window_close'}