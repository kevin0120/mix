# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import UserError

import pyecharts
from pyecharts import Bar , Pie


class ProductTemplate(models.Model):
    _inherit = 'product.template'

    sa_type = fields.Selection([('screw', 'Screw'), ('vehicle', 'Vechile')], default='vehicle', string='产品类型')


class ProductProduct(models.Model):
    _inherit = 'product.product'

    sa_type = fields.Selection(related='product_tmpl_id.sa_type')
    screw_type_code = fields.Char(string='螺栓编号', copy=False)
    vehicle_type_code = fields.Char(string="车型编码", copy=False)
    qp_count = fields.Integer(string='Quality Point Count', compute='_compute_product_quality_point_count')
    # active_bom_id = fields.Many2one('mrp.bom', string='Current Actived BOM', compute='_compute_actived_bom_id')

    # active_bom_line_ids = fields.One2many('mrp.bom.line', related='active_bom_id.bom_line_ids')

    description = fields.Text(string='Product Description')

    _sql_constraints = [
        ('screw_type_code_uniq', 'unique(screw_type_code)', _("A Screw Code can only be assigned to one SCREW !")),
        ('vehicle_type_code_uniq', 'unique(vehicle_type_code)', _("A Vehicle Code can only be assigned to one Vehicle !")),
    ]

    def _compute_product_quality_point_count(self):
        for product in self:
            product.qp_count = self.env['quality.point'].search_count([('product_id','=', product.id)])

    @api.multi
    def copy(self, default=None):
        raise UserError(_('Product can not be copy by User!'))

    @api.multi
    @api.depends('name', 'screw_type_code', 'vehicle_type_code')
    def name_get(self):
        res = []
        for product in self:
            if product.sa_type == 'vehicle':
                name = u"[{0}] {1}".format(product.vehicle_type_code, product.name)
            else:
                name = u"[{0}] {1}".format(product.screw_type_code, product.name)
            res.append((product.id, name))
        return res

    # @api.multi
    # @api.depends('bom_ids')
    # def _compute_actived_bom_id(self):
    #     for product in self:
    #         product.active_bom_id = product.bom_ids.filtered("active")[0] if product.bom_ids else False

    ###此方法打开相应的页面
    @api.multi
    def action_sa_view_bom(self):
        action = self.env.ref('sa_configuration.sa_product_open_bom').read()[0]
        template_ids = self.mapped('product_tmpl_id').ids
        # bom specific to this variant or global to template
        action['context'] = {
            'default_product_tmpl_id': template_ids[0],
            'default_product_id': self.ids[0],
        }
        action['domain'] = ['|', ('product_id', 'in', [self.ids]), '&', ('product_id', '=', False),
                            ('product_tmpl_id', 'in', template_ids)]
        return action


