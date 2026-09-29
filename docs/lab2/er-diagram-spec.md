# ER-диаграмма — данные для StarUML

Создать: Model → добавить `ER Diagram` (или Data Model), в нём — 3 сущности
(Entity), связи через Relationship (1 — *).

## Таблица 1: communal_resources

| Поле | Тип | Ключ |
|---|---|---|
| id | SERIAL | PK |
| name | VARCHAR(100) | not null |
| description | VARCHAR(500) | |
| status | VARCHAR(15) | not null, default 'draft' |
| image_url | VARCHAR(255) | |
| video_url | VARCHAR(255) | |
| tariff_rate | NUMERIC(10,2) | |
| consumption_norm | NUMERIC(10,3) | |
| measurement_unit | VARCHAR(20) | |
| created_at | TIMESTAMP | not null |
| published_at | TIMESTAMP | nullable |
| creator_id | INT | FK → users.id, not null |

## Таблица 2: likes

| Поле | Тип | Ключ |
|---|---|---|
| id | SERIAL | PK |
| user_id | INT | FK → users.id, not null |
| communal_resource_id | INT | FK → communal_resources.id, not null |

Уникальный составной индекс: (user_id, communal_resource_id) — один
пользователь не может лайкнуть одну услугу дважды.

## Таблица 3: users

| Поле | Тип | Ключ |
|---|---|---|
| id | SERIAL | PK |
| login | VARCHAR(50) | unique, not null |
| password | VARCHAR(100) | not null |
| is_moderator | BOOLEAN | not null, default false |

## Связи

- `users` (1) → `communal_resources` (*) по `creator_id` — один пользователь
  создаёт много услуг
- `users` (1) → `likes` (*) по `user_id`
- `communal_resources` (1) → `likes` (*) по `communal_resource_id`

Каскадное удаление НЕ используется нигде (требование методички) — удаление
только логическое, через поле `status`.
