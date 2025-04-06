create table viento_items (
    id integer primary key autoincrement,
    name text,
    price decimal(10, 2),
    ean text
);

CREATE INDEX idx_operation_date ON operations (operation_date);

