# service-chat

1. cmd - точка входа в программу, здесь лежит файл main
2. config - здесь лежит файл local.yaml со всеми параметрами конфигурации сервиса и базы данных
3. db - всё необходимое для работы базы данных
4. internal - вся внутрення кухня

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

