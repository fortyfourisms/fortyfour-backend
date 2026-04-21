CREATE TABLE se_edit_request (
    id CHAR(36) PRIMARY KEY,
    id_se CHAR(36) NOT NULL,
    id_user CHAR(36) NOT NULL,
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    catatan_user TEXT,
    catatan TEXT,
    data_perubahan JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_ser_se FOREIGN KEY (id_se) REFERENCES se(id) ON DELETE CASCADE,
    CONSTRAINT fk_ser_user FOREIGN KEY (id_user) REFERENCES users(id) ON DELETE CASCADE
);
