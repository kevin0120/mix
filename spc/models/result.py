# -*- coding: utf-8 -*-

from odoo import fields,models,api,_
from odoo.exceptions import ValidationError
import odoo.addons.decimal_precision as dp
from datetime import datetime


class OperationResult(models.HyperModel):
    _name = "operation.result"

    workorder_id = fields.Many2one('mrp.workorder', 'Operation')
    workcenter_id = fields.Many2one('mrp.workcenter', related='workorder_id.workcenter_id', store=True, readonly=True)  # TDE: necessary ?
    production_id = fields.Many2one('mrp.production', 'Production Order')

    pset_strategy = fields.Selection([('AD', 'Torque tightening'),
                                         ('AW', 'Angle tightening'),
                                         ('ADW', 'Torque/Angle tightening'),
                                         ('LN', 'Loosening'),
                                         ('AN', 'Number of Pulses tightening'),
                                         ('AT', 'Time tightening')])

    pset_m_max = fields.Float(string='Set Max Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_min = fields.Float(string='Set Min Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_threshold = fields.Float(string='Set Threshold Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_Target = fields.Float(string='Set Target Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_w_max = fields.Float(string='Set Max Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_min = fields.Float(string='Set Min Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_threshold = fields.Float(string='Set Threshold Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_Target = fields.Float(string='Set Target Angel(grad)', digits=dp.get_precision('Quality Tests'))

    cur_file = fields.Char(string='Waveform File')

    name = fields.Char('Name', default=lambda self: _('New'))
    point_id = fields.Many2one('quality.point', 'Control Point')
    quality_state = fields.Selection([
        ('none', 'To do'),
        ('pass', 'Passed'),
        ('fail', 'Failed')], string='Status', default='none', copy=False)

    user_id = fields.Many2one('res.users', 'Responsible')

    product_id = fields.Many2one('product.product', 'Product',
                                 domain="[('type', 'in', ['consu', 'product'])]", required=True)

    measure_torque = fields.Float('Measure Torque(NM)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_degree = fields.Float('Measure Degree(grad)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_t_don = fields.Float('Measure Time Done(ms)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_result = fields.Selection([
        ('none', 'No measure'),
        ('ok', 'OK'),
        ('nok', 'NOK')], string="Measure Success", default='none')

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        if 'display_name' in fields:
            fields.remove('display_name')
        if '__last_update' in fields:
            fields.remove('__last_update')
        return super(OperationResult, self).read(fields, load=load)

    @api.model
    def create(self, vals):
        if 'name' not in vals or vals['name'] == _('New'):
            vals['name'] = self.env['ir.sequence'].next_by_code('quality.check') or _('New')
        return super(OperationResult, self).create(vals)

    @api.multi
    def do_fail(self):
        self.write({
            'quality_state': 'fail',
            'user_id': self.env.user.id,
            'time': datetime.now()})
        return True

    @api.multi
    def do_pass(self):
        self.write({'quality_state': 'pass',
                    'user_id': self.env.user.id,
                    'time': datetime.now()})
        return True

    @api.multi
    def do_measure(self):
        if self.measure_torque < self.point_id.tolerance_min or self.measure_torque > self.point_id.tolerance_max:
            return self.do_fail()
        elif self.measure_degree < self.point_id.tolerance_min_degree or self.measure_degree > self.point_id.tolerance_max_degree:
            return self.do_fail()
        else:
            return self.do_pass()

