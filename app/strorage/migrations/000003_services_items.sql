CREATE TABLE postings (
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


CREATE TABLE iut (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    operation_id INTEGER,
    name TEXT,
    price REAL
);

CREATE TABLE operation (
    id INTEGER PRIMARY KEY,
    operation_id INTEGER,
    operation_type TEXT,
    operation_date TEXT,
    operation_type_name TEXT,
    delivery_charge REAL,
    return_delivery_charge REAL,
    accruals_for_sale REAL,
    sale_commission REAL,
    amount REAL,
    transaction_type TEXT,
    posting_id INTEGER,
    items INTEGER,
    services INTEGER
);

