
-- transaction_source table to store sources of transactions
create table transaction_source (
    -- control
    id bigint primary key,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    -- values
    name varchar(20) not null, -- 1 - webservice, 2 - intercambio, 3 - gestao
    priority int not null -- 10 - gestao, 20 - intercambio, 30 - webservice
);


-- transaction table is the main table to store transaction data
create table transaction (
    -- control fields
    id bigserial primary key,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    -- keys fields 
    key1 varchar(50) not null,
    key2 varchar(50) null,
    key3 varchar(50) null,
    -- establishment fields
    establishment_code numeric(50) null,
    establishment_mcc numeric(4) null,
    establishment_terminal_code numeric(50) null,
    -- authorization fields
    bin numeric(11) not null,
    authorization_code varchar(20) null,
    transaction_nsu varchar(20) not null, 
    -- transaction fields
    transaction_date date not null,
    transaction_amount numeric(15,2) not null,
    transaction_installments numeric(2) not null, -- 1 to 12
    transaction_installments_type varchar(10) not null, -- 'loja', 'emissor'
    transaction_brand varchar(2) not null, -- 'V' Visa, 'M' Mastercard, 'E' Elo
    transaction_product varchar(2) not null, -- 'DB' debit, 'CR' credit, 'PX' Pix, 'PP' Pre-pago
    transaction_capture varchar(3) not null, -- POS, TEF
    -- financial values
    revenue_mdr_value numeric(15,2) null,
    cost_interchange_value numeric(15,2) null
    -- references/control
    high_source_priority int not null, -- 30 - webservice, 20 - intercambio, 10 - gestao (o de maior prioridade)
    -- closing fields
    period_date date null -- inicialmente nulo
);

-- transaction_events table to store events related to transactions
create table transaction_event (
    -- control
    id bigserial primary key,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    -- ref
    transaction_id bigserial not null,
    -- values
    source_id int not null, -- 1 - webservice, 2 - intercambio, 3 - gestao
    source_name varchar(20) not null, -- 1 - webservice, 2 - intercambio, 3 - gestao
    source_datetime timestamp not null, -- data do reconhecimento da transacao no source
    source_key varchar(50) not null, -- id da tabela bruta
    -- structure
    foreign key (transaction_id) references transaction(id)
);

-- transaction_log table to store logs related to transactions
create table transaction_event_gap (
    -- control
    id bigserial primary key,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    -- reference/control
    event_id bigserial not null,
    high_source_id int not null,
    minor_source_id int not null,
    updated BOOLEAN not null,
    -- values
    field_name varchar(50) not null,
    high_value varchar(100) not null,
    minor_value varchar(100) not null,
    foreign key (event_id) references transaction_events(id)
);