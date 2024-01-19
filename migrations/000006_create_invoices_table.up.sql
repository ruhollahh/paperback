CREATE TABLE IF NOT EXISTS invoices
(
    id         bigserial PRIMARY KEY,
    order_id   bigint   NOT NULL REFERENCES orders ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    status     text                        NOT NULL,
    version     integer                     NOT NULL DEFAULT 1
);