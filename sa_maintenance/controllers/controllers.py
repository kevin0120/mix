# -*- coding: utf-8 -*-
from odoo import http

# class SaMaintenance(http.Controller):
#     @http.route('/sa_maintenance/sa_maintenance/', auth='public')
#     def index(self, **kw):
#         return "Hello, world"

#     @http.route('/sa_maintenance/sa_maintenance/objects/', auth='public')
#     def list(self, **kw):
#         return http.request.render('sa_maintenance.listing', {
#             'root': '/sa_maintenance/sa_maintenance',
#             'objects': http.request.env['sa_maintenance.sa_maintenance'].search([]),
#         })

#     @http.route('/sa_maintenance/sa_maintenance/objects/<model("sa_maintenance.sa_maintenance"):obj>/', auth='public')
#     def object(self, obj, **kw):
#         return http.request.render('sa_maintenance.object', {
#             'object': obj
#         })