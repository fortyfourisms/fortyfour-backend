ALTER TABLE users
    ADD COLUMN status ENUM('Aktif', 'Suspend', 'Nonaktif') NOT NULL DEFAULT 'Aktif' AFTER id_perusahaan,
    ADD COLUMN password_changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER status,
    ADD COLUMN login_attempts INT NOT NULL DEFAULT 0 AFTER password_changed_at;

CREATE INDEX idx_users_status ON users(status);
