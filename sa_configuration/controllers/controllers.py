# -*- coding: utf-8 -*-
from odoo import http

# class SaConfiguration(http.Controller):
#     @http.route('/sa_configuration/sa_configuration/', auth='public')
#     def index(self, **kw):
#         return "Hello, world"

#     @http.route('/sa_configuration/sa_configuration/objects/', auth='public')
#     def list(self, **kw):
#         return http.request.render('sa_configuration.listing', {
#             'root': '/sa_configuration/sa_configuration',
#             'objects': http.request.env['sa_configuration.sa_configuration'].search([]),
#         })

#     @http.route('/sa_configuration/sa_configuration/objects/<model("sa_configuration.sa_configuration"):obj>/', auth='public')
#     def object(self, obj, **kw):
#         return http.request.render('sa_configuration.object', {
#             'object': obj
#         })