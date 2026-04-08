CREATE TABLE IF NOT EXISTS ikas_audit_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ikas_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    action VARCHAR(50) NOT NULL,
    changes JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_ikas_id (ikas_id),
    INDEX idx_user_id (user_id)
);
