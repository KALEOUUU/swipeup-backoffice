-- +mig Up
-- Create activity_logs table for tracking user activities
CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Index for performance
    INDEX idx_activity_logs_user_id (id_user),
    INDEX idx_activity_logs_action (action),
    INDEX idx_activity_logs_created_at (created_at)
);

-- +mig Down
DROP TABLE IF EXISTS activity_logs;