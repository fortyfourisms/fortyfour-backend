CREATE TABLE se_edit_request (
    id CHAR(36) PRIMARY KEY,

    responden_id INT NOT NULL,
    risiko_id INT NOT NULL,
    user_id CHAR(36) NOT NULL,

    status ENUM('pending', 'approved', 'rejected') 
        NOT NULL DEFAULT 'pending',

    catatan_user TEXT,
    catatan_admin TEXT,

    data_perubahan JSON NOT NULL,

    approved_by CHAR(36) DEFAULT NULL,
    approved_at TIMESTAMP NULL DEFAULT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
        ON UPDATE CURRENT_TIMESTAMP,

    -- RELATION
    CONSTRAINT fk_ser_responden 
        FOREIGN KEY (responden_id) 
        REFERENCES responden(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_ser_risiko 
        FOREIGN KEY (risiko_id) 
        REFERENCES risiko(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_ser_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,

    CONSTRAINT fk_ser_approved 
        FOREIGN KEY (approved_by) 
        REFERENCES users(id) 
        ON DELETE SET NULL
);