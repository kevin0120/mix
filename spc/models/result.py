# -*- coding: utf-8 -*-

from odoo import fields,models,api,_
from odoo.exceptions import ValidationError
import odoo.addons.decimal_precision as dp
from datetime import datetime, timedelta, date
from dateutil.relativedelta import relativedelta
from odoo.addons.spc.controllers.result import _post_aiis_result_package
import json
import logging

_logger = logging.getLogger(__name__)



class OperationResult(models.HyperModel):
    _name = "operation.result"

    workorder_id = fields.Many2one('mrp.workorder', 'Operation')
    workcenter_id = fields.Many2one('mrp.workcenter',  readonly=True)
    production_id = fields.Many2one('mrp.production', 'Production Order')

    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line', readonly=True)

    pset_strategy = fields.Selection([('AD', 'Torque tightening'),
                                         ('AW', 'Angle tightening'),
                                         ('ADW', 'Torque/Angle tightening'),
                                         ('LN', 'Loosening'),
                                         ('AN', 'Number of Pulses tightening'),
                                         ('AT', 'Time tightening')])

    pset_m_max = fields.Float(string='Set Max Torque(NM)', digits=dp.get_precision('Operation Result'))

    pset_m_min = fields.Float(string='Set Min Torque(NM)', digits=dp.get_precision('Operation Result'))

    pset_m_threshold = fields.Float(string='Set Threshold Torque(NM)', digits=dp.get_precision('Operation Result'))

    pset_m_target = fields.Float(string='Set Target Torque(NM)', digits=dp.get_precision('Operation Result'))

    pset_w_max = fields.Float(string='Set Max Angel(grad)', digits=dp.get_precision('Operation Result'))

    pset_w_min = fields.Float(string='Set Min Angel(grad)', digits=dp.get_precision('Operation Result'))

    pset_w_threshold = fields.Float(string='Set Threshold Angel(grad)', digits=dp.get_precision('Operation Result'))

    pset_w_target = fields.Float(string='Set Target Angel(grad)', digits=dp.get_precision('Operation Result'))

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

    measure_torque = fields.Float('Measure Torque(NM)', default=0.0, digits=dp.get_precision('Operation Result'))

    measure_degree = fields.Float('Measure Degree(grad)', default=0.0, digits=dp.get_precision('Operation Result'))

    measure_t_don = fields.Float('Measure Time Done(ms)', default=0.0, digits=dp.get_precision('Operation Result'))

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

    sent = fields.Boolean('Have sent to aiis', default=False)

    @api.multi
    def sent_aiis(self):
        domain = [('measure_result', 'in', ['ok', 'nok'])]
        results = self.search(domain)
        if not results:
            return True
        aiis_urls = self.env['ir.config_parameter'].sudo().get_param('aiis.urls').split(',')
        ret = _post_aiis_result_package(aiis_urls, results)
        return True

    @api.one
    @api.constrains('cur_objects')
    def _constraint_curs(self):
        try:
            obj = json.loads(self.cur_objects)
            if not isinstance(obj, list):
                raise ValidationError('Waveform File is not a list')
        except ValueError:
            raise ValidationError('Waveform File is not Schema correct')

    @api.multi
    @api.depends('measure_result', 'op_time')
    def _compute_result_pass(self):
        for result in self:
            result.one_time_pass = 'fail'
            result.final_pass = 'fail'
            if result.measure_result != 'ok':
                return
            result.final_pass = 'pass'
            if result.op_time == 1:
                result.one_time_pass = 'pass'

    @api.multi
    # @api.depends('measure_torque', 'measure_degree', 'measure_result')
    @api.depends('measure_result')
    def _compute_result_lacking(self):
        for result in self:
            result.lacking = 'lack'
            if result.measure_result != 'none': # if result.measure_torque not in [False, 0.0] and result.measure_degree not in [False, 0.0] and result.measure_result != 'none':
                result.lacking = 'normal'

    @api.multi
    def read(self, fields=None, load='_classic_read'):
        if not fields:
            return super(OperationResult, self).read(fields, load=load)
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

    @api.model
    def _bulk_create(self, all_vals):
        # low-level implementation of create()
        if self.is_transient():
            self._transient_vacuum()

        all_updates = []
        for vals in all_vals:
            if 'name' not in vals or vals['name'] == _('New'):
                vals['name'] = self.env['ir.sequence'].next_by_code('quality.check') or _('New')
            # data of parent records to create or update, by model
            tocreate = {
                parent_model: {'id': vals.pop(parent_field, None)}
                for parent_model, parent_field in self._inherits.iteritems()
            }

            # list of column assignments defined as tuples like:
            #   (column_name, format_string, column_value)
            #   (column_name, sql_formula)
            # Those tuples will be used by the string formatting for the INSERT
            # statement below.
            updates = [
                ('id', "%s", "nextval('%s')" % self._sequence),
            ]

            upd_todo = []
            unknown_fields = []
            protected_fields = []
            for name, val in vals.items():
                field = self._fields.get(name)
                if not field:
                    unknown_fields.append(name)
                    del vals[name]
                elif field.inherited:
                    tocreate[field.related_field.model_name][name] = val
                    del vals[name]
                elif not field.store:
                    del vals[name]
                elif field.inverse:
                    protected_fields.append(field)
            if unknown_fields:
                _logger.warning('No such field(s) in model %s: %s.', self._name, ', '.join(unknown_fields))

            # set boolean fields to False by default (to make search more powerful)
            for name, field in self._fields.iteritems():
                if field.type == 'boolean' and field.store and name not in vals:
                    vals[name] = False

            # determine SQL values
            for name, val in vals.iteritems():
                field = self._fields[name]
                if field.store and field.column_type:
                    updates.append((name, field.column_format, field.convert_to_column(val, self)))
                else:
                    upd_todo.append(name)

                if hasattr(field, 'selection') and val:
                    self._check_selection_field_value(name, val)

            if self._log_access:
                updates.append(('create_uid', '%s', self._uid))
                updates.append(('write_uid', '%s', self._uid))
                updates.append(('create_date', "(now() at time zone 'UTC')"))
                updates.append(('write_date', "(now() at time zone 'UTC')"))
            all_updates.append(updates)
            # insert a row for this record
        cr = self._cr
        t = [tuple(u[2] for u in update if len(u) > 2) for update in all_updates]
        query = """INSERT INTO "%s" (%s) VALUES %s RETURNING id""" % (
            self._table,
            ', '.join('"%s"' % u[0] for u in all_updates[0]),
            ','.join("(nextval('%s')," % self._sequence + str(_t[1:])[1:] for _t in t),
        )

        cr.execute(query)

        # from now on, self is the new record
        ids_news = cr.fetchall()
        return [ids[0] for ids in ids_news]

    @api.model
    def bulk_create(self, vals):
        vals = [self._add_missing_default_values(val)for val in vals]
        self._bulk_create(vals)

    @api.multi
    def write(self, vals):
        # if 'cur_objects' in vals:
        #     cur_objects = list(set(json.loads(vals['cur_objects'])))  # 去重
        #     vals.update({'cur_objects': json.dumps(cur_objects)})
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

