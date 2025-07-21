CREATE INDEX IF NOT EXISTS idx_service_user ON subscriptions (service_name, user_id);
CREATE INDEX IF NOT EXISTS idx_created_at ON subscriptions (start_date);