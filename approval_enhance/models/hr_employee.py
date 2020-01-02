# -*- coding: utf-8 -*-

from odoo import api, fields, models, SUPERUSER_ID, _


class HrEmployee(models.Model):
    _inherit = 'hr.employee'

    @api.onchange('agent_enable')
    def _change_agent_enable(self):
        super(HrEmployee,self)._change_agent_enable()
        if not self.agent_enable:
            self.agent_group_ids = []
    
    @api.one    
    @api.depends('agent_user_id.user_id.groups_id','user_id.groups_id')        
    def _compute_agent_groups(self):
        self.agent_group_ids = self.user_id.groups_id - self.agent_user_id.user_id.groups_id
        
        
    def _set_agent_groups(self):
        self.agent_group_ids_store |= self.agent_group_ids

    agent_group_ids = fields.Many2many('res.groups', string=u'代理的权限', compute="_compute_agent_groups", inverse='_set_agent_groups',store=False)
    agent_group_ids_store = fields.Many2many('res.groups', string=u'代理的权限')


