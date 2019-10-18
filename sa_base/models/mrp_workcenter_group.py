from odoo import api, fields, models,_


class MrpWorkcenterGroup(models.Model):
    _name = 'mrp.workcenter.group'
    _description = 'Work Center Group'
    _order = "code"

    code = fields.Char('Reference', copy=False, required=True)
    name = fields.Char('Work Center Group')
    sa_workcenter_ids = fields.Many2many('mrp.workcenter', 'mrp_workcenter_rel', 'group_id', 'workcenter_id',
                                     string="Workcenters", copy=False)

    active = fields.Boolean(
        'Active', default=True,
        help="If the active field is set to False, it will allow you to hide the bills of material without removing it.")

    _sql_constraints = [('code_uniq', 'unique(code)', 'Only one code per Work Center Group is allowed')]

    @api.multi
    def name_get(self):
        res = []
        for center_group in self:
            res.append((center_group.id, _('[%s] %s') % (center_group.code, center_group.name)))
        return res

    @api.model
    def create(self, vals):
        ret = super(MrpWorkcenterGroup, self).create(vals)
        workgroup_id = ret
        for tool in workgroup_id.sa_workcenter_ids.mapped('equipment_ids'):
            val = {
                "workgroup_id": workgroup_id.id,
                "workcenter_id": tool.workcenter_id.id,
                "tool_id": tool.id,
            }
            self.env['mrp.workcenter.group.tool'].sudo().create(val)
        return ret

    @api.multi
    def write(self, vals):
        super(MrpWorkcenterGroup, self).write(vals)
        if 'sa_workcenter_ids' in vals:
            ret = self.env['mrp.workcenter.group.tool'].search([('workgroup_id', 'in', self.ids)])
            ret.sudo().unlink()
            self.ensure_one()
            for tool in self.sa_workcenter_ids.mapped('equipment_ids'):
                val = {
                    "workgroup_id": self.id,
                    "workcenter_id": tool.workcenter_id.id,
                    "tool_id": tool.id,
                }
                self.env['mrp.workcenter.group.tool'].sudo().create(val)
        return True


    @api.multi
    def unlink(self):
        ret =self.env['mrp.workcenter.group.tool'].search([('workgroup_id', 'in', self.ids)])
        ret.sudo().unlink()
        return super(MrpWorkcenterGroup, self).unlink()
    #
    # @api.onchange('sa_workcenter_ids', 'sa_workcenter_ids.equipment_ids')
    # def _onchange_operations(self):
    #     ret =self.env['mrp.workcenter.group.tool'].search([('workgroup_id', 'in', self.ids)])
    #     ret.sudo().unlink()
    #     self.ensure_one()
    #     workgroup_id = self.id
    #     for workcenter in self.sa_workcenter_ids:
    #         for gun in workcenter.mapped('equipment_ids'):
    #             gun_id = gun.id
    #             val = {
    #                 "workgroup_id": workgroup_id,
    #                 "workcenter_id": workcenter.id,
    #                 "tool_id": gun_id,
    #             }
    #             self.env['mrp.workcenter.group.tool'].sudo().create(val)
    #     return


class MrpWorkcenterGroupTool(models.Model):
    _name = 'mrp.workcenter.group.tool'
    _description = 'Work Center Group Tool'
    _order = "id"

    workgroup_id = fields.Many2one('mrp.workcenter.group', string='Work Group', copy=False,
                              required=True)

    workcenter_id = fields.Many2one('mrp.workcenter', string='Work Centre', copy=False,
                              required=True)

    tool_id = fields.Many2one('maintenance.equipment', string='Gun', copy=False,
                              required=True)


    @api.multi
    def name_get(self):
        res = []
        for tool_group in self:
            res.append((tool_group.id, _('[%s]@ %s@%s') % (tool_group.tool_id.code, tool_group.workcenter_id.name,tool_group.workgroup_id.name)))
        return res