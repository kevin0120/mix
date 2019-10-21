# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError
import logging
import json

_logger = logging.getLogger(__name__)


class MrpWOConsu(models.Model):
    _name = 'mrp.wo.consu.line'
    _inherits = {'sa.quality.check': 'check_id'}
    _order = "sequence"
    _log_access = False

    sequence = fields.Integer(string='Sequence')

    workorder_id = fields.Many2one('mrp.workorder')

    bom_line_id = fields.Many2one('mrp.bom.line')

    qcp_id = fields.Many2one('sa.quality.point', related='check_id.point_id')
    check_id = fields.Many2one('sa.quality.check', ondelete='restrict', required=True,
                               string='Quality Check For Work Step')

    product_id = fields.Many2one('product.product', related='check_id.product_id', string='Consume Product')

    product_qty = fields.Float('Consume Product Qty')

    tool_id = fields.Many2one('maintenance.equipment', string='Screw Gun', copy=False)

    program_id = fields.Many2one('controller.program')

    @api.model
    def _bulk_create(self, all_vals):
        # low-level implementation of create()
        if self.is_transient():
            self._transient_vacuum()

        all_updates = []
        for vals in all_vals:
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


class MrpWorkorder(models.Model):
    _inherit = 'mrp.workorder'

    track_no = fields.Char('Finished Product Tracking Number', related='production_id.track_no')

    consu_bom_line_ids = fields.One2many('mrp.wo.consu.line', 'workorder_id', string='Consume Product')

    @api.multi
    def unlink(self):
        raise ValidationError(u'不允许删除工单')


class DispatchingWorkOrder(models.Model):
    _name = 'dispatch.mrp.workorder'
    _order = "sequence"
    _log_access = False

    @api.depends('workcenter_ids')
    def _compute_operation_workcenter_domain(self):
        for dispatch_wo in self:
            if dispatch_wo.workcenter_id not in dispatch_wo.workcenter_ids:
                dispatch_wo.workcenter_id = False
            if not dispatch_wo.workcenter_ids:
                dispatch_wo.operation_workcenter_domain = json.dumps([])
            dispatch_wo.operation_workcenter_domain = json.dumps([('id', 'in', dispatch_wo.workcenter_ids.ids)])

    sequence = fields.Integer(string='Sequence', default=1)

    is_dispatched = fields.Boolean('Is Dispatched', default=False)

    production_id = fields.Many2one('mrp.production', string='Manufacture Order', index=True)

    workorder_id = fields.Many2one('mrp.workorder')

    user_id = fields.Many2one('res.users', string='Responsible Person')

    routing_id = fields.Many2one('mrp.routing', 'Routing', related='production_id.routing_id', store=True,
                                 readonly=True)

    operation_id = fields.Many2one('mrp.routing.workcenter')

    workcenter_ids = fields.Many2many('mrp.workcenter', related='operation_id.workcenter_ids',
                                      copy=False, readonly=True)
    workcenter_id = fields.Many2one('mrp.workcenter', string='Operate In Work Center')

    operation_workcenter_domain = fields.Char(
        compute=_compute_operation_workcenter_domain,
        readonly=True,
        default=json.dumps([]),
        store=False,
    )

    @api.onchange('workcenter_id')
    def _onchange_workcenter_id(self):
        self.ensure_one()
        if not self.workcenter_id:
            self.user_id = False
        self.user_id = self.workcenter_id.user_ids and self.workcenter_id.user_ids[0]  # 第一个用户


class MrpProduction(models.Model):
    _inherit = 'mrp.production'

    product_uom_id = fields.Many2one(
        'product.uom', 'Product Unit of Measure',
        oldname='product_uom', readonly=True, required=False,
        states={'draft': [('readonly', False)]})  # 重新定义产品单位字段

    state = fields.Selection([
        ('draft', 'Draft'),
        ('confirmed', 'Confirmed'),
        ('planned', 'Planned'),
        ('progress', 'In Progress'),
        ('done', 'Done'),
        ('cancel', 'Cancelled')], string='State',
        copy=False, default='draft')  # 重新定义生产订单的阶段，将草稿阶段移至最前

    track_no = fields.Char('Finished Product Tracking Number')

    dispatch_workorder_ids = fields.One2many('dispatch.mrp.workorder', 'production_id', string='Dispatching Work Order')

    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line ID')

    # 重写此方法，禁止生成库存移动
    @api.multi
    def _generate_moves(self):
        return True

    @api.multi
    @api.depends('bom_id.routing_id', 'bom_id.routing_id.sa_operation_ids')
    def _compute_routing(self):
        for production in self:
            if production.bom_id.routing_id.sa_operation_ids:
                production.routing_id = production.bom_id.routing_id.id
            else:
                production.routing_id = False

    @api.model
    def create(self, vals):
        ret = super(MrpProduction, self).create(vals)
        if ret:
            ret._create_dispatch_workorder()
        return ret

    @api.multi
    def do_confirm(self):
        can_confirm_productions = self.env['mrp.production']
        for production in self:
            # todo: 确认是否可以确认需求
            can_confirm_productions |= production
            production.dispatch_workorder_ids.write({'is_dispatched': True})
        can_confirm_productions.write({'state': 'confirm'})

    @api.multi
    def button_dispatching(self):
        for production in self:
            if not production.routing_id:
                _logger.error("Production: {0} Can Not Create Dispatching Info. Cause It Is Not Define Routing".format(
                    production.name))
                continue

    @api.multi
    def _create_dispatch_workorder(self):
        ret = True
        self.ensure_one()
        routing_id = self.routing_id
        if not routing_id:
            _logger.error(
                "Production: {0} Can Not Create Dispatching Work Order. Cause It Is Not Define Routing".format(
                    self.name))
            return ret
        for idx, operation_id in enumerate(routing_id.sa_operation_ids):
            val = {
                'sequence': idx,
                'operation_id': operation_id.id,
                'routing_id': routing_id.id,
                'production_id': self.id,
                'workcenter_id': operation_id.workcenter_id and operation_id.workcenter_id.id,  # 设置优先选择工位，方便后续快速排产
                'user_id': operation_id.workcenter_id.user_ids and operation_id.workcenter_id.user_ids[0].id
            }
            self.env['dispatch.mrp.workorder'].sudo().create(val)  # 创建派工信息

        return ret
