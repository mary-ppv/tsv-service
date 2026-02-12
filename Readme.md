# tsv-service

Сервис для импорта TSV-файлов, сохранения данных в PostgreSQL и генерации отчётов по unit_guid.

Стек:
•	Go
•	Gin
•	PostgreSQL
•	SQLBoiler

## Описание работы

Сервис отслеживает директорию INPUT_DIR. При появлении TSV-файла:
1.	Файл регистрируется в таблице tsv_files.
2.	Статус файла меняется на processing.
3.	Строки файла парсятся.
4.	Валидные строки сохраняются в таблицу tsv_records.
5.	Ошибки сохраняются в таблицу parse_errors.
6.	Для каждого unit_guid формируется отчёт.
7.	Статус файла меняется на done или failed.

⸻

## Требования
•	Go (1.21+)
•	Docker и docker compose
•	PostgreSQL client (psql)
•	sqlboiler и sqlboiler-psql (для генерации моделей)

⸻

## Переменные окружения

- Обязательные переменные:

POSTGRES_DSN
строка подключения к базе
пример:
postgres://postgres:postgres@localhost:5433/tsv?sslmode=disable

INPUT_DIR
директория для входящих TSV-файлов
пример:
./var/input

OUTPUT_DIR
директория для сформированных отчётов
пример:
./var/output

- Необязательные:

HTTP_ADDR
по умолчанию :8080

POLL_INTERVAL_SEC
по умолчанию 5

WORKERS
по умолчанию 4

⸻

## Файл .env.example

HTTP_ADDR=:8080
POSTGRES_DSN=postgres://postgres:postgres@localhost:5433/tsv?sslmode=disable
INPUT_DIR=./var/input
OUTPUT_DIR=./var/output
POLL_INTERVAL_SEC=5
WORKERS=4

В репозиторий кладётся .env.example.
Файл .env создаётся локально и не коммитится.

⸻

## Запуск PostgreSQL через Docker

docker compose up -d db

Проверить подключение:

psql “postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable” -c ‘\l’

Если базы tsv нет, создать:

psql “postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable” -c “CREATE DATABASE tsv;”

⸻

## Миграции

psql “$POSTGRES_DSN” -f db/migrations/000001_init.sql

⸻

## Генерация моделей SQLBoiler

После применения миграций:

make sqlboiler

Важно: версия sqlboiler в go.mod и установленный бинарник должны совпадать.

⸻

## Запуск сервиса

set -a; source .env; set +a; 
make run

Сервис стартует на HTTP_ADDR (по умолчанию :8080).

⸻

## API

Базовый путь: /api/admin

Получить записи по unit_guid:
GET /api/admin/units/:unit_guid/records?page=1&limit=20

Получить отчёты по unit_guid:
GET /api/admin/units/:unit_guid/reports?page=1&limit=20

⸻

Проверка обработки TSV
1.	Поместить файл в INPUT_DIR.
2.	Вотчер обнаружит файл.
3.	Воркер обработает его.
4.	Данные появятся в базе.
5.	Отчёты будут созданы в OUTPUT_DIR.


⸻

## Структура проекта

cmd/ - запуск приложения

db/ - описание миграции

internal/: 
    app/        – инициализация приложения
    repository/ – слой доступа к данным (SQLBoiler)
    services/   – бизнес-логика
    worker/     – обработка файлов и генерация отчётов
    transport/  – HTTP-слой (Gin)

Архитектура разделена на уровни:
repository → service → controller
