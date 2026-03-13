-- Таблица учреждений
CREATE TABLE IF NOT EXISTS institutions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(100),                      -- 'Children' (Детский дом), 'Elderly' (Дом престарелых), 'Disabled' (Инвалиды)
    city VARCHAR(100),
    region VARCHAR(100),
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(150),
    description TEXT,                       -- Описание учреждения
    activity_hours TEXT,                    -- Часы для посещения (для волонтеров)
    latitude DECIMAL(9,6),                  -- Координаты
    longitude DECIMAL(9,6),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    needs_count INT DEFAULT 0,
    events_count INT DEFAULT 0
);

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    institution_id INT REFERENCES institutions(id) ON DELETE SET NULL, -- NULL для Админов и Волонтёров, ID для Работников
    full_name VARCHAR(150),
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(150) UNIQUE,
    password TEXT,                          -- Хеш пароля
    role VARCHAR(50) NOT NULL DEFAULT 'volunteer', -- 'super_admin', 'employee', 'volunteer'
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Категории нужд (Продукты, Гигиена, Одежда...)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Таблица нужд (Потребности)
CREATE TABLE IF NOT EXISTS needs (
    id SERIAL PRIMARY KEY,
    institution_id INT NOT NULL REFERENCES institutions(id) ON DELETE CASCADE,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,             -- Название (напр. "Подгузники")
    description TEXT,                       -- Детали (напр. "Размер 4, марка Pampers")
    unit VARCHAR(50) NOT NULL,              -- Ед. измерения (шт, кг, л, уп)
    required_qty DECIMAL(10,2) NOT NULL,    -- Сколько нужно
    received_qty DECIMAL(10,2) DEFAULT 0,   -- Сколько уже собрали
    urgency VARCHAR(20) DEFAULT 'medium',   -- 'low', 'medium', 'high' (для сортировки)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS needs_history (
    id SERIAL PRIMARY KEY,
    need_id INT NOT NULL REFERENCES needs(id) ON DELETE CASCADE,
    comment TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);


-- Таблица OTP (Одноразовые пароли для входа/регистрации)
CREATE TABLE IF NOT EXISTS otp (
    id SERIAL PRIMARY KEY,
    attempt INTEGER DEFAULT 0,
    receiver VARCHAR(100) NOT NULL,         -- Телефон или Email
    method VARCHAR(50),                     -- 'sms', 'email'
    otp_code VARCHAR(20) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Таблица бронирований (откликов волонтеров на нужды)
CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    need_id INT NOT NULL REFERENCES needs(id) ON DELETE CASCADE,
    quantity DECIMAL(10,2) NOT NULL,         -- Количество, которое волонтер готов принести
    note TEXT,                               -- Сообщение от волонтера
    status VARCHAR(20) DEFAULT 'pending',    -- 'pending', 'approved', 'rejected'
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Индексы для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_bookings_need_id ON bookings(need_id);
CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status ON bookings(status);
CREATE UNIQUE INDEX IF NOT EXISTS bookings_user_need_active_uniq
  ON bookings(user_id, need_id)
  WHERE status NOT IN ('cancelled', 'rejected') AND is_deleted = false;

-- Таблица волонтёрских событий
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMPTZ NOT NULL,        -- Дата и время события
    institution_id INT NOT NULL REFERENCES institutions(id) ON DELETE CASCADE,
    creator_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);

-- Таблица участников событий (M2M)
CREATE TABLE IF NOT EXISTS event_participants (
    event_id INT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (event_id, user_id)         -- Уникальная пара, нельзя записаться дважды
);

-- Индексы для событий
CREATE INDEX IF NOT EXISTS idx_events_event_date ON events(event_date);
CREATE INDEX IF NOT EXISTS idx_events_institution_id ON events(institution_id);
CREATE INDEX IF NOT EXISTS idx_events_creator_id ON events(creator_id);
-- Таблица для хранения и ротации refresh-токенов
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id           SERIAL PRIMARY KEY,
    user_id      INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash   TEXT NOT NULL UNIQUE,
    expires_at   TIMESTAMPTZ NOT NULL,
    is_revoked   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- Уникальный индекс: один пользователь — одна активная заявка на need
CREATE UNIQUE INDEX IF NOT EXISTS bookings_user_need_active_uniq
  ON bookings(user_id, need_id)
  WHERE status NOT IN ('cancelled', 'rejected');
