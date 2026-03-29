PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS residents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    r_floor INTEGER NOT NULL,
    r_nr INTEGER NOT NULL,
    name TEXT NOT NULL,
    telephone TEXT NOT NULL,
    removed_on TEXT
);

CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    beverage_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    date TEXT DEFAULT CURRENT_TIMESTAMP,
    resident_id INTEGER NOT NULL,
    paid_on TEXT,
    FOREIGN KEY(beverage_id) REFERENCES beverages(id),
    FOREIGN KEY(resident_id) REFERENCES residents(id)
);

CREATE TABLE IF NOT EXISTS beverages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    removed_on TEXT
);

CREATE TABLE IF NOT EXISTS checkouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS key_value (
    key TEXT PRIMARY KEY,
    value TEXT,
);

CREATE UNIQUE INDEX idx_unique_floor_nr_removed
ON residents (r_floor, r_nr)
WHERE removed_on IS NULL;

CREATE UNIQUE INDEX idx_unique_name_date
ON beverages (name)
WHERE removed_on IS NULL;

