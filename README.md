# service-chat

# Запуск сервиса
### 1. Поднимаем контейнеры с сервисом и базой данных
```bash
make up-service
```

# Debug режим

### 1. Для отладки необходимо использовать файлы Dockerfile.debug и compose.debug.yaml
### 2. Запустите файл compose.debug.yaml прямо с зеленой кнопки в IDE
### 3. Необходимо настроить запуск в IDE согласно этим двум инструкциям:
https://blog.jetbrains.com/go/2020/05/06/debugging-a-go-application-inside-a-docker-container/
https://blog.jetbrains.com/go/2020/05/08/running-go-applications-using-docker-compose-in-goland/
### 4. Также в файле local.yaml необходимо указать параметр для базы данных - host: clair-postgres-debug

# Swagger документация
### 1. Можно посмотреть по ссылке:
http://localhost:9000/swagger/index.html/
### 2. Генерация схемы - swag init -g cmd/service/main.go

# Структура проекта
## 1. cmd - точка входа в программу, здесь лежит файл main
## 2. config - здесь лежит файл local.yaml со всеми параметрами конфигурации сервиса и базы данных
## 3. internal - вся внутрення кухня  
### 3.1 /config - считываем файл local.yaml и перекладываем параметры в структуру для дальнейшего использования
### 3.2 /db - всё для работы с базой данных  
#### 3.2.1 db/schema - файлы миграции базы postgres (создание структуры таблиц и удаление)
#### 3.2.2 db/db - здесь интерфейсы для слоя базы данных и конструктор
#### 3.2.3 db/postgres - подключение к базе postgres  
### 3.3 /encryption - шифрование
### 3.4. /entity - сущности нашего приложения, с которыми мы можем работать на всех слоях приложения
### 3.5 /handler - здесь живёт слой обработки запросов
### 3.6 /logger - всё для логирования
### 3.7 /router - маршрутизатор запросов
### 3.8 /service - слой бизнес-логики
#### 3.8.1 service/service - здесь интерфейсы для слоя бизнес-логики нашего приложения
## 4. server/server - всё для запуска и остановки сервера

# План проекта

```
// TODO: init config: cleanenv
// библиотека cleanenv удобная минималистичная библиотека в отличии viper или cobra, в ней есть всё необходимое
// умеет читать из всех популярных файлов yaml, json, toml, .env и др.
// также удобно использовать struct tags, можем задавать required, default значение и др.

// TODO: init logger: slog
// библиотека slog является стандартной с версии go 1.21 и она самая актуальная

// TODO: init storage (db): postgres
// наверное самая популярная реляционная база данных

// TODO: init router: chi, chi render
// удобный, минималистичный, активно развивается, совместим с http/net стандартным пакетом
// chi render это один из пакетов chi для работы с структурами запросов и ответов

// TODO: handler sign-up и sign-in
// реализуем handlers для регистрации и авторизации с помощью JWT

// TODO: handlers for work service
// реализуем все остальные handlers для работы основной логики сервиса

// TODO: schema swagger
// делаем свагер схему

// TODO: tests
// делаем unit тесты

// TODO: run server
// просто запускаем сервер
```

