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

INSERT INTO vacancies (title, description, type, experience, workload)
SELECT 'SMM Специалист', 'Мы ищем волонтера для ведения наших социальных сетей. Рассказывайте о наших проектах и мероприятиях, формируйте бренд Hadaf. Очень коротко, но живо.', 'Волонтёрство', 'От 1 года', '1 час в неделю'
WHERE NOT EXISTS (SELECT 1 FROM vacancies WHERE title = 'SMM Специалист');

INSERT INTO events (title, description, event_date, institution_id, creator_id, status)
SELECT 
    'Благотворительная ярмарка и концерт', 
    'Организуем ярмарку поделок детей, вырученные средства пойдут на закупку медицинского оборудования.', 
    NOW() + INTERVAL '7 days', 
    i.id, 
    u.id, 
    'approved'
FROM institutions i, users u 
WHERE i.name = '[ТЕСТ] Центр социальной поддержки "Нур"' 
AND u.role = 'superadmin'
AND NOT EXISTS (SELECT 1 FROM events WHERE title = 'Благотворительная ярмарка и концерт');

INSERT INTO events (title, description, event_date, institution_id, creator_id, status)
SELECT 
    'Субботник по уборке территории', 
    'Собираемся дружной командой волонтёров и приводим в порядок территорию вокруг учреждения. Приносите перчатки!', 
    NOW() + INTERVAL '14 days', 
    i.id, 
    u.id, 
    'approved'
FROM institutions i, users u 
WHERE i.name = '[ТЕСТ] Дом-интернат для пожилых "Мехрубон"' 
AND u.role = 'superadmin'
AND NOT EXISTS (SELECT 1 FROM events WHERE title = 'Субботник по уборке территории');

-- Seed данных команды
INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Сиёвуш Хамидов', 'Founder & CTO', '/team/siyovush_hamidov.jpg', 'Хочу, чтобы помощь была прозрачной, адресной и по-человечески тёплой.', 'https://t.me/siyovush_hamidov', 'https://www.linkedin.com/in/siyovush-hamidov/', 1
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Сиёвуш Хамидов');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Саидмехродж Шукурзода', 'Backend Разработчик', 'Пишу код, который стоит за каждой доставленной помощью.', 'https://t.me/s_mehroj', 'https://www.linkedin.com/in/shukurzodasaidmehroj/', 2
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Саидмехродж Шукурзода');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Манучехр Олимов', 'Backend Разработчик', 'Верю, что надёжная система — это тоже форма заботы.', 'https://t.me/olimov_manu', 'https://www.linkedin.com/in/manuchehr-olimov/', 3
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Манучехр Олимов');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Раббони Бафоев', 'Backend Разработчик', 'Каждая строка кода — шаг к тому, чтобы кто-то получил помощь вовремя.', 'https://t.me/AB001G', 'https://www.linkedin.com/in/rabboni-bafoev-45052131b/', 4
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Раббони Бафоев');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Илёс Джабборов', 'Frontend Разработчик', 'Делаю так, чтобы помогать было просто и приятно.', 'https://t.me/si_senor9', '', 5
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Илёс Джабборов');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Абукахор Якуби', 'Mobile Разработчик', 'Добро должно быть в кармане — буквально.', 'https://t.me/yakubiam', '', 6
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Абукахор Якуби');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Алиджон Тошев', 'QA Инженер', 'Моя задача — чтобы ни одна ошибка не встала между помощью и человеком.', 'https://t.me/alijon07_t', '', 7
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Алиджон Тошев');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Анна Сардаковская', 'UX/UI Дизайнер', 'Хороший интерфейс — когда бабушка разберётся без инструкции.', 'https://t.me/netta_sardakovskaia', 'https://www.linkedin.com/in/anna-sardakovskaya/', 8
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Анна Сардаковская');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Елизавета Сирошенко', 'UX/UI Дизайнер', 'Красота — это не про декор. Это про ощущение, что тебе здесь рады.', 'https://t.me/liza_oper', 'https://www.linkedin.com/in/elizaveta-siroshenko-b51350327/', 9
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Елизавета Сирошенко');

INSERT INTO team_members (full_name, role, quote, telegram, linkedin, sort_order)
SELECT 'Фарзона Ахмедова', 'Маркетолог', 'Рассказываю истории, которые вдохновляют людей помогать.', 'https://t.me/farzona_Akhmedova', 'https://www.linkedin.com/in/farzona-akhmedova-952082213/', 10
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Фарзона Ахмедова');