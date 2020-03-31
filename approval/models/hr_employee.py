# -*- coding: utf-8 -*-

from odoo import api, fields, models, SUPERUSER_ID, _


class HrEmployee(models.Model):
    _inherit = 'hr.employee'

    @api.onchange('agent_enable')
    def _change_agent_enable(self):
        if self.agent_enable:
            self.agent_start_time = fields.Date.today()

    flow_ids = fields.One2many('approval.flow', 'employee_id', string='Approvals')
    agent_user_id = fields.Many2one('hr.employee', string='Agent', require=True)
    agent_enable = fields.Boolean(string="Enable Agent", default=False, help="Enable Agent.")
    agent_start_time = fields.Date(string='Agent Start Time')
    agent_stop_time = fields.Date(string='Agent Stop Time', help="Set the agent stop time.If not set"
                                                                 ", It will agent forever")
    todo_approval_count = fields.Integer('Todo Approval Request')
    favorite_request_count = fields.Integer('Launched Approval')

    @api.multi
    def toggle_agent_enable(self):
        self.ensure_one()
        self.agent_enable = not self.agent_enable
        return {'url': '/web', 'type': 'ir.actions.act_url', 'target': 'self'}

