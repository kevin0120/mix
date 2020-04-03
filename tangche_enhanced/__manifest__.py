# -*- coding: utf-8 -*-
{
    'name': "tangche_enhanced",

    'summary': """
        TS002""",

    'description': """
        TS002
    """,

    'author': "My Company",
    'website': "http://www.yourcompany.com",

    # Categories can be used to filter modules in modules listing
    # Check https://github.com/odoo/odoo/blob/10.0/odoo/addons/base/module/module_data.xml
    # for the full list
    'category': 'customize',
    'version': '0.1',

    # any module necessary for this one to work correctly
    'depends': ['api', 'spc', 'common_sa_utils'],

    # always loaded
    'data': [
        # 'security/ir.model.access.csv',
        'views/maintenance_views.xml',
        'views/mrp_workcenter.xml',
        'views/mrp_production_views.xml',
        'views/mrp_workorder_views.xml',
        'views/product_views.xml',
        'views/quality_views.xml',
        'views/mrp_routing_views.xml',
        'views/menu_hide_views.xml',
        # 'demo/demo.xml',
        'data/quality_data.xml',
        'data/mes_data.xml',
    ],
    # only loaded in demonstration mode
    'demo': [
        # 'demo/demo.xml',
    ],
}
