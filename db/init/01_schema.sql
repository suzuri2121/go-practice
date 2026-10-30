CREATE TABLE customers (
    id         integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       text NOT NULL,
    prefecture text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id       integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     text NOT NULL,
    category text NOT NULL,
    price    integer NOT NULL CHECK (price >= 0)
);

CREATE TABLE orders (
    id          integer GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_id integer NOT NULL REFERENCES customers(id),
    ordered_at  timestamptz NOT NULL
);

CREATE TABLE order_items (
    order_id   integer NOT NULL REFERENCES orders(id),
    product_id integer NOT NULL REFERENCES products(id),
    quantity   integer NOT NULL CHECK (quantity > 0),
    unit_price integer NOT NULL,
    PRIMARY KEY (order_id, product_id)
);

INSERT INTO customers (name, prefecture) VALUES
    ('佐藤', '東京'),
    ('鈴木', '神奈川'),
    ('高橋', '大阪'),
    ('田中', '東京'),
    ('伊藤', '福岡');

INSERT INTO products (name, category, price) VALUES
    ('キーボード', '周辺機器', 12000),
    ('マウス', '周辺機器', 4000),
    ('モニター', '周辺機器', 35000),
    ('Go入門', '書籍', 3200),
    ('SQL実践', '書籍', 2800),
    ('コーヒー豆', '食品', 1800);

INSERT INTO orders (customer_id, ordered_at) VALUES
    (1, '2026-08-03 10:00+09'),
    (2, '2026-08-10 14:30+09'),
    (1, '2026-09-01 09:15+09'),
    (3, '2026-09-15 20:00+09'),
    (4, '2026-09-20 12:00+09'),
    (2, '2026-10-02 18:45+09'),
    (1, '2026-10-10 08:30+09'),
    (3, '2026-10-21 22:10+09');

INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
    (1, 1, 1, 12000), (1, 2, 1, 4000),
    (2, 4, 2, 3200),
    (3, 3, 1, 35000),
    (4, 5, 1, 2800), (4, 4, 1, 3200),
    (5, 2, 2, 4000),
    (6, 1, 1, 11000),
    (7, 5, 3, 2800),
    (8, 3, 1, 33000), (8, 2, 1, 4000);