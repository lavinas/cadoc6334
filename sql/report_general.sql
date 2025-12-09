with ws as (
select *
from conciliacao.vw_ws t0
where t0.dt_fis >= current_date - interval '7 day'
--and t0.dt_fis < '20231107'
), tc as (
select *
from conciliacao.vw_tc57 t0
where t0.dt_fis >= current_date - interval '7 day'
--and t0.dt_fis < '20231107'
), gestao as (
select *
from conciliacao.vw_transacao_gestao t0
where t0.dt_processamento >= current_date - interval '7 day'
--and t0.dt_processamento < '20231107'
), interc as (
select *
from conciliacao.vw_transacao_interc t0
where t0.transaction_date >= current_date - interval '15 day'
--and t0.transaction_date < '20231107'
)
select  coalesce ( z0.ref_num_fis, z1.ref_num_fis) ref_num_fis
, coalesce ( z0.ref_num_bnd , z1.ref_num_bnd) ref_num_bnd
, case when z0.ref_num_bnd is null then false else true end ws
, case when z1.ref_num_bnd is null then false else true end tc57
, coalesce (z2.cap_status, false) gestao_cap
, coalesce (z2.fin_status, false) gestao_fin
, coalesce (z2.schedule_status, false) gestao_agenda
, coalesce (z3.tran_status, false) interc_cap
, coalesce (z3.schedule_status, false) interc_fin
, coalesce ( z0.bandeira , z1.bandeira) bandeira
, coalesce ( z0.tipo_transacao , z1.tipo_transacao) tipo_transacao
, coalesce ( z0.dt_fis , z1.dt_fis) dt_fis
, coalesce ( z0.dt_pos , z1.dt_pos) dt_pos
, coalesce ( cast(z0.term_loc as int8) , z1.term_loc) term_loc
, coalesce ( z0.term_id , z1.term_id) term_id
--coalesce ( z0.tran_amount , z1.tran_amount) tran_amount,
--coalesce ( z0.qtd_parc , z1.qtd_parc) num_parc,
, coalesce ( z1.bin, z0.bin) bin
, z3.emissor
--, z2.tp_pessoa
, z2.nm_fantasia
--, z2.razao_social
, z2.arranjo_pagamento arranjo_gestao
, z2.cd_mcc
--, z2.descricao_mcc
, cd_banco, nm_banco
--, case when z0.ref_num_bnd is not null and (z2.installment_number = 1 or z3.installment_number = 1 or (z2.ref_num_bnd is null and z3.ref_num_bnd is null)) then 1 else 0 end qtd_ws
--, case when z0.ref_num_bnd is not null and (z2.installment_number = 1 or z3.installment_number = 1 or (z2.ref_num_bnd is null and z3.ref_num_bnd is null)) then z0.tran_amount else 0 end vlr_ws
, case when z1.ref_num_bnd is not null and (z2.installment_number = 1 or z3.installment_number = 1 or (z2.ref_num_bnd is null and z3.ref_num_bnd is null)) then 1 else 0 end qtd_tc57
, case when z1.ref_num_bnd is not null and (z2.installment_number = 1 or z3.installment_number = 1 or (z2.ref_num_bnd is null and z3.ref_num_bnd is null)) then z1.tran_amount else 0 end vlr_tc57
--, z2.cap_qty qtd_gestao_cap
--, z2.fin_qty qtd_gestao_fin
, z2.schedule_qty qtd_gestao_agenda
--, z2.cap_amount vlr_gestao_cap
--, z2.fin_amount vlr_gestao_fin
, z2.installment_number parcela_gestao_num
, z2.installment_gross_amount vlr_gestao_agenda_bruto
, z2.installment_net_amount vlr_gestao_agenda_liquido
, z2.installment_tax_amount vlr_gestao_agenda_taxa
, z2.tax_mdr_applied * 100.00 tx_mdr_aplicada
, z2.tax_rav_ec * 100.00 tx_rav_cadastrada
, z2.dt_agenda dt_gestao_agenda
, z2.agenda_status
, z2.dt_efetivacao_pagamento
, z2.liquidacao_tipo tipo_liquidacao_gestao
, z2.liquidacao_status status_liquidacao_gestao
--, z3.tran_qty qtd_interc_cap
--, z3.schedule_qty qtd_interc_agenda
, z3.outgoing_qty qtd_interc_outgoing
, z3.tran_amount vlr_interc_cap
, z3.installment_number parcela_interc_num
, z3.installment_amount vlr_interc_agenda_bruto
, z3.exchange_rate_amount tx_interc_agenda
, z3.exchange_amount vlr_interc_agenda_taxa
, z3.installment_net_amount vlr_interc_agenda_liquido
, z3.tipo_apresentacao tipo_apresentacao_intercambio
, z3.scheduled_date dt_interc_liberacao
, z3.payment_scheduled_date dt_interc_previsao_recebimento
, z1.tc57_file_name
, z3.outgoing_file_name
from ws z0
full outer join tc z1 on z1.chave_join = z0.chave_join
full outer join gestao z2 on z2.chave_join = coalesce ( z0.chave_join , z1.chave_join)
full outer join interc z3 on z3.chave_join = coalesce ( z0.chave_join , z1.chave_join, z2.chave_join) and z2.installment_number = z3.installment_number
Where 1 = 1 and coalesce ( z0.dt_fis::date , z1.dt_fis::date) = to_date('16/09/2025', 'dd/MM/yyyy');

