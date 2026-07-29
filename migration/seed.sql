-- 1. КАТЕГОРИИ (Их можно оставить без приставки "Тест", так как это справочник)
INSERT INTO categories (name) VALUES 
('Продукты питания'), ('Одежда и обувь'), ('Медикаменты'), 
('Бытовая химия'), ('Канцелярия'), ('Строительные материалы'), ('Мебель')
ON CONFLICT (name) DO NOTHING;

-- 2. УЧРЕЖДЕНИЯ (Таджикистан) - Добавлена приставка [ТЕСТ]
INSERT INTO institutions (name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude, wards_count) VALUES 
(
    '[ТЕСТ] Тестовое Учреждение 1', 'Elderly', 'Турсунзаде', 'РРП', 
    'ул. Парк Победы, 45', '+992 900 11 22 33', 'test_1@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Государственное учреждение для пожилых людей и инвалидов.', '08:00 - 17:00', 38.5134, 68.2256, 120
),
(
    '[ТЕСТ] Тестовое Учреждение 2', 'Children', 'Душанбе', 'Душанбе', 
    'ул. Айни, 126', '+992 918 44 55 66', 'test_2@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Школа-интернат для детей-сирот и детей из малообеспеченных семей.', '24/7', 38.5598, 68.7870, 350
),
(
    '[ТЕСТ] Тестовое Учреждение 3', 'Children', 'Худжанд', 'Согд', 
    'ул. Сирдарья, 12', '+992 927 00 11 22', 'test_3@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Учреждение для детей дошкольного возраста.', '09:00 - 18:00', 40.2833, 69.6167, 85
),
(
    '[ТЕСТ] Тестовое Учреждение 4', 'Disabled', 'Хорог', 'ГБАО', 
    'ул. Шош Саидов, 5', '+992 935 88 99 00', 'test_4@hadaf.tj', 
    '[ТЕСТОВЫЕ ДАННЫЕ] Помощь людям с ограниченными возможностями в горном регионе.', '10:00 - 16:00', 37.4894, 71.5529, 45
)
ON CONFLICT (name) DO NOTHING;

-- 3. НУЖДЫ (Генерация через безопасный поиск ID)

-- --- НУЖДЫ ДЛЯ БАТОШ ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Мука (1 сорт)', '[ТЕСТОВЫЙ ЗАПРОС] Для выпечки свежего хлеба', 'кг', 500, 50, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 1' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Мука (1 сорт)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Тонометры', '[ТЕСТОВЫЙ ЗАПРОС] Автоматические тонометры на плечо', 'шт', 10, 2, 'medium'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 1' AND c.name = 'Медикаменты'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Тонометры' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Стиральный порошок', '[ТЕСТОВЫЙ ЗАПРОС] Для автоматических машин', 'кг', 100, 0, 'low'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 1' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Стиральный порошок' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ДУШАНБЕ (ШКОЛА №1) ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Тетради (12 л.)', '[ТЕСТОВЫЙ ЗАПРОС] В клетку и в линейку для начальных классов', 'шт', 2000, 450, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 2' AND c.name = 'Канцелярия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Тетради (12 л.)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Матрасы', '[ТЕСТОВЫЙ ЗАПРОС] Ортопедические детские матрасы 160x70', 'шт', 40, 5, 'medium'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 2' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Матрасы' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХУДЖАНДА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Подгузники (Size 4)', '[ТЕСТОВЫЙ ЗАПРОС] Для детей от 9 до 14 кг', 'упак', 60, 10, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 3' AND c.name = 'Бытовая химия'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Подгузники (Size 4)' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Детское питание (смесь)', '[ТЕСТОВЫЙ ЗАПРОС] Гипоаллергенные смеси для новорожденных', 'банка', 100, 20, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 3' AND c.name = 'Продукты питания'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Детское питание (смесь)' AND institution_id = i.id);

-- --- НУЖДЫ ДЛЯ ХОРОГА ---
INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Обогреватели', '[ТЕСТОВЫЙ ЗАПРОС] Масляные радиаторы для зимнего периода', 'шт', 15, 3, 'high'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 4' AND c.name = 'Мебель'
AND NOT EXISTS (SELECT 1 FROM needs WHERE name = '[ТЕСТ] Обогреватели' AND institution_id = i.id);

INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
SELECT i.id, c.id, '[ТЕСТ] Цемент М500', '[ТЕСТОВЫЙ ЗАПРОС] Для ремонта пандуса у входа', 'мешок', 20, 0, 'low'
FROM institutions i, categories c WHERE i.name = '[ТЕСТ] Тестовое Учреждение 4' AND c.name = 'Строительные материалы'
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
WHERE i.name = '[ТЕСТ] Тестовое Учреждение 4' 
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
WHERE i.name = '[ТЕСТ] Тестовое Учреждение 3' 
AND u.role = 'superadmin'
AND NOT EXISTS (SELECT 1 FROM events WHERE title = 'Субботник по уборке территории');

INSERT INTO events (title, description, event_date, institution_id, creator_id, status)
SELECT 
    'Сбор зимней одежды и обуви', 
    'Собираем теплые куртки, шапки и зимнюю обувь для воспитанников к предстоящим холодам. Нужно собрать 50 комплектов (рост 120-150 см). Приносите новые вещи или в отличном состоянии.', 
    NOW() + INTERVAL '5 days', 
    i.id, 
    u.id, 
    'approved'
FROM institutions i, users u 
WHERE i.name = '[ТЕСТ] Тестовое Учреждение 3' 
AND u.role = 'superadmin'
AND NOT EXISTS (SELECT 1 FROM events WHERE title = 'Сбор зимней одежды и обуви');

INSERT INTO events (title, description, event_date, institution_id, creator_id, status)
SELECT 
    'Мастер-класс по рисованию и сбор красок', 
    'Ищем волонтёров, умеющих рисовать! Организуем творческий вечер для детей. Собираем материалы: 30 альбомов, гуашь, кисточки и цветные карандаши. Подарим детям праздник.', 
    NOW() + INTERVAL '10 days', 
    i.id, 
    u.id, 
    'approved'
FROM institutions i, users u 
WHERE i.name = '[ТЕСТ] Тестовое Учреждение 2' 
AND u.role = 'superadmin'
AND NOT EXISTS (SELECT 1 FROM events WHERE title = 'Мастер-класс по рисованию и сбор красок');

-- Seed данных команды
INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Сиёвуш Хамидов', 'Founder & CTO', '/team/siyovush_hamidov.webp', 'Хочу, чтобы помощь была прозрачной, адресной и по-человечески тёплой.', 'https://t.me/siyovush_hamidov', 'https://www.linkedin.com/in/siyovush-hamidov/', 1
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Сиёвуш Хамидов');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Саидмехродж Шукурзода', 'Lead Backend Разработчик', '/team/saidmehroj_shukurzoda.webp', 'Пишу код, который стоит за каждой доставленной помощью.', 'https://t.me/s_mehroj', 'https://www.linkedin.com/in/shukurzodasaidmehroj/', 2
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Саидмехродж Шукурзода');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Илёс Джабборов', 'Lead Frontend Разработчик', '/team/ilyos_djabborov.webp', 'Делаю так, чтобы помогать было просто и приятно.', 'https://t.me/si_senor9', '', 3
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Илёс Джабборов');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Абдукахор Якуби', 'Lead Mobile Разработчик', '/team/abdukakhor_yakubi.webp', 'Добро должно быть в кармане — буквально.', 'https://t.me/yakubiam', '', 4
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Абдукахор Якуби');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Манучехр Олимов', 'Backend Разработчик', '/team/manuchehr_olimov.webp', 'Верю, что надёжная система — это тоже форма заботы.', 'https://t.me/olimov_manu', 'https://www.linkedin.com/in/manuchehr-olimov/', 7
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Манучехр Олимов');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Раббони Бафоев', 'Backend Разработчик', '/team/rabboni_bafoev.webp', 'Каждая строка кода — шаг к тому, чтобы кто-то получил помощь вовремя.', 'https://t.me/AB001G', 'https://www.linkedin.com/in/rabboni-bafoev-45052131b/', 8
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Раббони Бафоев');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Шарипов Шухрат', 'Frontend Разработчик', '/team/sharipov_shukhrat.webp', 'Интерфейс должен быть не только красивым, но и помогать делать добро без лишних кликов.', 'https://t.me/sharipovsh13', '', 10
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Шарипов Шухрат');

