# -*- coding: utf-8 -*-
{
    'name': "approval enhance",

    'summary': """
        approval function""",

    'description': """
        make approval flows' generation be conditional
    """,

    'author': "YDIT",
    'website': "http://www.empower.cn",

    # Categories can be used to filter modules in modules listing
    # Check https://github.com/odoo/odoo/blob/master/odoo/addons/base/module/module_data.xml
    # for the full list
    'version': '0.1',

    # any module necessary for this one to work correctly
    'depends': ['approval'],

    # always loaded
    'data': [
        'wizard/approval_wizard_view.xml',
        'views/approval_views.xml',
        'views/approval_form_view_hide_edit_view.xml',
    ],

    'installable': True,
    'application': True,
    'auto_install': False,
    'qweb': [],

}
