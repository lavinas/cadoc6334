
-- Falta NSU, codigo de autorizacao, codigo do estabelecimento, codigo do terminal e chave unica (basse bruna)
SELECT
    fc.brand_transaction_identifier as cd_transacao_fin, -- nsu
    fc.pos_entry_mode as forma_captura,
    fc.transaction_date as dt_processamento,
    fc.transaction_amount as ValorTransacoes,
        case
        when fci.exchange_rate_type_name = 'PERCENT' then fci.exchange_rate       -- já é %
        when fci.exchange_rate_type_name = 'FIXED'   then (fci.exchange_rate / nullif(fc.transaction_amount,0)) * 100
        else null
    end as percentual_desconto,
    case
        when fci.exchange_rate_type_name = 'PERCENT' then round((fc.transaction_amount * fci.exchange_rate / 100),2) -- valor calculado
        when fci.exchange_rate_type_name = 'FIXED'   then fci.exchange_rate
        else null
    end as taxa_intercambio_valor,
        fc.brand_name as bandeira,
    fci.plan as parcela,
    fc.product_name as tipo_cartao,
    fc.merchant_category_code as segmento,
    fc.mask_account_number AS bin,
    NULL as modalidade_cartao,
    NULL as produto_cartao,
    NULL as cadoc_item_id
FROM interc.fin_contract fc
JOIN (
    SELECT
        fci.contract_id,
        fci.plan,
        fci.exchange_rate_amount as exchange_rate,
        fci.exchange_rate_type_name
    FROM interc.fin_contract_installment fci
    JOIN interc.fin_contract fc2
        ON fci.contract_id = fc2.id
    WHERE fc2.transaction_date >= ? AND fc2.transaction_date <= ?
    and fci.installment_number = 1
    GROUP BY fci.contract_id, fci.plan, fci.exchange_rate_type_name, fci.exchange_rate_amount
) fci
    ON fci.contract_id = fc.id



-- fin_contract.authorization_code
-- fin_contract.terminal_id
-- fin_contract.card_acceptor_id
-- fin_contract.created_date // fin_contract.updated_date -- sempre preenchido

-- fin_contract.contract_type_name -- 

    