# -*- coding: utf-8 -*-
from odoo import http

# class SvwExhanced(http.Controller):
#     @http.route('/svw_exhanced/svw_exhanced/', auth='public')
#     def index(self, **kw):
#         return "Hello, world"

#     @http.route('/svw_exhanced/svw_exhanced/objects/', auth='public')
#     def list(self, **kw):
#         return http.request.render('svw_exhanced.listing', {
#             'root': '/svw_exhanced/svw_exhanced',
#             'objects': http.request.env['svw_exhanced.svw_exhanced'].search([]),
#         })

#     @http.route('/svw_exhanced/svw_exhanced/objects/<model("svw_exhanced.svw_exhanced"):obj>/', auth='public')
#     def object(self, obj, **kw):
#         return http.request.render('svw_exhanced.object', {
#             'object': obj
#         })