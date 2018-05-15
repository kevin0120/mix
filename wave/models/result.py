# -*- coding: utf-8 -*-


from odoo import fields,models,api,_
from odoo.exceptions import ValidationError


class OperationResultLine(models.TransientModel):
    _name = "operation.result.line"

    wizard_id = fields.Many2one('wave.wave', ondelete='cascade')

    # result_id = fields.Many2one('operation.result', ondelete='')

    selected = fields.Boolean('Is selected to create Waveform', help='Whether add this result to create waveform')

    workorder_id = fields.Many2one('mrp.workorder')
    workcenter_id = fields.Many2one('mrp.workcenter')  # TDE: necessary ?
    # production_id = fields.Many2one('mrp.production', 'Production Order')
    #
    # cur_objects = fields.Char(string='Waveform Files')
    #
    # user_id = fields.Many2one('res.users', 'Responsible')
    #
    product_id = fields.Many2one('product.product', 'Vehicle')
    #
    consu_product_id = fields.Many2one('product.product', 'Screw' )
    #
    # control_date = fields.Datetime('Control Date')
    #
    # measure_torque = fields.Float('Measure Torque(NM)', default=0.0, digits=dp.get_precision('Operation Result'))
    #
    # measure_degree = fields.Float('Measure Degree(grad)', default=0.0, digits=dp.get_precision('Operation Result'))
    #
    # measure_t_don = fields.Float('Measure Time Done(ms)', default=0.0, digits=dp.get_precision('Operation Result'))
    #
    measure_result = fields.Selection([
        ('none', 'No measure'),
        ('ok', 'OK'),
        ('nok', 'NOK')], string="Measure Result", default='none')
    #
    # lacking = fields.Selection([('lack', 'Data Lacking'),
    #     ('normal', 'Normal')], string='Lacking', default='lack', compute='_compute_result_lacking', store=True)
    #
    # op_time = fields.Integer(string=u'第几次拧紧作业', default=1)
    #
    # one_time_pass = fields.Selection([('pass', 'One Time Passed'),
    #     ('fail', 'Failed')], string='One Time Pass', default='fail',
    #                                compute='_compute_result_pass', store=True)
    #
    # final_pass = fields.Selection([('pass', 'Final Passed'),
    #     ('fail', 'Failed')], string='Final Pass', default='fail',
    #                             compute='_compute_result_pass', store=True)