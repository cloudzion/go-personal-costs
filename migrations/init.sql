CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO categories (name) VALUES 
    ('Еда'), 
    ('Транспорт'), 
    ('Жильё'), 
    ('Развлечения'), 
    ('Прочее');

CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    description TEXT NOT NULL,
    date TEXT NOT NULL,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);
