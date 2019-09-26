# -*- coding: utf-8 -*-
{
    'name': "svw_enhanced",

    'summary': """
        Short (1 phrase/line) summary of the module's purpose, used as
        subtitle on modules listing or apps.openerp.com""",

    'description': """
        上海大众
    """,

    'author': "My Company",
    'website': "http://www.yourcompany.com",

    # Categories can be used to filter modules in modules listing
    # Check https://github.com/odoo/odoo/blob/master/odoo/addons/base/module/module_data.xml
    # for the full list
    'category': 'Manufacturing',
    'version': '0.1',

    # any module necessary for this one to work correctly
    'depends': ['sa_base', 'spc', 'wave'],

    # always loaded
    'data': [
        # 'security/ir.model.access.csv',
        'views/mrp_bom_views.xml',
        'views/hide_menu.xml',
        'views/mrp_routing_view.xml',
        'views/mrp_workcenter_views.xml',
        'views/operation_result_views.xml',
    ],
    # only loaded in demonstration mode
}