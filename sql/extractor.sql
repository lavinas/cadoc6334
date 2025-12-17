-- Extractor Process Tables
CREATE TABLE extractor_process (
    -- id
    id BIGSERIAL PRIMARY KEY,
    -- parameters
    name VARCHAR(50) NOT NULL
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Daily Extractor Process Configuration
CREATE TABLE extractor_process_daily (
    -- ID
    id BIGSERIAL PRIMARY KEY, 
    extractor_process_id BIGINT,
    -- parameters
    window_lag TIME NOT NULL DEFAULT '01:00:00',
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(extractor_process_id),
    FOREIGN KEY (extractor_process_id) REFERENCES extractor_process(id)
);

-- Reprocessing Extractor Process Configuration
CREATE TABLE extractor_process_reprocessing (
    -- id
    id BIGSERIAL PRIMARY KEY,
    extractor_process_id BIGINT,
    -- parameters
    window_lag TIME NOT NULL DEFAULT '01:00:00',
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(extractor_process_id),
    FOREIGN KEY extractor_process_id REFERENCES extractor_process(id)
);

-- Totals Extractor Process Configuration
CREATE TABLE extractor_process_totals (
    -- id
    id BIGSERIAL PRIMARY KEY,
    extractor_process_id BIGINT,
    -- parameters
    rollback_days INT NOT NULL,
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    UNIQUE(extractor_process_id),
    FOREIGN KEY extractor_process_id REFERENCES extractor_process(id)
);

-- Daily Extractor Process Control
CREATE TABLE extractor_daily_control (
    -- id
    id BIGSERIAL PRIMARY KEY,
    extractor_process_id BIGINT NOT NULL,
    -- control
    last_period_start TIMESTAMP,
    last_period_end TIMESTAMP,
    last_total INT,
    last_quantity INT,
    last_processing_start TIMESTAMP,
    last_processing_end TIMESTAMP,
    last_status VARCHAR(20),
    last_trace_id VARCHAR(50),
    -- status
    status_id INT NOT NULL,
    status_name VARCHAR(20) NOT NULL,
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (extractor_process_id),
    FOREIGN KEY (extractor_process_id) REFERENCES extractor_process(id)
);
 
-- Reprocessing Extractor Process Control
CREATE TABLE extractor_reprocessing_control (
    -- id
    id BIGSERIAL PRIMARY KEY,
    extractor_process_id BIGINT NOT NULL,
    -- entry
    required_period_start date NOT NULL,
    required_period_end date NOT NULL,
    -- control
    last_period_start TIMESTAMP,
    last_period_end TIMESTAMP,
    last_processing_start TIMESTAMP,
    last_processing_end TIMESTAMP,
    last_total INT,
    last_quantity INT,
    last_status VARCHAR(20),
    last_trace_id VARCHAR(50),
    -- status
    status_id INT NOT NULL,
    status_name VARCHAR(20) NOT NULL,
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (extractor_process_id) REFERENCES extractor_process(id)
);
 
-- Extractor Execution Records
CREATE TABLE extractor_execution (
    -- id
    id BIGSERIAL PRIMARY KEY,
    process_daily_id BIGINT,
    process_reprocessing_id BIGINT,
    execution_type VARCHAR(20) NOT NULL,  -- 'DAILY' ou 'REPROCESS'
    -- parameters
    trace_id VARCHAR(50) NOT NULL,
    period_start TIMESTAMP NOT NULL,
    period_end TIMESTAMP NOT NULL,
    execution_total INT NOT NULL DEFAULT 0,
    execution_quantity INT NOT NULL DEFAULT 0,
    execution_start TIMESTAMP,
    execution_end TIMESTAMP,
    -- status
    status_id INT NOT NULL,
    status_name VARCHAR(20) NOT NULL,
    error_message VARCHAR(300),
    -- structure
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (process_daily_id) REFERENCES extractor_process_daily(id),
    FOREIGN KEY (process_reprocessing_id) REFERENCES extractor_reprocessing_control(id)
);
 
 
 
 