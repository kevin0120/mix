# -*- coding: utf-8 -*-

from odoo import fields,models,api,_
from odoo.exceptions import ValidationError
import odoo.addons.decimal_precision as dp
from datetime import datetime, timedelta, date
from dateutil.relativedelta import relativedelta
import json


class OperationResult(models.HyperModel):
    _name = "operation.result"

    workorder_id = fields.Many2one('mrp.workorder', 'Operation')
    workcenter_id = fields.Many2one('mrp.workcenter', related='workorder_id.workcenter_id', store=True, readonly=True)  # TDE: necessary ?
    production_id = fields.Many2one('mrp.production', 'Production Order')

    pset_configuration = fields.Char(string='Pset Configuration')

    pset_strategy = fields.Selection([('AD', 'Torque tightening'),
                                         ('AW', 'Angle tightening'),
                                         ('ADW', 'Torque/Angle tightening'),
                                         ('LN', 'Loosening'),
                                         ('AN', 'Number of Pulses tightening'),
                                         ('AT', 'Time tightening')])

    pset_m_max = fields.Float(string='Set Max Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_min = fields.Float(string='Set Min Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_threshold = fields.Float(string='Set Threshold Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_m_target = fields.Float(string='Set Target Torque(NM)', digits=dp.get_precision('Quality Tests'))

    pset_w_max = fields.Float(string='Set Max Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_min = fields.Float(string='Set Min Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_threshold = fields.Float(string='Set Threshold Angel(grad)', digits=dp.get_precision('Quality Tests'))

    pset_w_target = fields.Float(string='Set Target Angel(grad)', digits=dp.get_precision('Quality Tests'))

    cur_objects = fields.Char(string='Waveform Files')

    name = fields.Char('Name', default=lambda self: _('New'))
    point_id = fields.Many2one('quality.point', 'Control Point')
    quality_state = fields.Selection([
        ('none', 'To do'),
        ('pass', 'Passed'),
        ('fail', 'Failed')], string='Status', default='none', copy=False)

    user_id = fields.Many2one('res.users', 'Responsible')

    product_id = fields.Many2one('product.product', 'Vehicle',
                                 domain="[('sa_type', '=', 'vehicle')]", required=True)

    consu_product_id = fields.Many2one('product.product', 'Screw',
                                 domain="[('sa_type', '=', 'screw')]", required=True)

    control_date = fields.Datetime('Control Date')

    measure_torque = fields.Float('Measure Torque(NM)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_degree = fields.Float('Measure Degree(grad)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_t_don = fields.Float('Measure Time Done(ms)', default=0.0, digits=dp.get_precision('Quality Tests'))

    measure_result = fields.Selection([
        ('none', 'No measure'),
        ('ok', 'OK'),
        ('nok', 'NOK')], string="Measure Success", default='none')

    lacking = fields.Selection([('lack', 'Data Lacking'),
        ('normal', 'Normal')], string='Lacking', default='lack', compute='_compute_result_lacking', store=True)

    op_time = fields.Integer(string=u'第几次拧紧作业', default=1)

    one_time_pass = fields.Selection([('pass', 'One Time Passed'),
        ('fail', 'Failed')], string='One Time Pass', default='fail',
                                   compute='_compute_result_pass', store=True)

    final_pass = fields.Selection([('pass', 'Final Passed'),
        ('fail', 'Failed')], string='Final Pass', default='fail',
                                compute='_compute_result_pass', store=True)

    @api.one
    @api.constrains('cur_objects')
    def _constraint_curs(self):
        try:
            obj = json.loads(self.cur_objects)
            if not isinstance(list, obj):
                raise ValidationError('Waveform File is not a list')
        except ValueError:
            raise ValidationError('Waveform File is not Schema correct')

    @api.multi
    @api.depends('measure_result', 'op_time')
    def _compute_result_pass(self):
        for result in self:
            if result.measure_result != 'ok':
                result.one_time_pass = 'fail'
                result.final_pass = 'fail'
                return
            result.final_pass = 'pass'
            if result.op_time == 1:
                result.one_time_pass = 'pass'

    @api.multi
    @api.depends('measure_torque', 'measure_degree', 'measure_result')
    def _compute_result_lacking(self):
        for result in self:
            if result.measure_torque not in [False, 0.0] and result.measure_degree not in [False, 0.0] and result.measure_result != 'none':
                result.lacking = 'normal'
            result.lacking = 'lack'

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
    def write(self, vals):
        if 'cur_objects' in vals:
            cur_objects = list(set(json.loads(vals['cur_objects'])))  # 去重
            vals.update({'cur_objects': json.dumps(cur_objects)})
        ret = super(OperationResult, self).write(vals)
        return ret

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

