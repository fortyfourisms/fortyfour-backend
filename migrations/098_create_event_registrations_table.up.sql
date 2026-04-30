CREATE TABLE event_registrations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    event_id BIGINT NOT NULL,
    nama VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    perusahaan VARCHAR(255) NOT NULL,
    jabatan VARCHAR(255) NOT NULL,
    no_hp VARCHAR(50) NOT NULL,
    sektor VARCHAR(255) NOT NULL,
    qr_payload JSON NOT NULL,
    qr_token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_event_registrations_event
        FOREIGN KEY (event_id) REFERENCES events(id)
        ON DELETE CASCADE,
    UNIQUE KEY uq_event_registration_email (event_id, email),
    UNIQUE KEY uq_event_registration_qr_token (qr_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
