CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    delivery_schema TEXT,
    order_date TEXT,
    posting_number TEXT,
    warehouse_id INTEGER
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    sku INTEGER,
    viento_id INTEGER,
    consumption INTEGER
);


CREATE TABLE services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation_id INTEGER,
    name TEXT,
    price REAL
);

-- транзакции при заказах, в том числе возвраты из файла
CREATE TABLE operations (
    operation_id           INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    operation_type         TEXT                              NOT NULL,
    operation_date         DATETIME                          NOT NULL,
    operation_type_name    TEXT                              NOT NULL,
    delivery_charge        REAL,
    return_delivery_charge REAL,
    accruals_for_sale      REAL,
    sale_commission        REAL,
    amount                 REAL,
    transaction_type       TEXT,
    posting_id             INTEGER,
    item_sku               INTEGER,
    services               INTEGER[]
);
