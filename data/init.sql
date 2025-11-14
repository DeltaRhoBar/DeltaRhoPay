PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS residents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    r_floor INTEGER NOT NULL,
    r_nr INTEGER NOT NULL,
    name TEXT NOT NULL,
    removed_on TEXT
);

CREATE TABLE IF NOT EXISTS debts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount INTEGER NOT NULL DEFAULT 0,
    date TEXT,
    resident_id INTEGER,
    FOREIGN KEY(resident_id) REFERENCES residents(id)
);

CREATE TRIGGER IF NOT EXISTS create_invoice_for_new_resident
AFTER INSERT ON residents
FOR EACH ROW
BEGIN
    INSERT INTO debts (date, resident_id)
    VALUES (NULL, NEW.id);
END;

CREATE TRIGGER IF NOT EXISTS create_new_debt_after_date_set
AFTER UPDATE ON debts
FOR EACH ROW
WHEN (OLD.date IS NULL OR OLD.date = '') AND (NEW.date IS NOT NULL AND NEW.date <> '')
BEGIN
    INSERT INTO debts (date, resident_id)
    VALUES (NULL, NEW.resident_id);
END;
