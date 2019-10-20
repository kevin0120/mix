# -*- coding: utf-8 -*-

from odoo import api, fields, models, _


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
        if not ret or 'sa_workcenter_ids' not in vals:
            return ret
        workgroup_id = ret
        tool_ids = workgroup_id.sa_workcenter_ids.mapped('equipment_ids')
        for tool in tool_ids:
            val = {
                "workgroup_id": workgroup_id.id,
                "workcenter_id": tool.workcenter_id.id,
                "tool_id": tool.id,
            }
            self.env['mrp.workcenter.group.tool'].sudo().create(val)
        return ret

    @api.multi
    def _update_create_workcenter_group_tool(self):
        tool_category_ids = self.env.ref('sa_base.equipment_Gun') + self.env.ref('sa_base.equipment_Wrench')
        for wg in self:
            already_workcenter_ids = self.env['mrp.workcenter.group.tool'].search([('workgroup_id', '=', wg.id)]).mapped('workcenter_id')
            workcenter_ids = wg.sa_workcenter_ids.filtered(lambda wc: wc not in already_workcenter_ids)
            if not workcenter_ids:
                continue
            tool_ids = self.env['maintenance.equipment'].search(
                [('category_id', 'in', tool_category_ids.ids), ('workcenter_id', 'in', workcenter_ids.ids)])
            for tool in tool_ids:
                val = {
                    "workgroup_id": wg.id,
                    "workcenter_id": tool.workcenter_id.id,
                    "tool_id": tool.id,
                }
                self.env['mrp.workcenter.group.tool'].sudo().create(val)

    @api.multi
    def _update_unlink_workcenter_group_tool(self):
        need_unlink_recs = self.env['mrp.workcenter.group.tool']
        for wg in self:
            recs = self.env['mrp.workcenter.group.tool'].search([('workgroup_id', '=', wg.id),
                                                                ('workcenter_id', 'not in', wg.sa_workcenter_ids.ids)])
            need_unlink_recs |= recs
        need_unlink_recs.sudo.unlink()

    @api.multi
    def write(self, vals):
        ret = super(MrpWorkcenterGroup, self).write(vals)
        if 'sa_workcenter_ids' not in vals:
            return ret
        self._update_create_workcenter_group_tool()
        # self._cr.commit()  # FIXME: 确保事务被提交
        self._update_unlink_workcenter_group_tool()
        return ret

    # @api.multi
    # def unlink(self):
    #     ret = self.env['mrp.workcenter.group.tool'].search([('workgroup_id', 'in', self.ids)])
    #     ret.sudo().unlink()
    #     return super(MrpWorkcenterGroup, self).unlink()


class MrpWorkcenterGroupTool(models.Model):
    _name = 'mrp.workcenter.group.tool'
    _description = 'Work Center Group Tool'
    _order = "id"

    workgroup_id = fields.Many2one('mrp.workcenter.group', string='Work Group', copy=False,
                                   ondelete='cascade', required=True)

    workcenter_id = fields.Many2one('mrp.workcenter', string='Work Centre', copy=False,
                                    ondelete='cascade', required=True)

    tool_id = fields.Many2one('maintenance.equipment', string='Tightening Tool(Tightening Gun/Wrench)', copy=False,
                              ondelete='cascade', required=True)

    @api.multi
    def name_get(self):
        res = []
        for tool_group_id in self:
            res.append((tool_group_id.id, _('[%s]@%s@%s') % (
                tool_group_id.tool_id.serial_no, tool_group_id.workcenter_id.name, tool_group_id.workgroup_id.name)))
        return res
