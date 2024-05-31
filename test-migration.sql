CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(50) PRIMARY KEY UNIQUE,
    track_number VARCHAR(50) UNIQUE,
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(50),
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(10),
    sm_id BIGINT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS delivery (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(50) REFERENCES orders(order_uid),
    name VARCHAR(100),
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100),
    address VARCHAR(200),
    region VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS payment (
    id SERIAL PRIMARY KEY,
    transaction VARCHAR(50) REFERENCES orders(order_uid),
    request_id VARCHAR(100),
    currency VARCHAR(10),
    provider VARCHAR(100),
    amount INT,
    payment_dt INT,
    bank VARCHAR(100),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    chrt_id INT,
    track_number VARCHAR(50) REFERENCES orders(track_number),
    price INT,
    rid VARCHAR(100),
    name VARCHAR(200),
    sale INT,
    size VARCHAR(20),
    total_price INT,
    nm_id INT,
    brand VARCHAR(100),
    status INT
);
