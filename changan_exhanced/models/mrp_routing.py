# -*- coding: utf-8 -*-

from odoo import models, fields, api

import json

import requests as Requests

from requests import ConnectionError, RequestException

# MASTER_WROKORDERS_API = '/rush/v1/mrp.routing.workcenter'
#
#
# class MrpRoutingWorkcenter(models.Model):
#     _inherit = 'mrp.routing.workcenter'
#
#