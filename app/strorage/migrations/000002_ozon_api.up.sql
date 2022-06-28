-- транзакции при заказах, в том числе возвраты
create table ozon_orders
(
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
    item_id                INTEGER,
    services               TEXT
);

