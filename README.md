# pvz-service

Это тестовое задание в рамках отбора на летнюю стажировку в компанию Avito на позицию Trainee Backend Engineer.

## TODO

- [x] Продумать архитектуру сервиса
- [x] Написать доменные сущности
- [x] Написать репозитории
- [x] Написать имплементации репозиториев
- [x] Написать миграции базы
- [ ] Реализовать все нужные `middleware` (`auth`, `logger`, `panic`)
- [ ] Реализовать все обработчики
- [ ] Добавить логи
- [ ] Добавить метрики
- [ ] Написать unit-тесты
- [ ] Написать интеграционные тесты
- [ ] Написать нагрузочные тесты
- [ ] Провести полное ручное тестирование всех сценариев
- [ ] Написать Dockerfile и docker-compose.yml

## Архитектура

В последнее время я придерживаюсь принципа **DRY**, поэтому использовал здесь **DDD** подход. Во главе всего стоит `domain` - доменные сущности и репозитории. API обращается к сервисам, сервисы "ходят" в домен и в хранилище.

Каждый сервис ограничен своим набором методов, которые он может выполнять с хранилищем.

На каждый набор ручек есть свой `Handler`, который содержит все нужные обработчики.

Я отказался от `go-playground/validator` и написал для каждого `<name>Params` свой метод `Validate()`, который валидирует поля этого `Params`. Это позволяет четко определить бизнес-правила валидации + быстрее, чем валидация через рефлексию. 

## Вопросы

1. Почему для `/register` в спецификации прописаны только `201` код и `400`? А если ошибка на стороне базы будет, то все равно `400` отдавать? Или если юзер уже существует, то почему не `StatusConflict`?