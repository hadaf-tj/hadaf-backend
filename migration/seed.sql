-- 1. КАТЕГОРИИ
INSERT INTO categories (name) VALUES 
('Продукты питания'), ('Одежда и обувь'), ('Медикаменты'), 
('Бытовая химия'), ('Канцелярия'), ('Строительные материалы'), ('Мебель')
ON CONFLICT (name) DO NOTHING;

-- 2. УЧРЕЖДЕНИЯ (Таджикистан)
INSERT INTO institutions (name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude) VALUES 
(
    'Дом-интернат для престарелых "Батош"', 'Elderly', 'Турсунзаде', 'РРП', 
    'ул. Парк Победы, 45', '+992 900 11 22 33', 'batosh@hadaf.tj', 
    'Государственное учреждение для пожилых людей и инвалидов.', '08:00 - 17:00', 38.5134, 68.2256
),
(
    'Республиканская школа-интернат №1', 'Children', 'Душанбе', 'Душанбе', 
    'ул. Айни, 126', '+992 918 44 55 66', 'internat1@hadaf.tj', 
    'Школа-интернат для детей-сирот и детей из малообеспеченных семей.', '24/7', 38.5598, 68.7870
),
(
    'Областной детский дом г. Худжанда', 'Children', 'Худжанд', 'Согд', 
    'ул. Сирдарья, 12', '+992 927 00 11 22', 'khujand_kids@hadaf.tj', 
    'Учреждение для детей дошкольного возраста.', '09:00 - 18:00', 40.2833, 69.6167
),
(
    'Центр социальной поддержки "Нур"', 'Disabled', 'Хорог', 'ГБАО', 
    'ул. Шош Саидов, 5', '+992 935 88 99 00', 'khorog_nur@hadaf.tj', 
    'Помощь людям с ограниченными возможностями в горном регионе.', '10:00 - 16:00', 37.4894, 71.5529
)
ON CONFLICT (name) DO NOTHING;

-- 3. НУЖДЫ (Генерация через безопасный поиск ID)

-- Помощник: Функция для вставки нужд, чтобы не дублировать код в SQL
-- Мы просто сделаем несколько прямых вставок через SELECT для надежности.

-- --- НУЖДЫ ДЛЯ БАТОШ ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Мука (1 сорт)', 'Для выпечки свежего хлеба', 'кг', 500, 50, 'high'
FROM institutions i, categories c WHERE i.name = 'Дом-интернат для престарелых "Батош"' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Мука (1 сорт)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Тонометры', 'Автоматические тонометры на плечо', 'шт', 10, 2, 'medium'
FROM institutions i, categories c WHERE i.name = 'Дом-интернат для престарелых "Батош"' AND c.name = 'Медикаменты'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Тонометры' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Стиральный порошок', 'Для автоматических машин', 'кг', 100, 0, 'low'
FROM institutions i, categories c WHERE i.name = 'Дом-интернат для престарелых "Батош"' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Стиральный порошок' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ДУШАНБЕ (ШКОЛА №1) ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Тетради (12 л.)', 'В клетку и в линейку для начальных классов', 'шт', 2000, 450, 'high'
FROM institutions i, categories c WHERE i.name = 'Республиканская школа-интернат №1' AND c.name = 'Канцелярия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Тетради (12 л.)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Матрасы', 'Ортопедические детские матрасы 160x70', 'шт', 40, 5, 'medium'
FROM institutions i, categories c WHERE i.name = 'Республиканская школа-интернат №1' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Матрасы' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХУДЖАНДА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Подгузники (Size 4)', 'Для детей от 9 до 14 кг', 'упак', 60, 10, 'high'
FROM institutions i, categories c WHERE i.name = 'Областной детский дом г. Худжанда' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Подгузники (Size 4)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Детское питание (смесь)', 'Гипоаллергенные смеси для новорожденных', 'банка', 100, 20, 'high'
FROM institutions i, categories c WHERE i.name = 'Областной детский дом г. Худжанда' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Детское питание (смесь)' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХОРОГА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Обогреватели', 'Масляные радиаторы для зимнего периода', 'шт', 15, 3, 'high'
FROM institutions i, categories c WHERE i.name = 'Центр социальной поддержки "Нур"' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Обогреватели' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, 'Цемент М500', 'Для ремонта пандуса у входа', 'мешок', 20, 0, 'low'
FROM institutions i, categories c WHERE i.name = 'Центр социальной поддержки "Нур"' AND c.name = 'Строительные материалы'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = 'Цемент М500' AND institution_id = i.id);