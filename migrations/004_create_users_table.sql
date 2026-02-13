CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL,
    foto_profile VARCHAR(255) DEFAULT NULL,
    banner VARCHAR(255) DEFAULT NULL,
    role_id CHAR(36) DEFAULT NULL,
    id_jabatan CHAR(36) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    -- Indexes
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role_id (role_id),
    INDEX idx_jabatan (id_jabatan),
    
    -- Foreign Keys
    CONSTRAINT fk_users_role 
        FOREIGN KEY (role_id) 
        REFERENCES roles(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE,
        
    CONSTRAINT fk_users_jabatan 
        FOREIGN KEY (id_jabatan) 
        REFERENCES jabatan(id) 
        ON DELETE SET NULL 
        ON UPDATE CASCADE
);