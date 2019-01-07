# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
import odoo.addons.decimal_precision as dp
from odoo.exceptions import UserError
from urlparse import urljoin
import uuid

import requests as Requests
import json
from requests import ConnectionError, RequestException, exceptions
import logging
PUSH_MAINTENANCE_REQ_URL = '/rush/v1/maintenance'

logger = logging.getLogger(__name__)


class MaintenanceCheckPointCategory(models.Model):
    _name = 'maintenance.cp.category'

    name = fields.Char('Check Point Name')
    code = fields.Char('Check Point Code')

    test_type = fields.Selection([
        ('passfail', 'Pass - Fail'),
        ('measure', 'Measure')], string="Test Type",
        default='passfail', required=True)


class MaintenanceCheckPoint(models.Model):
    _name = 'maintenance.cp'
    _description = 'Equipment Checklist for each maintenance request '

    name = fields.Char('CP Name')
    equipment_id = fields.Many2one('maintenance.equipment', string='Equipment', index=True)

    category_id = fields.Many2one('maintenance.cp.category')
    description = fields.Char('Maintenance Check Point Description')

    test_type = fields.Selection(related='category_id.test_type', readonly=True)

    norm = fields.Float('Normal', digits=dp.get_precision('Maintenance Tests'))  # TDE RENAME ?
    tolerance_min = fields.Float('Min Tolerance', digits=dp.get_precision('Maintenance Tests'))
    tolerance_max = fields.Float('Max Tolerance', digits=dp.get_precision('Maintenance Tests'))

    @api.onchange('norm')
    def onchange_norm(self):
        if self.tolerance_max == 0.0:
            self.tolerance_max = self.norm


class MaintenanceCheckPointAction(models.Model):
    _name = 'maintenance.cp.action'

    request_id = fields.Many2one('maintenance.request', string='Request', index=True)

    test_type = fields.Selection([
        ('passfail', 'Pass - Fail'),
        ('measure', 'Measure')], string="Test Type",
        default='passfail', required=True)

    category_id = fields.Many2one('maintenance.cp.category')

    description = fields.Char('Maintenance Check Point Description')

    norm = fields.Float('Normal', digits=dp.get_precision('Maintenance Tests'))  # TDE RENAME ?
    tolerance_min = fields.Float('Min Tolerance', digits=dp.get_precision('Maintenance Tests'))
    tolerance_max = fields.Float('Max Tolerance', digits=dp.get_precision('Maintenance Tests'))

    measure = fields.Float(digits=dp.get_precision('Maintenance Tests'))

    measure_success = fields.Selection([
        ('none', 'No measure'),
        ('pass', 'Pass'),
        ('fail', 'Fail')], string="Measure Success", compute="_compute_measure_success", store=True)

    @api.one
    @api.depends('measure')
    def _compute_measure_success(self):
        if self.test_type == 'passfail':
            self.measure_success = 'none'
        else:
            if self.measure < self.tolerance_min or self.measure > self.tolerance_max:
                self.measure_success = 'fail'
            else:
                self.measure_success = 'pass'

    @api.model
    def bulk_create(self, all_vals):

        all_updates = []
        for vals in all_vals:
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
                # updates.append(('create_date', '%s', '(now() at time zone \'UTC\')'))
                # updates.append(('write_date', '%s', '(now() at time zone \'UTC\')'))
            all_updates.append(updates)
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


class MaintenanceRequest(models.Model):
    _inherit = 'maintenance.request'

    def _default_access_token(self):
        return uuid.uuid4().hex

    access_token = fields.Char('Invitation Token', default=_default_access_token)

    check_point_action_ids = fields.One2many('maintenance.cp.action', 'request_id')

    @api.multi
    def assign_ticket_to_self(self):
        for r in self:
            r.technician_user_id = self.env.uid

    @api.multi
    def post_maintenance_req(self):
        self.ensure_one()
        ret = self.equipment_id
        if not ret:
            return
        if ret.category_name == 'Gun':
            master = ret._get_parent_masterpc()[0]
            if not master:
                return
            connections = master.connection_ids.filtered(lambda r: r.protocol == 'http') if master.connection_ids else None
            if not connections:
                return
            url = ['http://{0}:{1}{2}'.format(connect.ip, connect.port, PUSH_MAINTENANCE_REQ_URL) for connect in
                   connections][0]
            val = {
                "type": ret.category_name,
                "name": ret.display_name,
                "expire_time": fields.Date.today()
            }
            try:
                Requests.post(urljoin(url, PUSH_MAINTENANCE_REQ_URL), data=json.dumps(val), headers={'Content-Type': 'application/json'}, timeout=3)
            except ConnectionError as e:
                logger.debug(u'发送维护请求失败, 错误原因:{0}'.format(e.message))
            except RequestException as e:
                logger.debug(u'发送维护请求失败, 错误原因:{0}'.format(e.message))
            finally:
                return

    @api.model
    def create(self, vals):
        ret = super(MaintenanceRequest, self).create(vals)
        if 'equipment_id' in vals and vals.get('maintenance_type', 'corrective') == 'preventive':
            equipment_id = vals.get('equipment_id')
            check_point_ids = self.env['maintenance.cp'].sudo().search([('equipment_id', '=', equipment_id)])
            actions = []
            for check_point in check_point_ids:
                actions.append({
                    'point_id': check_point.id,
                    'category_id': check_point.category_id.id,
                    'request_id': ret.id,
                    'description': check_point.description if check_point.description else "",
                    'test_type': check_point.category_id.test_type,
                    'norm': check_point.norm,
                    'tolerance_min': check_point.tolerance_min,
                    'tolerance_max': check_point.tolerance_max,
                })
            if len(actions) > 0:
                self.env['maintenance.cp.action'].sudo().bulk_create(actions)

            ret.post_maintenance_req()  # 主动发送维护请求到HMI
            # template_id = self.env.ref('sa_maintenance.new_maintenance_request_email_template', False)
            # if template_id:
            #     rendering_context = dict(self._context)
            #     rendering_context.update({
            #         'dbname': self._cr.dbname,
            #         'base_url': self.env['ir.config_parameter'].sudo().get_param('web.base.url',
            #                                                                      default='http://localhost:8069')
            #     })
            #     template_id = template_id.with_context(rendering_context)
            #     mail_id = template_id.send_mail(ret.id)  # 先不要发送,之后调用send方法发送邮件
            #     current_mail = self.env['mail.mail'].browse(mail_id)
            #     # self.env["celery.task"].call_task("mail.mail", "send_async_by_id", mail_id=mail_id)
            #     current_mail.send()  # 发送邮件
        return ret


class MaintenanceEquipment(models.Model):
    _inherit = 'maintenance.equipment'

    check_point_ids = fields.One2many('maintenance.cp', 'equipment_id')

    @api.one
    def _get_parent_masterpc(self):
        cat = self
        while cat:
            if cat.category_name == 'MasterPC':
                return cat
            cat = cat.parent_id
        return None

