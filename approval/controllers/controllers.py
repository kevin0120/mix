# -*- coding: utf-8 -*-
import werkzeug

from odoo.api import Environment
import odoo.http as http

from odoo.http import request
from odoo import SUPERUSER_ID
from odoo import registry as registry_get


class ApprovalController(http.Controller):

    # @http.route('/approval/approved', type='http', auth="none")
    # def approve(self,  db, token, action, id, **kwargs):
    #     registry = registry_get(db)
    #     with registry.cursor() as cr:
    #         env = Environment(cr, SUPERUSER_ID, {})
    #         flow = env['approval.flow'].search([('access_token', '=', token), ('state', '!=', 'approved')])
    #         if flow:
    #             return self.view()
    #     return self.view(db, token, action, id, view='form')
    #
    # @http.route('/approval/rejected', type='http', auth="none")
    # def reject(self, db, token, action, id, **kwargs):
    #     registry = registry_get(db)
    #     with registry.cursor() as cr:
    #         env = Environment(cr, SUPERUSER_ID, {})
    #         flow = env['approval.flow'].search([('access_token', '=', token), ('state', '!=', 'rejected')])
    #         if flow:
    #             return flow.do_rejected()
    #     return self.view(db, token, action, id, view='form')

    @http.route('/approval/view', type='http', auth="user")
    def view(self, db, token, action, active_id, view='form'):
        registry = registry_get(db)
        with registry.cursor() as cr:
            # Since we are in auth=none, create an env with SUPERUSER_ID
            env = Environment(cr, SUPERUSER_ID, {})
            flow = env['approval.flow'].search([('access_token', '=', token)])
            approval = flow.approval_id

            # If user is logged, redirect to form view of event
            # otherwise, display the simplifyed web page with event informations
            if request.session.uid:
                return werkzeug.utils.redirect('/web?db=%s#id=%s&view_type=form&model=approval.approval&action=%s&active_id=%s' % (db, approval.id,action,active_id))
            else:
                return werkzeug.utils.redirect('/web/login')