# -*- coding: utf-8 -*-
from __future__ import division
from odoo import fields, models, api, _, SUPERUSER_ID
import logging

_logger = logging.getLogger(__name__)


class OperationResult(models.HyperModel):
    _inherit = "operation.result"

    def init(self):
        self.env.cr.execute("""
            CREATE OR REPLACE FUNCTION create_operation_result(pset_m_threshold numeric, pset_m_max numeric,
                                                   control_data timestamp without time zone, pset_w_max numeric,
                                                   user_id bigint, one_time_pass boolean,
                                                   pset_strategy varchar, measure_result varchar,
                                                   pset_w_threshold numeric, cur_objects varchar,
                                                   pset_m_target numeric, pset_m_min numeric,
                                                   final_pass varchar, measure_degree numeric,
                                                   measure_t_done numeric, measure_torque numeric, op_time integer,
                                                   pset_w_min numeric, pset_w_target numeric, lacking varchar,
                                                   quality_state varchar, exception_reason varchar, sent boolean,
                                                   batch varchar,
                                                   order_id bigint, nut_no varchar, r_tightening_id integer,
                                                   vin_code varchar, vehicle_type varchar, gun_sn varchar)
  returns BIGINT as
$$
DECLARE
  result_id            bigint;
  r_vin_code           varchar;
  r_job                varchar;
  r_one_time_pass      varchar = 'fail';
  r_qcp_id             BIGINT  = null;
  consu_bom_id         BIGINT  = null;
  r_consu_product_id   BIGINT  = null;
  r_production_id      BIGINT  = null;
  r_workcenter_id      BIGINT;
  r_gun_id             BIGINT;
  r_product_id         BIGINT;
  r_program_id         BIGINT;
  r_assembly_id        BIGINT;
  r_bom_line_id        BIGINT  = null;
  r_operation_point_id BIGINT  = null;
  r_measure_result     varchar;
BEGIN
  select me.id into r_gun_id
  from public.maintenance_equipment me
  where me.serial_no = gun_sn; /*永远获取真实的拧紧工具*/
  case pset_strategy
    when 'LN'
      then r_measure_result = 'lsn';
    ELSE r_measure_result = measure_result;
    end case;
  if one_time_pass
  then
    r_one_time_pass = 'pass';
  end if;
  if order_id != 0
  then
    select mp.track_no,
           qp.id,
           co.id,
           mp.id,
           wo.workcenter_id,
           mp.product_id,
           co.program_id,
           co.product_id,
           mp.assembly_line_id,
           mbl.id,
           mbl.operation_point_id
           into r_vin_code, r_qcp_id, consu_bom_id, r_production_id, r_workcenter_id, r_product_id, r_program_id, r_consu_product_id, r_assembly_id,r_bom_line_id,r_operation_point_id
    from public.mrp_workorder wo,
         public.mrp_wo_consu co,
         public.sa_quality_point qp,
         public.mrp_production mp,
         public.product_product pp,
         public.maintenance_equipment me,
         public.mrp_bom_line mbl
    where wo.id = order_id
      and co.workorder_id = order_id
      and pp.default_code = nut_no
      and co.bom_line_id = qp.bom_line_id
      and wo.production_id = mp.id
      and mbl.id = co.bom_line_id
      and me.serial_no = gun_sn
      and co.product_id = pp.id;
  else
    r_vin_code = vin_code;
    order_id = null;
    r_bom_line_id = null;
    select pp2.id,
           tp.cou_pid,
           tp.cou_program_code,
           tp.cou_program_id,
           tp.cou_workcenter_id,
           tp.cou_qcp_id,
           tp.operation_point_id
           into r_product_id, r_consu_product_id,r_job, r_program_id, r_workcenter_id,r_qcp_id,r_operation_point_id
    from (select sqp.product_id   cou_pid,
                 sqp.operation_id cou_opd,
                 op.routing_id    cou_routing_id,
                 sqp.id           cou_qcp_id,
                 opp.id           operation_point_id,
                 op.workcenter_id cou_workcenter_id,
                 sa_program.id    cou_program_id,
                 sa_program.code  cou_program_code
          from public.product_product pp,
               public.operation_point opp,
               public.sa_quality_point sqp,
               public.mrp_routing_workcenter op,
               public.controller_program sa_program,
               public.mrp_routing routing
          where pp.id = sqp.product_id
            and opp.qcp_id = sqp.id /* 找到拧紧点 */
            and sqp.name = nut_no /* 找到质量控制点 */
            and op.id = sqp.operation_id
            and sa_program.id = sqp.program_id
         ) as tp, /* 找到拧紧点相关拧紧编号的清单 */
         public.product_product pp2,
         public.mrp_bom bom
    where pp2.default_code = vehicle_type
      and bom.product_id = pp2.id; /* 找到物料清单 */
  end if;
  INSERT INTO public.operation_result (pset_m_threshold, pset_m_max, control_date, pset_w_max, user_id,
                                       one_time_pass, pset_strategy, measure_result, pset_w_threshold,
                                       cur_objects, pset_m_target, pset_m_min, final_pass, measure_degree,
                                       measure_t_don,
                                       measure_torque, op_time, pset_w_min, pset_w_target, lacking, quality_state,
                                       exception_reason, sent, batch, track_no,
                                       qcp_id, bom_line_id, operation_point_id, workorder_id,
                                       consu_product_id, consu_bom_line_id,
                                       production_id, tool_id, program_id, product_id, assembly_line_id, workcenter_id,
                                       tightening_id, job,
                                       time)
  VALUES (pset_m_threshold, pset_m_max, control_data, pset_w_max, user_id, r_one_time_pass, pset_strategy,
          r_measure_result,
          pset_w_threshold, cur_objects, pset_m_target, pset_m_min, final_pass, measure_degree,
          measure_t_done, measure_torque, op_time,
          pset_w_min, pset_w_target, lacking,
          quality_state, exception_reason,
          sent, batch, r_vin_code,
          r_qcp_id, r_bom_line_id, r_operation_point_id, order_id,
          r_consu_product_id,
          consu_bom_id, r_production_id,
          r_gun_id, r_program_id, r_product_id, r_assembly_id,
          r_workcenter_id, r_tightening_id, r_job,
          NOW());
  result_id = lastval();
  RETURN result_id;
END;
$$
  LANGUAGE plpgsql;
            """)
