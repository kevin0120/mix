# -*- coding: utf-8 -*-
from odoo import http

# class Spc(http.Controller):
#     @http.route('/spc/spc/', auth='public')
#     def index(self, **kw):
#         return "Hello, world"

#     @http.route('/spc/spc/objects/', auth='public')
#     def list(self, **kw):
#         return http.request.render('spc.listing', {
#             'root': '/spc/spc',
#             'objects': http.request.env['spc.spc'].search([]),
#         })

#     @http.route('/spc/spc/objects/<model("spc.spc"):obj>/', auth='public')
#     def object(self, obj, **kw):
#         return http.request.render('spc.object', {
#             'object': obj
#         })