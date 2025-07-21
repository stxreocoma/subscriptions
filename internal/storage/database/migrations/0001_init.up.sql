CREATE TABLE IF NOT EXISTS subscriptions (
    service_name VARCHAR(100) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (service_name, user_id),
    start_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_date TIMESTAMP
);