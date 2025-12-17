-- process table to store monitoring process definitions
create table process (
    -- id
    id bigserial primary key,
    -- parameters
    name varchar(50) not null,
    description varchar(255) null,
    -- classification
    flow_id int not null,
    flow_name varchar(50) not null,
    source_id int not null,
    source_name varchar(50) not null,
    -- structure
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

-- process time limit table to store time limits for each process
create table process_time_limit (
    -- id
    id bigserial primary key,
    process_id bigint not null,
    -- parameters
    periodicity varchar(20) not null, -- daily, hourly, weekly
    time_limit time not null,
    -- structure
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    unique(process_id),
    foreign key (process_id) references process(id)
);

-- process message table to store messages related to processes
create table process_message (
    -- id
    id bigserial primary key,
    process_id bigint not null,
    -- message
    message_type_id int not null, -- 1 - timeout, 2 - error, 3 - indicator
    message_subject varchar(200) not null,
    message_body text not null,
    -- structure
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    foreign key (process_id) references process(id)
);

-- process error table to store errors related to processes
create table process_error (
    -- id
    id bigserial primary key,
    process_id bigint not null,
    -- parameters
    error_key varchar(100) not null,
    description varchar(255) not null,
    generate_call BOOLEAN not null,
    message_body text not null,
    -- structure
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    foreign key (process_id) references process(id)
);

-- process_indicator table to store indicators for monitoring processes
create table process_indicator (
    -- id
    id bigserial primary key,
    process_id bigint not null,
    name varchar(100) not null,
    -- parameters
    process_reference_id bigint not null,
    under_var numeric(5,4) not null,
    over_var numeric(5,4) not null,
    message_body text not null,
    -- structure
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp,
    foreign key (process_id) references process(id),
    foreign key (process_reference_id) references process(id)
);

-- process_execution - substitui antiga process_daily_processing e process_indicadot_processing
create table process_execution (
    -- id
    id bigserial primary key,
    -- parameters
    process_id bigint not null,
    reference_date date not null,
    -- processing status
    status_id int not null,
    status_name varchar(20) not null,
    remarks varchar(300) null,
);



