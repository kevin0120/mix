# -*- coding: utf-8 -*-
from odoo import models, api


class HyperModel(models.Model):
    _auto = True  # automatically create database backend
    _register = False  # not visible in ORM registry, meant to be python-inherited only
    _abstract = False  # not abstract
    _transient = False  # not transient
    _order = 'time'  # default order for searching results

    _log_access = False

    _hyper = True


    @api.model
    def _add_magic_fields(self):
        # cyclic import
        from odoo import fields

        # this field 'id' must override any other column or field
        self._add_field('id', fields.Id(automatic=True))

        self._add_field('time', fields.Date())