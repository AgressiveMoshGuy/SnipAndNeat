
-- начисления от озона из отчета о реализации
CREATE TABLE ozon_pay (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    date_accrual DATETIME NOT NULL,
    ozon_id TEXT NOT NULL,
    date_of_purchase DATETIME NOT NULL,
    warehouse TEXT,
    sku INTEGER NOT NULL,
    vendor_code TEXT,
    description TEXT,
    amount INTEGER,
    raw_sum REAL,
    commission_bid INTEGER,
    commission_sum REAL,
    last_mile REAL,
    return_mile REAL,
    return_treatment REAL,
    logistic REAL,
    back_logistic REAL,
    total REAL
);


INSERT INTO ozon_pay (title, date_accrual, ozon_id, date_of_purchase, warehouse, sku, vendor_code, description, amount, raw_sum, commission_bid, commission_sum, last_mile, return_mile, return_treatment, logistic, back_logistic, total)
VALUES
    ('test1', '2022-12-01 00:00:00', '1234567890', '2022-11-01 00:00:00', 'test', 1, 'test', 'test', 1, 100.0, 1, 10.0, 10.0, 10.0, 10.0, 10.0, 10.0, 100.0),
    ('test2', '2022-12-02 00:00:00', '1234567891', '2022-11-02 00:00:00', 'test', 2, 'test', 'test', 2, 200.0, 2, 20.0, 20.0, 20.0, 20.0, 20.0, 20.0, 200.0),
    ('test3', '2022-12-03 00:00:00', '1234567892', '2022-11-03 00:00:00', 'test', 3, 'test', 'test', 3, 300.0, 3, 30.0, 30.0, 30.0, 30.0, 30.0, 30.0, 300.0),
    ('test4', '2022-12-04 00:00:00', '1234567893', '2022-11-04 00:00:00', 'test', 4, 'test', 'test', 4, 400.0, 4, 40.0, 40.0, 40.0, 40.0, 40.0, 40.0, 400.0),
    ('test5', '2022-12-05 00:00:00', '1234567894', '2022-11-05 00:00:00', 'test', 5, 'test', 'test', 5, 500.0, 5, 50.0, 50.0, 50.0, 50.0, 50.0, 50.0, 500.0),
    ('test6', '2022-12-06 00:00:00', '1234567895', '2022-11-06 00:00:00', 'test', 6, 'test', 'test', 6, 600.0, 6, 60.0, 60.0, 60.0, 60.0, 60.0, 60.0, 600.0),
    ('test7', '2022-12-07 00:00:00', '1234567896', '2022-11-07 00:00:00', 'test', 7, 'test', 'test', 7, 700.0, 7, 70.0, 70.0, 70.0, 70.0, 70.0, 70.0, 700.0),
    ('test8', '2022-12-08 00:00:00', '1234567897', '2022-11-08 00:00:00', 'test', 8, 'test', 'test', 8, 800.0, 8, 80.0, 80.0, 80.0, 80.0, 80.0, 80.0, 800.0),
    ('test9', '2022-12-09 00:00:00', '1234567898', '2022-11-09 00:00:00', 'test', 9, 'test', 'test', 9, 900.0, 9, 90.0, 90.0, 90.0, 90.0, 90.0, 90.0, 900.0),
    ('test10', '2022-12-10 00:00:00', '1234567899', '2022-11-10 00:00:00', 'test', 10, 'test', 'test', 10, 1000.0, 10, 100.0, 100.0, 100.0, 100.0, 100.0, 100.0, 1000.0);
