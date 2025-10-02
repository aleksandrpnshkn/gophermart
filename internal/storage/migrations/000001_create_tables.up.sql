CREATE TABLE users (
    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login VARCHAR(30) NOT NULL UNIQUE,
    password_hash BYTEA NOT NULL
);

CREATE TABLE orders (
    number VARCHAR(30) NOT NULL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    accrual DECIMAL(12, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    uploaded_at TIMESTAMP NOT NULL,

    CONSTRAINT fk_orders_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

CREATE TABLE balance_logs (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(30) NOT NULL,
    user_id BIGINT NOT NULL,
    amount DECIMAL(12, 2) NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_balance_order_number
    FOREIGN KEY (order_number)
    REFERENCES orders (number)
    ON UPDATE CASCADE
    ON DELETE RESTRICT,

    CONSTRAINT fk_balance_user_id
    FOREIGN KEY (user_id)
    REFERENCES users (id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);
