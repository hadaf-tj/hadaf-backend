CREATE TABLE IF NOT EXISTS institutions (
                              id SERIAL PRIMARY KEY,
                              name VARCHAR(255) NOT NULL,
                              type VARCHAR(100),                     -- тип учреждения: детский дом, дом престарелых и т.д.
                              city VARCHAR(100),
                              region VARCHAR(100),
                              address TEXT,
                              phone VARCHAR(50),
                              latitude DECIMAL(9,6),                 -- координаты для карты
                              longitude DECIMAL(9,6),
                              created_at TIMESTAMPZ DEFAULT NOW(),
                              updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
                       id SERIAL PRIMARY KEY,
                       full_name VARCHAR(150),
                       phone VARCHAR(20) UNIQUE,
                       email VARCHAR(150) UNIQUE,
                       password TEXT,
                       role TEXT,    -- enum
                       is_active BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMPZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS categories (
                            id SERIAL PRIMARY KEY,
                            name VARCHAR(100) UNIQUE NOT NULL,
                            created_at TIMESTAMPZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS needs (
                       id SERIAL PRIMARY KEY,
                       institution_id INT NOT NULL REFERENCES institutions(id) ON DELETE CASCADE,
                       category_id INT REFERENCES categories(id) ON DELETE SET NULL,
                       name VARCHAR(255) NOT NULL,                -- что нужно
                       unit VARCHAR(50) NOT NULL,                 -- единица измерения: шт, кг, л
                       required_qty DECIMAL(10,2) NOT NULL,       -- требуется
                       received_qty DECIMAL(10,2) DEFAULT 0,      -- уже получено
                       created_at TIMESTAMPZ DEFAULT NOW()
                       updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS otp (
    id SERIAL PRIMARY KEY,
    attempt    INTEGER,
    receiver VARCHAR(100),
    method      VARCHAR(50),
    otp_code     VARCHAR(20),
    is_verified  BOOLEAN,
    sent_at      TIMESTAMP,
    expires_at   TIMESTAMP,
    updated_at   TIMESTAMP
);