INSERT INTO team_members (full_name, role, photo_url, quote, telegram, linkedin, sort_order)
SELECT 'Алиджон Тошев', 'QA Инженер', '/team/alijon_toshev.webp', 'Моя задача — чтобы ни одна ошибка не встала между помощью и человеком.', 'https://t.me/alijon07_t', '', 12
WHERE NOT EXISTS (SELECT 1 FROM team_members WHERE full_name = 'Алиджон Тошев');

-- Запрещённые и рекомендуемые вещи для учреждений
UPDATE institutions SET
  prohibited_items = 'Алкогольные напитки и табачные изделия
Просроченные продукты питания
Просроченные или неупакованные лекарства
Колюще-режущие предметы
Скоропортящиеся продукты без сертификатов
Бытовая химия без маркировки
Одежда и постельное бельё в неудовлетворительном состоянии
Электроприборы без сертификатов безопасности',
  recommended_items = 'Средства личной гигиены (мыло, шампунь, зубная паста)
Постельное бельё и полотенца (новые)
Тонометры и глюкометры
Подгузники для взрослых
Тёплая одежда по сезону (халаты, носки, тапочки)
Нескоропортящиеся продукты (крупы, масло, чай, сахар)
Настольные игры и книги с крупным шрифтом
Моющие средства (стиральный порошок, средство для мытья посуды)'
WHERE name = '[ТЕСТ] Тестовое Учреждение 1';

UPDATE institutions SET
  prohibited_items = 'Сладости и конфеты в большом количестве
Мягкие игрушки бывшие в употреблении
Острые и колющие предметы (ножницы без закруглённых концов, ножи)
Продукты питания без сертификатов
Лекарственные препараты (передаются только через медперсонал)
Электронные гаджеты (телефоны, планшеты) без согласования с администрацией
Одежда и обувь в изношенном состоянии
Продукты, содержащие аллергены, без маркировки',
  recommended_items = 'Канцелярские принадлежности (тетради, ручки, карандаши, альбомы)
Одежда и обувь по сезону (новые, размеры от 120 до 170 см)
Настольные и развивающие игры
Спортивный инвентарь (мячи, скакалки, бадминтон)
Книги для детей и подростков
Средства гигиены (зубные щётки, мыло, шампунь)
Постельное бельё (полуторное, новое)
Фрукты и соки в заводской упаковке
Рюкзаки и школьные сумки'
WHERE name = '[ТЕСТ] Тестовое Учреждение 2';

UPDATE institutions SET
  prohibited_items = 'Детское питание и смеси без сертификатов качества
Ходунки и коляски бывшие в употреблении без дезинфекции
Мягкие игрушки б/у (гигиенические нормы)
Мелкие предметы и игрушки с мелкими деталями (опасность проглатывания)
Продукты с истекшим сроком годности
Бытовая химия с резким запахом
Непромаркированные лекарства и витамины',
  recommended_items = 'Гипоаллергенные детские смеси (по согласованию с врачом)
Подгузники (размеры 3, 4, 5)
Детская одежда (от 0 до 5 лет, новая)
Развивающие игрушки с сертификатами безопасности
Влажные салфетки и детские кремы
Пелёнки и одеяла (новые)
Детские книжки с картинками
Моющие средства для детской посуды и стирки'
WHERE name = '[ТЕСТ] Тестовое Учреждение 3';

UPDATE institutions SET
  prohibited_items = 'Электроприборы без сертификатов безопасности
Продукты, содержащие распространённые аллергены, без маркировки
Медицинское оборудование без документации
Алкоголь и табачные изделия
Острые и тяжёлые предметы без упаковки
Бытовая химия с агрессивными компонентами
Мебель без устойчивых креплений (риск опрокидывания)',
  recommended_items = 'Тёплая верхняя одежда и обувь (горный регион, размеры по запросу)
Обогреватели с сертификатом безопасности
Средства реабилитации (костыли, трости, противоскользящие коврики)
Средства личной гигиены
Нескоропортящиеся продукты (крупы, консервы, масло, мука)
Постельное бельё и тёплые одеяла
Батарейки, фонарики и аккумуляторы
Настольные игры и материалы для рукоделия
Строительные материалы для ремонта пандусов и входов'
WHERE name = '[ТЕСТ] Тестовое Учреждение 4';