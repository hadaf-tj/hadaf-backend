-- 1. КАТЕГОРИИ (Их можно оставить без приставки "Тест", так как это справочник)
INSERT INTO categories (name) VALUES 
('Продукты питания'), ('Одежда и обувь'), ('Медикаменты'), 
('Бытовая химия'), ('Канцелярия'), ('Строительные материалы'), ('Мебель')
ON CONFLICT (name) DO NOTHING;

-- 2. УЧРЕЖДЕНИЯ (Таджикистан) - Добавлена приставка [ТЕСТ]
INSERT INTO institutions (name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude) VALUES 
(
    '[ТЕСТ] Дом-интернат для престарелых "Батош"', 'Elderly', 'Турсунзаде', 'РРП', 
    'ул. Парк Победы, 45', '+992 900 11 22 33', 'test_batosh@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Государственное учреждение для пожилых людей и инвалидов.', '08:00 - 17:00', 38.5134, 68.2256
),
(
    '[ТЕСТ] Республиканская школа-интернат №1', 'Children', 'Душанбе', 'Душанбе', 
    'ул. Айни, 126', '+992 918 44 55 66', 'test_internat1@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Школа-интернат для детей-сирот и детей из малообеспеченных семей.', '24/7', 38.5598, 68.7870
),
(
    '[ТЕСТ] Областной детский дом г. Худжанда', 'Children', 'Худжанд', 'Согд', 
    'ул. Сирдарья, 12', '+992 927 00 11 22', 'test_khujand_kids@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Учреждение для детей дошкольного возраста.', '09:00 - 18:00', 40.2833, 69.6167
),
(
    '[ТЕСТ] Центр социальной поддержки "Нур"', 'Disabled', 'Хорог', 'ГБАО', 
    'ул. Шош Саидов, 5', '+992 935 88 99 00', 'test_khorog_nur@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Помощь людям с ограниченными возможностями в горном регионе.', '10:00 - 16:00', 37.4894, 71.5529
)
ON CONFLICT (name) DO NOTHING;

-- 3. НУЖДЫ (Генерация через безопасный поиск ID)

-- --- НУЖДЫ ДЛЯ БАТОШ ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Мука (1 сорт)', '[ТЕСТОВЫЙ ЗАПРОС] Для выпечки свежего хлеба', 'кг', 500, 50, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Дом-интернат для престарелых "Батош"' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Мука (1 сорт)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Тонометры', '[ТЕСТОВЫЙ ЗАПРОС] Автоматические тонометры на плечо', 'шт', 10, 2, 'medium'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Дом-интернат для престарелых "Батош"' AND c.name = 'Медикаменты'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Тонометры' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Стиральный порошок', '[ТЕСТОВЫЙ ЗАПРОС] Для автоматических машин', 'кг', 100, 0, 'low'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Дом-интернат для престарелых "Батош"' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Стиральный порошок' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ДУШАНБЕ (ШКОЛА №1) ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Тетради (12 л.)', '[ТЕСТОВЫЙ ЗАПРОС] В клетку и в линейку для начальных классов', 'шт', 2000, 450, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Республиканская школа-интернат №1' AND c.name = 'Канцелярия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Тетради (12 л.)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Матрасы', '[ТЕСТОВЫЙ ЗАПРОС] Ортопедические детские матрасы 160x70', 'шт', 40, 5, 'medium'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Республиканская школа-интернат №1' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Матрасы' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХУДЖАНДА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Подгузники (Size 4)', '[ТЕСТОВЫЙ ЗАПРОС] Для детей от 9 до 14 кг', 'упак', 60, 10, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Областной детский дом г. Худжанда' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Подгузники (Size 4)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Детское питание (смесь)', '[ТЕСТОВЫЙ ЗАПРОС] Гипоаллергенные смеси для новорожденных', 'банка', 100, 20, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Областной детский дом г. Худжанда' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Детское питание (смесь)' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХОРОГА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Обогреватели', '[ТЕСТОВЫЙ ЗАПРОС] Масляные радиаторы для зимнего периода', 'шт', 15, 3, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Центр социальной поддержки "Нур"' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Обогреватели' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Цемент М500', '[ТЕСТОВЫЙ ЗАПРОС] Для ремонта пандуса у входа', 'мешок', 20, 0, 'low'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Центр социальной поддержки "Нур"' AND c.name = 'Строительные материалы'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Цемент М500' AND institution_id = i.id);