# SnipAndNeat

Проект для работы с API Ozon и анализа транзакций.

## Структура проекта

### Основные директории

- `app/` - Основной код приложения
  - `api/` - HTTP API эндпоинты
    - `get_list_transactions.go` - Получение списка транзакций
    - `get_sum_services.go` - Получение суммы услуг
    - `get_viento_product.go` - Получение информации о продукте
    - `health.go` - Проверка здоровья сервиса
    - `http.go` - Основной HTTP сервер
    - `post_update_items_barcode.go` - Обновление штрих-кодов товаров
  - `config/` - Конфигурация приложения
    - `config.go` - Основной конфигурационный файл
  - `service/` - Бизнес-логика сервисов
    - `ozon/` - Сервис для работы с Ozon API
      - `errors.go` - Определения ошибок
      - `list_product.go` - Работа со списком продуктов
      - `list_transactions.go` - Работа со списком транзакций
      - `metrics.go` - Метрики Prometheus
      - `ozon.go` - Основной клиент Ozon API
      - `profit_analytics.go` - Анализ прибыли
      - `product_info.go` - Информация о продуктах
      - `sum_proficiency.go` - Расчет прибыльности
      - `sum_services.go` - Расчет услуг
      - `udate_items.go` - Обновление товаров
    - `metrics/` - Метрики приложения
    - `scheded/` - Планировщик задач
    - `mailer/` - Сервис отправки почты
  - `strorage/` - Работа с хранилищем данных
    - `db.go` - Основной интерфейс базы данных
    - `items.go` - Работа с товарами
    - `operations.go` - Работа с операциями
    - `posting.go` - Работа с отправками
    - `services.go` - Работа с услугами
    - `viento_products.go` - Работа с продуктами Viento
    - `migrations/` - Миграции базы данных
  - `application.go` - Основная логика приложения
  - `launcher.go` - Запуск приложения

- `cmd/` - Исполняемые файлы
  - `SnipAndNeat/` - Основное приложение
    - `main.go` - Точка входа

- `common/` - Общие утилиты и компоненты

### Конфигурационные файлы

- `docker-compose.yaml` - Конфигурация Docker Compose
- `prometheus.yml` - Конфигурация Prometheus
- `grafana-dashboard.json` - Дашборд Grafana
- `.env` - Переменные окружения
- `.ogen.yml` - Конфигурация генерации OpenAPI
- `openapi.yml` - OpenAPI спецификация

### Мониторинг

Проект включает в себя систему мониторинга на основе:
- Prometheus для сбора метрик
- Grafana для визуализации

Основные метрики:
- Количество и сумма транзакций
- Прибыль по дням и категориям
- Комиссии по типам
- Возвраты и их стоимость

## Запуск проекта

1. Убедитесь, что установлен Docker и Docker Compose
2. Скопируйте `.env.example` в `.env` и настройте переменные окружения
3. Запустите проект:
```bash
docker-compose up -d
```

## Доступные сервисы

- Основное приложение: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (логин: admin, пароль: admin)

## API Endpoints

- GET /health - Проверка здоровья сервиса
- GET /api/v1/transactions - Получение списка транзакций
- GET /api/v1/services/sum - Получение суммы услуг
- GET /api/v1/products/viento - Получение информации о продукте
- POST /api/v1/items/barcode - Обновление штрих-кодов товаров
