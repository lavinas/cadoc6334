SELECT t0.ref_num_fis,
        CASE
            WHEN t0.brand_id::text = '04'::text THEN 'V'::text
            WHEN t0.brand_id::text = '05'::text THEN 'M'::text
            WHEN t0.brand_id::text = '42'::text THEN 'E'::text
            ELSE 'N'::text
        END || t0.ref_num_bnd::text AS ref_num_bnd,
        CASE
            WHEN t0.brand_id::text = '04'::text THEN 'V'::text
            WHEN t0.brand_id::text = '05'::text THEN 'M'::text
            WHEN t0.brand_id::text = '42'::text THEN 'E'::text
            ELSE 'N'::text
        END ||
        CASE
            WHEN t0.brand_id::text = '05'::text THEN ("left"(t0.ref_num_bnd::text, 9) || to_char(t0.dt_fis, 'MMDD'::text)::character varying::text)::character varying
            ELSE t0.ref_num_bnd
        END::text AS chave_join,
    t2.descricao AS bandeira,
    t1.descricao AS tipo_transacao,
    t0.dt_fis,
    t0.dt_pos,
    t0.term_loc,
    t0.term_id,
    t0.tran_amount,
    t0.num_parc AS qtd_parc,
    "left"(t0.masked_pan::text, 6) AS bin
   FROM fis.transaction_detail t0
     JOIN fis.terminal_transaction t1 ON t1.id = t0.fis_term_tran_detid::bpchar AND t1.terminal_transaction_type = 1
     LEFT JOIN fis.brand t2 ON t2.id::text = t0.brand_id::text
  WHERE t0.brand_resp_code::text = '00'::text AND NOT (EXISTS (
     SELECT z0.id,
            z0.id_transactionrequest,
            z0.term_id,
            z0.term_serial_number,
            z0.term_loc,
            z0.dt_pos,
            z0.dt_fis,
            z0.trace,
            z0.fis_type_tranid,
            z0.masked_pan,
            z0.term_entry_mode,
            z0.fis_cardholder_idmeth,
            z0.chip_app_cryptogram,
            z0.ref_num_fis,
            z0.destination_code,
            z0.brand_id,
            z0.network_code,
            z0.ref_num_bnd,
            z0.auth_code,
            z0.term_resp_code,
            z0.brand_resp_code,
            z0.inf_adic_respcode,
            z0.tran_amount,
            z0.tipo_parc,
            z0.num_parc,
            z0.orig_ref_num_fis,
            z0.fis_term_tran_detid,
            z1.id,
            z1.descricao,
            z1.tp_operacao,
            z1.cd_operacao,
            z1.integra_gestao,
            z1.tp_parcela,
            z1.cd_tp_cartao,
            z1.terminal_transaction_type
           FROM fis.transaction_detail z0
             JOIN fis.terminal_transaction z1 ON z1.id = z0.fis_term_tran_detid::bpchar 
              AND (z1.terminal_transaction_type = ANY (ARRAY[2, 3]))
          WHERE z0.brand_resp_code::text = '00'::text AND z0.orig_ref_num_fis::text = t0.ref_num_fis::text
          )
        );
