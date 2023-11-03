GRANT ALL PRIVILEGES ON DATABASE L0_db TO some_user;

\c L0_db;

CREATE TABLE IF NOT EXISTS orders
(
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(10),
    sm_id              INTEGER,
    date_created       TIMESTAMP,
    oof_shard          VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS delivery
(
    order_uid varchar(255) references orders,
    name      varchar(255),
    phone     varchar(255),
    zip       varchar(255),
    city      varchar(255),
    address   text,
    region    varchar(255),
    email     varchar(255)
);

CREATE TABLE IF NOT EXISTS payment
(
    order_uid     varchar(255) references orders,
    transaction   varchar(255),
    request_id    varchar(255),
    currency      varchar(3),
    provider      varchar(255),
    amount        numeric(10, 2),
    payment_dt    integer,
    bank          varchar(255),
    delivery_cost numeric(10, 2),
    goods_total   numeric(10, 2),
    custom_fee    numeric(10, 2)

);

CREATE TABLE IF NOT EXISTS items
(
    order_uid    varchar(255) references orders,
    chrt_id      bigint,
    track_number varchar(255),
    price        numeric(10, 2),
    rid          varchar(255),
    name         text,
    sale         integer,
    size         varchar(255),
    total_price  numeric(10, 2),
    nm_id        bigint,
    brand        varchar(255),
    status       integer
);

GRANT ALL PRIVILEGES ON TABLE orders TO some_user;
GRANT ALL PRIVILEGES ON TABLE delivery TO some_user;
GRANT ALL PRIVILEGES ON TABLE payment TO some_user;
GRANT ALL PRIVILEGES ON TABLE items TO some_user;
