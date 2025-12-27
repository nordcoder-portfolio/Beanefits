-- +goose Up
CREATE TABLE user_roles
(
    user_id   BIGINT NOT NULL,
    role_code TEXT   NOT NULL,
    CONSTRAINT pk_user_roles PRIMARY KEY (user_id, role_code),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_code) REFERENCES roles (code) ON DELETE RESTRICT
);

-- +goose Down
DROP TABLE user_roles;
