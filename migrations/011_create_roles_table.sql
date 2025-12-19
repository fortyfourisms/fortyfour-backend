CREATE TABLE roles (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_role_name (name)
);

INSERT INTO roles (id, name, description) VALUES 
(UUID(), 'admin', 'Administrator dengan akses penuh ke seluruh sistem'),
(UUID(), 'user', 'User biasa dengan akses terbatas');