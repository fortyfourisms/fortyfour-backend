DROP INDEX idx_users_perusahaan ON users;

ALTER TABLE users
    DROP FOREIGN KEY fk_users_perusahaan,
    DROP COLUMN id_perusahaan;
