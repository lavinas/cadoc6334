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
        END::text AS key1,
    t2.descricao AS transaction_brand,
    t1.descricao AS transaction_product,
    t0.dt_fis AS transaction_date,
    t0.dt_pos AS dt_pos_?,
    t0.term_loc AS establishment_terminal_code,
    t0.term_id AS term_id_?,
    t0.tran_amount AS transaction_amount,
    t0.num_parc AS qtd_parc AS  transaction_installments,
    "left"(t0.masked_pan::text, 11) AS bin
    -- forma de captura -- transaction_capture
    -- codigo_autorizacao - authorization_code
    -- nsu -- transaction_nsu
    -- codigo estabelecimento - establishment_code
    -- tipo de parcelamento - transaction_installments_type
   FROM fis.transaction_detail t0
     JOIN fis.terminal_transaction t1 ON t1.id = t0.fis_term_tran_detid::bpchar AND t1.terminal_transaction_type = 1
     LEFT JOIN fis.brand t2 ON t2.id = t0.brand_id
  WHERE t0.brand_resp_code::text = '00'::text AND NOT (EXISTS (
     SELECT 1
           FROM fis.transaction_detail z0
             JOIN fis.terminal_transaction z1 ON z1.id = z0.fis_term_tran_detid::bpchar 
              AND (z1.terminal_transaction_type = ANY (ARRAY[2, 3]))
          WHERE z0.brand_resp_code::text = '00'::text AND z0.orig_ref_num_fis::text = t0.ref_num_fis::text             
          )   
        )
    AND t0.dt_fis >= %dt_inicio% AND t0.dt_fis < %dt_fim%;