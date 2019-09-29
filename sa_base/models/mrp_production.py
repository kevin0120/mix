# -*- coding: utf-8 -*-

from odoo import models, fields, api, _
from odoo.exceptions import ValidationError
import math, json
import logging

_logger = logging.getLogger(__name__)


class MrpWOConsu(models.Model):
    _name = 'mrp.wo.consu'

    _order = "sequence"

    _log_access = False

    sequence = fields.Integer(string='Sequence')

    workorder_id = fields.Many2one('mrp.workorder')

    bom_line_id = fields.Many2one('mrp.bom.line')

    product_id = fields.Many2one('product.product', string='Consume Product')
    qty = fields.Float('Consume Product Qty')

    gun_id = fields.Many2one('maintenance.equipment', string='Screw Gun', copy=False)

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

    # vin = fields.Char(string='VIN', related='production_id.vin', copy=False, store=True)

    consu_bom_line_ids = fields.One2many('mrp.wo.consu', 'workorder_id', string='Consume Product')

    @api.multi
    def unlink(self):
        raise ValidationError(u'不允许删除工单')


class MrpProduction(models.Model):
    _inherit = 'mrp.production'

    equipment_name = fields.Char(string='Equipment Name')
    factory_name = fields.Char(string='Factory Name')
    year = fields.Integer(string='Year')
    production_routings = fields.Char(string='Production Routings')
    pin = fields.Integer(string='PIN')
    vin = fields.Char(string='VIN', copy=False)
    pin_check_code = fields.Integer(string='PIN check Code')
    assembly_line_id = fields.Many2one('mrp.assemblyline', string='Assembly Line ID')
    lnr = fields.Char(string='Line Number')
    knr = fields.Char(string='KNR', store=True, compute='_compute_long_pin', copy=False)
    long_pin = fields.Char(string='LongPIN', store=True, compute='_compute_long_pin', copy=False)
    display_long_pin = fields.Char(string='Display LongPIN', store=True, compute='_compute_long_pin')

    _sql_constraints = [('vin_uniq', 'unique(vin)', 'Only one VIN per MO is allowed'),
                        ('pin_check_uniq', 'unique(pin,pin_check_code)', 'Only one KNR per MO is allowed')]

    @api.depends('year', 'factory_name', 'pin', 'pin_check_code')
    def _compute_long_pin(self):
        for mo in self:
            mo.long_pin = u'{0}{1}{2}{3}'.format(mo.factory_name, mo.year, mo.pin, mo.pin_check_code)
            mo.knr = u'{0}{1}'.format(mo.pin, mo.pin_check_code)
            mo.display_long_pin = u'{0}-{1}-{2}={3}'.format(mo.factory_name, mo.year, mo.pin, mo.pin_check_code)

    @api.constrains('year')
    def _constraint_mo_year(self):
        if len(str(self.year)) != 4:
            raise ValidationError(u'不是年份')

    @api.multi
    def _generate_moves(self):  ### 直接返回，不创建调拨单
        return True

    @api.multi
    def plan_by_prs(self):
        """ Create work orders. And probably do stuff, like things. """
        orders_to_plan = self.filtered(lambda order: order.routing_id and order.state == 'confirmed')
        for order in orders_to_plan:
            quantity = order.product_uom_id._compute_quantity(order.product_qty,
                                                              order.bom_id.product_uom_id) / order.bom_id.product_qty
            boms, lines = order.bom_id.explode(order.product_id, quantity, picking_type=order.bom_id.picking_type_id)
            order._generate_workorders_by_prs(boms)
        # for production in self:
        #     if not production.workorder_ids.mapped('check_ids'):
        #         production.workorder_ids._create_checks()
        return orders_to_plan.write({'state': 'planned'})

    @api.multi
    def _generate_workorders_by_prs(self, exploded_boms):
        workorders = self.env['mrp.workorder']
        for bom, bom_data in exploded_boms:
            # If the routing of the parent BoM and phantom BoM are the same, don't recreate work orders, but use one master routing
            if bom.routing_id.id and (
                    not bom_data['parent_line'] or bom_data['parent_line'].bom_id.routing_id.id != bom.routing_id.id):
                workorders += self._workorders_create_by_prs(bom, bom_data)
        return workorders

    def _workorders_create_by_prs(self, bom, bom_data):
        """
        :param bom: in case of recursive boms: we could create work orders for child
                    BoMs
        """
        workorders = self.env['mrp.workorder']
        consume_bom_lines = self.env['mrp.wo.consu']

        consume_bom_line_vals = []
        bom_qty = bom_data['qty']

        # Initial qty producing
        if self.product_id.tracking == 'serial':
            quantity = 1.0
        else:
            quantity = self.product_qty - sum(self.move_finished_ids.mapped('quantity_done'))
            quantity = quantity if (quantity > 0) else 0

        need_plan_prs = self.env['mrp.routing.workcenter']
        for pr in json.loads(self.production_routings):
            need_plan_prs += self.env['mrp.routing.workcenter'].search(
                [('routing_id', '=', bom.routing_id.id), ('group_id.code', '=', pr['pr_group']),
                 ('name', '=', pr['pr_value'])])

        for operation in need_plan_prs:
            # create workorder
            # cycle_number = math.ceil(bom_qty / operation.workcenter_id.capacity)  # TODO: float_round UP
            # duration_expected = (operation.workcenter_id.time_start +
            #                      operation.workcenter_id.time_stop +
            #                      cycle_number * operation.time_cycle * 100.0 / operation.workcenter_id.time_efficiency)
            match_bom_line_ids = bom.bom_line_ids.filtered(lambda r: r.operation_id == operation)
            workorder = workorders.create({
                'name': operation.name,
                'production_id': self.id,
                'workcenter_id': operation.workcenter_id.id,
                'operation_id': operation.id,
                'date_planned_start': self.date_planned_start,
                # 'duration_expected': duration_expected,
                'state': len(workorders) == 0 and 'ready' or 'pending',
                'qty_producing': quantity,
                'capacity': operation.workcenter_id.capacity,
            })
            for idx, line_id in enumerate(match_bom_line_ids):
                val = {
                    'sequence': idx + 1,
                    'workorder_id': workorder.id,
                    'product_id': line_id.product_id.id or None,
                    'qty': line_id.product_qty or None,
                    'gun_id': line_id.gun_id.id,
                    'program_id': line_id.program_id.id,
                    'bom_line_id': line_id.id
                }
                consume_bom_line_vals.append(val)

            if workorders:
                workorders[-1].next_work_order_id = workorder.id
            workorders += workorder

            # assign moves; last operation receive all unassigned moves (which case ?)
            # moves_raw = self.move_raw_ids.filtered(lambda move: move.operation_id == operation)
            # if len(workorders) == len(bom.routing_id.operation_ids):
            #     moves_raw |= self.move_raw_ids.filtered(lambda move: not move.operation_id)
            # moves_finished = self.move_finished_ids.filtered(
            #     lambda move: move.operation_id == operation)  # TODO: code does nothing, unless maybe by_products?
            # moves_raw.mapped('move_lot_ids').write({'workorder_id': workorder.id})
            # (moves_finished + moves_raw).write({'workorder_id': workorder.id})
            #
            # workorder._generate_lot_ids()
        if len(consume_bom_line_vals):
            consume_bom_lines._bulk_create(consume_bom_line_vals)
        return workorders

    def _workorders_create(self, bom, bom_data):
        workorders = self.env['mrp.workorder']
        consume_bom_lines = self.env['mrp.wo.consu']
        bom_qty = bom_data['qty']

        consume_bom_line_vals = []

        # Initial qty producing
        if self.product_id.tracking == 'serial':
            quantity = 1.0
        else:
            quantity = self.product_qty - sum(self.move_finished_ids.mapped('quantity_done'))
            quantity = quantity if (quantity > 0) else 0

        for operation in bom.routing_id.operation_ids:
            # create workorder
            cycle_number = math.ceil(bom_qty / operation.workcenter_id.capacity)  # TODO: float_round UP
            duration_expected = (operation.workcenter_id.time_start +
                                 operation.workcenter_id.time_stop +
                                 cycle_number * operation.time_cycle * 100.0 / operation.workcenter_id.time_efficiency)
            match_bom_line_ids = bom.bom_line_ids.filtered(lambda r: r.operation_id == operation)
            workorder = workorders.create({
                'name': operation.name,
                'production_id': self.id,
                'workcenter_id': operation.workcenter_id.id,
                'operation_id': operation.id,
                'duration_expected': duration_expected,
                'state': len(workorders) == 0 and 'ready' or 'pending',
                'qty_producing': quantity,
                'capacity': operation.workcenter_id.capacity,
            })
            for idx, line_id in enumerate(match_bom_line_ids):
                val = {
                    'sequence': idx + 1,
                    'workorder_id': workorder.id,
                    'product_id': line_id.product_id.id or None,
                    'qty': line_id.product_qty or None,
                    'gun_id': line_id.gun_id.id,
                    'program_id': line_id.program_id.id,
                    'bom_line_id': line_id.id
                }
                consume_bom_line_vals.append(val)
            if workorders:
                workorders[-1].next_work_order_id = workorder.id
            workorders += workorder

            # assign moves; last operation receive all unassigned moves (which case ?)
            moves_raw = self.move_raw_ids.filtered(lambda move: move.operation_id == operation)
            if len(workorders) == len(bom.routing_id.operation_ids):
                moves_raw |= self.move_raw_ids.filtered(lambda move: not move.operation_id)
            moves_finished = self.move_finished_ids.filtered(
                lambda move: move.operation_id == operation)  # TODO: code does nothing, unless maybe by_products?
            moves_raw.mapped('move_lot_ids').write({'workorder_id': workorder.id})
            (moves_finished + moves_raw).write({'workorder_id': workorder.id})

            workorder._generate_lot_ids()

        if len(consume_bom_line_vals):
            consume_bom_lines._bulk_create(consume_bom_line_vals)
        return workorders

    @api.multi
    def unlink(self):
        raise ValidationError(u'不允许删除生产订单')
