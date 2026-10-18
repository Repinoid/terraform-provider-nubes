> **Инструкция из личного кабинета**
> Это руководство по сервису из личного кабинета Nubes (ЛК). Текст приведён
> без изменений, «как есть».
>
# Инструкция по управлению кластерами и пользователями ClickHouse

---

## Общая информация

ClickHouse — это колоночная система управления базами данных, предназначенная для обработки аналитических запросов в реальном времени. Платформа позволяет создавать, удалять, приостанавливать и возобновлять кластеры ClickHouse, а также управлять базами данных и пользователями.

Кластер создаётся через оператор Altinity (`https://github.com/Altinity/clickhouse-operator`) в ресурсной платформе типа `k8s`

Кластеры могут создаваться шардированными с произвольным количеством реплик. Для обеспечения отказоустойчивости и координации обязательно включение параметра `Создать ClickHouse-Keeper`.

Кластера создаются с обязательным бекапом. Для этого необходимо предварительно создать услугу `S3 Object Storage`.
---

## Доступные операции

* **create**  — Создание кластера
* **delete**  — Удаление кластера
* **suspend** — Приостановка работы кластера
* **resume**  — Возобновление работы кластера
* **create_database** — Создание базы данных
* **delete_database** — Удаление базы данных
* **create_user** — Создание пользователя и настройка прав доступа к базе данных
* **delete_user** — Удаление пользователя

## Процесс создания и первого подключения

1. Выполнить операцию `create` - будет создан кластер
2. Выполнить операцию `create_database` - будет создана база данных
3. Выполнить операцию `create_user` - будет создан пользователь (выбрать доступ подключения к БД, созданной в п.2)

## Примеры подключения

### Подключение через HTTP-интерфейс

```bash
# Проверка доступности кластера
curl "http://<username>:<password>@<host>:8123/" -d "SELECT 1"

# Выполнение запроса
curl "http://<username>:<password>@<host>:8123/" -d "SHOW DATABASES"
```

## Простые кейсы
*Создать БД и обычной таблицы*
```clickhouse
CREATE TABLE database.table
(
    id UInt64,
    value String
)
ENGINE = MergeTree()
ORDER BY id;

INSERT INTO database.table SELECT number, toString(number) FROM numbers(10000000);

SELECT count() FROM database.table;
```
Базы и таблицы не реплицируются без специального ключа (`ON CLUSTER 'clickhousek8s'`)
В случае создания в стиле, описанным выше, необходимо будет выполнять действия на всех шардах/репликах

*Создание БД и шардированной таблицы*
```clickhouse
CREATE TABLE database.localUsers
(
    id UInt64,
    value String
)
ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE database.shardedTable
AS localUsers
ENGINE = Distributed('clickhousek8s', 'database', 'localUsers', id);

INSERT INTO shardedTable SELECT number, toString(number) FROM numbers(10000000);

SELECT count() FROM shardedTable;
```
Важный момент - в данном примере таблицы придётся также выполнять на всех нодах ClickHouse

*Создание реплицируемой и шардированной таблицы*
```clickhouse
CREATE TABLE IF NOT EXISTS database.localUsers ON cluster 'clickhousek8s'
(
    user_id UInt64,
    name String,
    email String
)
ENGINE = ReplicatedMergeTree(
    '/clickhouse/tables/{shard}/localUsers',
    '{replica}'
)
ORDER BY user_id;

CREATE TABLE IF NOT EXISTS database.users ON cluster 'clickhousek8s'
AS database.localUsers
ENGINE = Distributed('clickhousek8s', 'database', 'localUsers', rand());

INSERT INTO database.users VALUES (1, 'Alice', 'alice@example.com');
INSERT INTO database.users VALUES (2, 'Bob', 'bob@example.com');
INSERT INTO database.users VALUES (3, 'Charlie', 'charlie@example.com');
INSERT INTO database.users VALUES (4, 'David', 'david@example.com');
INSERT INTO database.users VALUES (5, 'Gosha', 'gosha@example.com');
INSERT INTO database.users VALUES (6, 'Eve', 'eve@example.com');
INSERT INTO database.users VALUES (7, 'Frank', 'frank@example.com');
INSERT INTO database.users VALUES (8, 'Grace', 'grace@example.com');
INSERT INTO database.users VALUES (9, 'Hannah', 'hannah@example.com');
INSERT INTO database.users VALUES (10, 'Ivan', 'ivan@example.com');

SELECT * FROM database.localUsers;
SELECT * FROM database.users;

SELECT COUNT(*) AS total_users       FROM database.users;
SELECT COUNT(*) AS total_users_local FROM database.localUsers;
```
