CREATE TABLE IF NOT EXISTS conversation_sessions (
    id VARCHAR(36) PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_message_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    feedback_prompt_sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_session_status CHECK (status IN ('active', 'waiting_feedback', 'closed'))
);

CREATE INDEX IF NOT EXISTS idx_sessions_phone_number ON conversation_sessions(phone_number);
CREATE INDEX IF NOT EXISTS idx_sessions_status ON conversation_sessions(status);
CREATE INDEX IF NOT EXISTS idx_sessions_last_message_at ON conversation_sessions(last_message_at);

CREATE TABLE IF NOT EXISTS feedbacks (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    rating INT NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_feedback_session FOREIGN KEY (session_id) REFERENCES conversation_sessions(id) ON DELETE CASCADE,
    CONSTRAINT chk_feedback_rating CHECK (rating >= 1 AND rating <= 5)
);

CREATE INDEX IF NOT EXISTS idx_feedbacks_session_id ON feedbacks(session_id);
CREATE INDEX IF NOT EXISTS idx_feedbacks_phone_number ON feedbacks(phone_number);
CREATE INDEX IF NOT EXISTS idx_feedbacks_rating ON feedbacks(rating);
CREATE INDEX IF NOT EXISTS idx_feedbacks_created_at ON feedbacks(created_at DESC);
