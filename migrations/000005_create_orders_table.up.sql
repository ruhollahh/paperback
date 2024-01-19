CREATE TABLE IF NOT EXISTS orders
(
    id          bigserial PRIMARY KEY,
    created_at  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    user_id     bigint                      NOT NULL REFERENCES users ON DELETE CASCADE,
    status      text                        NOT NULL,
    total_price integer                     NOT NULL,
    version     integer                     NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id_status ON orders (user_id, status);

CREATE TABLE IF NOT EXISTS order_items
(
    order_id   bigint REFERENCES orders ON DELETE CASCADE,
    product_id bigint REFERENCES products ON DELETE CASCADE,
    quantity   integer NOT NULL,
    price      integer NOT NULL,
    version    integer NOT NULL DEFAULT 1,
    PRIMARY KEY (order_id, product_id)
);