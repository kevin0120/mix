# -*- coding: utf-8 -*-

from odoo import models, fields,api, _
from odoo.exceptions import ValidationError
from odoo.osv import expression
from validators import ipv4, ValidationFailure


class EquipmentConnection(models.Model):
    _name = 'maintenance.equipment.connection'

    ip = fields.Char(string='IP')
    tty = fields.Char(string='Serial TTY')

    equipment_id = fields.Many2one('Equipment')

    port = fields.Integer(string='port')
    unitid = fields.Integer(string='Unit ID', help='Modbus need this ID for identification')
    protocol = fields.Selection([('modbustcp','ModbusTCP'),('modbusrtu', 'ModbusRTU'),
                                 ('rawtcp','TCP'),('rawudp','UDP')], string='Protocol')

    @api.constrains('ip')
    def _constraint_ip(self):
        try:
            ipv4(self.ip)
        except ValidationFailure:
            raise ValidationError(_('is NOT valid IP Address!'))

    @api.multi
    def name_get(self):
        def get_names(cat):
            if cat.protocol == 'modbustcp':
                return u"modbustcp://{0}:{1}/{2}".format(cat.ip, cat.port, cat.unitid)
            if cat.protocol == 'modbusrtu':
                return u"modbusrtu://{0}/{1}".format(cat.tty, cat.unitid)
            if cat.protocol == 'rawtcp':
                return u"tcp://{0}:{1}".format(cat.ip, cat.port)
            if cat.protocol == 'modbustcp':
                return u"udp://{0}:{1}".format(cat.ip, cat.port)

        return [(cat.id, get_names(cat)) for cat in self]


class MaintenanceEquipment(models.Model):
    _inherit = 'maintenance.equipment'
    _parent_name = "parent_id"
    _parent_store = True
    _order = 'parent_left'
    _parent_order = 'name'

    parent_id = fields.Many2one('maintenance.equipment', 'Parent Equipment', index=True, ondelete='cascade')
    child_ids = fields.One2many('maintenance.equipment', 'parent_id', 'Child Equipments')

    parent_left = fields.Integer('Left Parent', index=1)
    parent_right = fields.Integer('Right Parent', index=1)
    #
    child_equipments_count = fields.Integer(compute='_compute_child_equipments_count')

    category_name = fields.Char(compute='_compute_categroy_name', default='', store=True)

    connection_ids = fields.One2many('maintenance.equipment.connection', 'equipment_id', 'Connection Information')




    @api.multi
    def _compute_child_equipments_count(self):
        for equipment in self:
            equipment.child_equipments_count = len(equipment.child_ids)

    @api.multi
    @api.depends('category_id')
    def _compute_categroy_name(self):
        for equipment in self:
            equipment.category_name = equipment.category_id.name or ''

    @api.constrains('parent_id')
    def _check_category_recursion(self):
        if not self._check_recursion():
            raise ValidationError(_('Error ! You cannot create recursive Equipments.'))
        return True

    @api.multi
    def name_get(self):
        def get_names(cat):
            """ Return the list [cat.name, cat.parent_id.name, ...] """
            res = []
            while cat:
                if cat.name and cat.serial_no:
                    res.append(cat.name + '#' + cat.serial_no)
                if cat.name and not cat.serial_no:
                    res.append(cat.name)
                cat = cat.parent_id
            return res

        return [(cat.id, " / ".join(reversed(get_names(cat)))) for cat in self]

    @api.model
    def name_search(self, name, args=None, operator='ilike', limit=100):
        if not args:
            args = []
        if name:
            # Be sure name_search is symetric to name_get
            category_names = name.split(' / ')
            parents = list(category_names)
            child = parents.pop()
            domain = [('name', operator, child)]
            if parents:
                names_ids = self.name_search(' / '.join(parents), args=args, operator='ilike', limit=limit)
                category_ids = [name_id[0] for name_id in names_ids]
                if operator in expression.NEGATIVE_TERM_OPERATORS:
                    categories = self.search([('id', 'not in', category_ids)])
                    domain = expression.OR([[('parent_id', 'in', categories.ids)], domain])
                else:
                    domain = expression.AND([[('parent_id', 'in', category_ids)], domain])
                for i in range(1, len(category_names)):
                    domain = [[('name', operator, ' / '.join(category_names[-1 - i:]))], domain]
                    if operator in expression.NEGATIVE_TERM_OPERATORS:
                        domain = expression.AND(domain)
                    else:
                        domain = expression.OR(domain)
            categories = self.search(expression.AND([domain, args]), limit=limit)
        else:
            categories = self.search(args, limit=limit)
        return categories.name_get()


