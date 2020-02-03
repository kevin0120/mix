# -*- coding: utf-8 -*-
{
    'name': "sanyi_enhanced",

    'summary': """
        TS005""",

    'description': """
        TS005
    """,

    'author': "Saturn",
    'website': "http://www.yourcompany.com",

    # Categories can be used to filter modules in modules listing
    # Check https://github.com/odoo/odoo/blob/10.0/odoo/addons/base/module/module_data.xml
    # for the full list
    'category': 'customized',
    'version': '10.0.0.1',

    # any module necessary for this one to work correctly
    'depends': ['sa_base', 'sa_maintenance', 'spc'],

    # always loaded
    'data': [
        'views/mrp_workcenter_views.xml',
        'views/maintenance_views.xml',
        'views/quality_views.xml'
    ]
}