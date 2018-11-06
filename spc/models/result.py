# -*- coding: utf-8 -*-

from odoo import fields,models,api,_
from odoo.exceptions import ValidationError
import odoo.addons.decimal_precision as dp
from datetime import datetime, timedelta, date
from dateutil.relativedelta import relativedelta
from odoo.addons.spc.controllers.result import _post_aiis_result_package
import json
import logging

from collections import defaultdict, MutableMapping, OrderedDict
from odoo.tools import frozendict

_logger = logging.getLogger(__name__)


class OperationResult(models.HyperModel):
    _name = "operation.result"

    sequence = fields.Integer('sequence', default=1)

    workorder_id = fields.Many2one('mrp.workorder', 'Operation')
    workcenter_id = fields.Many2one('mrp.workcenter',  readonly=True)
    production_id = fields.Many2one('mrp.production', 'Production Order')

    vin = fields.Char('Vin')  # 如果无法通过工单查询,使用此字段

    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line', readonly=True)

    qcp_id = fields.Many2one('quality.point')

    gun_id = fields.Many2one('maintenance.equipment', string='Screw Gun', copy=False)

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
        ('exception', 'Exception'),
        ('pass', 'Passed'),
        ('fail', 'Failed')], string='Status', default='none', copy=False)

    exception_reason = fields.Char('Exception Reason')

    user_id = fields.Many2one('res.users', 'Responsible')

    product_id = fields.Many2one('product.product', 'Vehicle',
                                 domain="[('sa_type', '=', 'vehicle')]")

    consu_bom_line_id = fields.Many2one('mrp.wo.consu')

    consu_product_id = fields.Many2one('product.product')

    control_date = fields.Datetime('Control Date')

    measure_torque = fields.Float('Measure Torque(NM)', default=0.0, digits=dp.get_precision('Operation Result'))

    measure_degree = fields.Float('Measure Degree(grad)', default=0.0, digits=dp.get_precision('Operation Result'))

    measure_t_don = fields.Float('Measure Time Done(ms)', default=0.0, digits=dp.get_precision('Operation Result'))

    measure_result = fields.Selection([
        (_('none'), 'No measure'),
        (_('ok'), 'OK'),
        (_('nok'), 'NOK')], string="Measure Success", default='none')

    lacking = fields.Selection([(_('lack'), _('Data Lacking')),
        (_('normal'), _('Normal'))], string='Lacking', default='lack')

    op_time = fields.Integer(string=u'第几次拧紧作业', default=1)

    one_time_pass = fields.Selection([('pass', _('One Time Passed')),
        ('fail', _('Failed'))], string='One Time Pass', default='fail')

    final_pass = fields.Selection([('pass', _('Final Passed')),
        ('fail', _('Failed'))], string='Final Pass', default='fail')

    sent = fields.Boolean('Have sent to aiis', default=False)

    batch = fields.Char('Batch')  # 批次信息

    @api.multi
    def sent_aiis(self):
        results = self.filtered(lambda r: r.measure_result in ['ok', 'nok'])
        if not results:
            return True
        aiis_urls = self.env['ir.config_parameter'].sudo().get_param('aiis.urls')
        if not aiis_urls:
            return
        aiis_urls = aiis_urls.split(',')
        ret = _post_aiis_result_package(aiis_urls, results)
        return True

    # @api.one
    # @api.constrains('cur_objects')
    # def _constraint_curs(self):
    #     try:
    #         obj = json.loads(self.cur_objects)
    #         if not isinstance(obj, list):
    #             raise ValidationError('Waveform File is not a list')
    #     except ValueError:
    #         raise ValidationError('Waveform File is not Schema correct')

    # @api.multi
    # @api.depends('measure_result', 'op_time')
    # def _compute_result_pass(self):
    #     for result in self:
    #         result.one_time_pass = 'fail'
    #         result.final_pass = 'fail'
    #         if result.measure_result != 'ok':
    #             continue
    #         result.final_pass = 'pass'
    #         if result.op_time == 1:
    #             result.one_time_pass = 'pass'

    # @api.multi
    # # @api.depends('measure_torque', 'measure_degree', 'measure_result')
    # @api.depends('measure_result')
    # def _compute_result_lacking(self):
    #     for result in self:
    #         result.lacking = 'lack'
    #         if result.measure_result != 'none': # if result.measure_torque not in [False, 0.0] and result.measure_degree not in [False, 0.0] and result.measure_result != 'none':
    #             result.lacking = 'normal'

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
    def _bulk_write(self, all_vals):
        all_updates = []
        update_wo_ids = []
        for vals in all_vals:
            updates = []  # list of (column, expr) or (column, pattern, value)
            splite_updates = []
            upd_todo = []  # list of column names to set explicitly
            updend = []  # list of possibly inherited field names
            direct = []  # list of direcly updated columns
            has_trans = self.env.lang and self.env.lang != 'en_US'
            single_lang = len(self.env['res.lang'].get_installed()) <= 1
            for name, val in vals.iteritems():
                field = self._fields[name]
                if field and field.deprecated:
                    _logger.warning('Field %s.%s is deprecated: %s', self._name, name, field.deprecated)
                if field.store:
                    if hasattr(field, 'selection') and val:
                        self._check_selection_field_value(name, val)
                    if field.column_type:
                        # if single_lang or not (has_trans and field.translate is True):
                        #     # val is not a translation: update the table
                        #     val = field.convert_to_column(val, self)
                        updates.append((name, field.column_format, '''timestamp '%s' ''' % val if field.column_type[0] == 'timestamp' else val))
                        if name != 'id':
                            splite_updates.append((name, field.column_format, '''timestamp '%s' ''' % val if field.column_type[0] == 'timestamp' else val))
                        direct.append(name)
                    else:
                        upd_todo.append(name)
                else:
                    updend.append(name)

            if self._log_access:
                updates.append(('write_uid', '%s', self._uid))
                updates.append(('write_date', "(now() at time zone 'UTC')"))
                direct.append('write_uid')
                direct.append('write_date')
            all_updates.append(updates)
            update_wo_ids.append(splite_updates)

        ### 添加需要重计算字段
        self.modified([u[0] for u in update_wo_ids[0]])
        cr = self._cr
        t = [tuple(u[2] for u in update if len(u) > 2) for update in all_updates]
        x = []
        for _t in t:
            s = '(%s)' % (','.join("'%s'" % _s if (isinstance(_s, str) or isinstance(_s, unicode)) and _s.find('timestamp') < 0 else '%s' % str(_s) for _s in _t))
            x.append(s)
        query = """UPDATE "%s" AS o SET (%s) = (%s) FROM(VALUES %s) AS s (%s) WHERE o.id = s.id""" % (
            self._table,
            ', '.join('%s' % u[0] for u in update_wo_ids[0]),
            ','.join("s.%s" % u[0] for u in update_wo_ids[0]),
            ','.join('%s'% _t for _t in x),
            ', '.join('%s' % u[0] for u in all_updates[0]),
        )

        cr.execute(query)
        map(lambda vals: self.browse(vals['id'])._validate_fields(vals), all_vals)

        # recompute new-style fields
        # if self.env.recompute and self._context.get('recompute', True):
        #     self._bulk_recompute()

        return True

    @api.model
    def _bulk_recompute(self):
        """ Recompute stored function fields. The fields and records to
            recompute have been determined by method :meth:`modified`.
        """
        while self.env.has_todo():
            field, recs = self.env.get_todo()
            # determine the fields to recompute
            fs = self.env[field.model_name]._field_computed[field]
            ns = [f.name for f in fs if f.store]
            # evaluate fields, and group record ids by update
            updates = defaultdict(set)
            for rec in recs.exists():
                vals = rec._convert_to_write({n: rec[n] for n in ns})
                updates[frozendict(vals)].add(rec.id)
            # update records in batch when possible
            with recs.env.norecompute():
                for vals, ids in updates.iteritems():
                    recs.browse(ids)._write(dict(vals))
            # mark computed fields as done
            map(recs._recompute_done, fs)

    @api.multi
    def bulk_write(self, vals):
        if not self:
            return True
        self._check_concurrency()
        self._bulk_write(vals)
        return True

    @api.multi
    def do_fail(self):
        self.write({
            'quality_state': 'fail',
            'user_id': self.env.user.id,
            'time': datetime.now()})
        return True

    @api.multi
    def do_exception(self):
        self.write({
            'quality_state': 'exception',
            'user_id': self.env.user.id,
            'time': datetime.now()})
        return True

    @api.multi
    def get_torques(self, args, limit=1000):
        query = self._where_calc(args)
        self._apply_ir_rules(query, 'read')
        from_clause, where_clause, where_clause_params = query.get_sql()

        where_str = where_clause and (" WHERE %s" % where_clause) or ''

        limit_str = limit and ' limit %d' % limit or ''
        query_str = 'SELECT "%s".measure_torque FROM ' % self._table + from_clause + where_str + limit_str
        self._cr.execute(query_str, where_clause_params)
        res = self._cr.fetchall()

        # TDE note: with auto_join, we could have several lines about the same result
        # i.e. a lead with several unread messages; we uniquify the result using
        # a fast way to do it while preserving order (http://www.peterbe.com/plog/uniqifiers-benchmark)
        def _uniquify_list(seq):
            seen = set()
            return [x for x in seq if x not in seen and not seen.add(x)]

        return _uniquify_list([x[0] for x in res])

    @api.multi
    def get_angles(self, args, limit=1000):
        query = self._where_calc(args)
        self._apply_ir_rules(query, 'read')
        from_clause, where_clause, where_clause_params = query.get_sql()

        where_str = where_clause and (" WHERE %s" % where_clause) or ''

        limit_str = limit and ' limit %d' % limit or ''
        query_str = 'SELECT "%s".measure_degree FROM ' % self._table + from_clause + where_str + limit_str
        self._cr.execute(query_str, where_clause_params)
        res = self._cr.fetchall()

        # TDE note: with auto_join, we could have several lines about the same result
        # i.e. a lead with several unread messages; we uniquify the result using
        # a fast way to do it while preserving order (http://www.peterbe.com/plog/uniqifiers-benchmark)
        def _uniquify_list(seq):
            seen = set()
            return [x for x in seq if x not in seen and not seen.add(x)]

        return _uniquify_list([x[0] for x in res])

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

