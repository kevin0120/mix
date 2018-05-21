# -*- coding: utf-8 -*-

from odoo import models, fields, api,_
from odoo.exceptions import ValidationError


class MrpBom(models.Model):
    _inherit = 'mrp.bom'

    @api.onchange('routing_id','product_id')
    def _onchange_routing_id(self):
        self.code = u'[{0}]{1}'.format(self.routing_id.name, self.product_id.name)

    @api.constrains('product_id', 'product_tmpl_id')
    def _product_tmpl_product_constraint(self):
        if self.product_id.product_tmpl_id.id != self.product_tmpl_id.id:
            raise ValidationError(_(u'The product template "%s" is invalid on product with name "%s"') % (self.product_tmpl_id.name, self.product_id.name))

    @api.constrains('product_id','routing_id','active')
    def _constraint_active_product_routing(self):
        if not self.active:
            return
            ###只有激活状态才检查
        count = self.env['mrp.bom'].search_count([('id','!=',self.id),('product_id','=',self.product_id.id),('routing_id','=',self.routing_id.id),('active','=',True)])
        if count:
            raise ValidationError(_(u'The product had a related routing config "%s" been actived!') % (self.product_id.name))

    @api.constrains('product_tmpl_id', 'routing_id', 'active')
    def _constraint_active_product_tmpl_routing(self):
        if not self.active:
            return
            ###只有激活状态才检查
        count = self.env['mrp.bom'].search_count(
            [('id','!=',self.id),('product_tmpl_id', '=', self.product_tmpl_id.id), ('routing_id', '=', self.routing_id.id), ('active', '=', True)])
        if count:
            raise ValidationError(
                _(u'The product Template had a related routing config "%s" been actived!') % (self.product_tmpl_id.name))


class MrpBomLine(models.Model):
    _inherit = 'mrp.bom.line'

    group_id = fields.Many2one('mrp.routing.group', related="operation_id.group_id", string='Routing Group')

    program_id = fields.Many2one('controller.program',  related="operation_id.program_id", string='程序号')

    workcenter_id = fields.Many2one('mrp.workcenter', related="operation_id.workcenter_id", string='Work Center')

    _sql_constraints = [
        ('unique_operation_bom_id', 'unique(bom_id,operation_id)', 'Every Bom unique operation'),
    ]

    @api.model
    def create(self, vals):
        line = super(MrpBomLine, self).create(vals)
        vals = {
            'product_id': line.bom_id.product_id.id,
            'product_tmpl_id': line.bom_id.product_tmpl_id.id,
            'operation_id': line.operation_id.id,
            'picking_type_id':
                self.env['stock.picking.type'].search_read(domain=[('code', '=', 'mrp_operation')], fields=['id'],
                                                           limit=1)[0]['id'],
            'workcenter_id': line.operation_id.workcenter_id.id,
            'times': line.product_qty,
            'test_type': 'measure',
        }
        self.env['quality.point'].sudo().create(vals)
        return line

    @api.multi
    def write(self, vals):
        res = super(MrpBomLine, self).write(vals)
        if 'product_qty' in vals:
            for line in self:
                rec = self.env['quality.point'].search([('product_id', '=', line.bom_id.product_id.id),
                                                        ('product_tmpl_id', '=', line.bom_id.product_tmpl_id.id),
                                                        ('operation_id', '=', line.operation_id.id)])
                rec.sudo().write({'times': line.product_qty})
        return res


    @api.multi
    def unlink(self):
        quality_points = self.env['quality.point']
        for line in self:
            rec = self.env['quality.point'].search([('product_id', '=', line.bom_id.product_id.id),
                                                    ('product_tmpl_id', '=', line.bom_id.product_tmpl_id.id),
                                                    ('operation_id', '=', line.operation_id.id)])
            quality_points += rec
        quality_points.sudo().unlink()
        ret = super(MrpBomLine, self).unlink()
        return ret
