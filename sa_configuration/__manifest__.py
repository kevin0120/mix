# -*- coding: utf-8 -*-
{
    'name': "sa_configuration",

    'summary': """
        Short (1 phrase/line) summary of the module's purpose, used as
        subtitle on modules listing or apps.openerp.com""",

    'description': """
        Long description of module's purpose
    """,

    'author': "My Company",
    'website': "http://www.yourcompany.com",

    # Categories can be used to filter modules in modules listing
    # Check https://github.com/odoo/odoo/blob/master/odoo/addons/base/module/module_data.xml
    # for the full list
    'category': 'Uncategorized',
    'version': '0.1',

    # any module necessary for this one to work correctly
    'depends': ['mrp_maintenance', 'maintenance','mrp'],

    # always loaded
    'data': [
        # 'security/ir.model.access.csv',
        'views/sa_views_menus.xml',
        'views/product_views.xml',
        'views/mrp_routing_view.xml',
        'views/mrp_workorder_views.xml',
        'views/mrp_bom_views.xml',
        'data/maintenance_data.xml',
    ],
    # only loaded in demonstration mode
}