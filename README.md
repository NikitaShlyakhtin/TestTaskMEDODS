# Тестовое задание в компанию MEDODS

Это часть сервиса аутентификации, разработанная в качестве тестового задания на позицию Junior Backend Developer в компанию MEDODS.

## Используемые технологии

- Go
- JWT
- MongoDB

## Задание

Задача состоит в том, чтобы реализовать часть сервиса аутентификации c двумя REST-маршрутами:

1. Первый маршрут выдаёт пару токенов Access и Refresh для пользователя c указанным идентификатором (GUID) в параметре запроса.
2. Второй маршрут выполняет операцию обновления пары токенов Access и Refresh.

## Требования

- Access токен имеет тип JWT c алгоритмом SHA512. Хранение его в базе данных строго запрещено.
- Refresh токен имеет произвольный тип с форматом base64. Он хранится в базе данных исключительно в виде bcrypt хеша и должен быть защищен от изменений на стороне клиента и попыток повторного использования.
- Access и Refresh токены взаимосвязаны. Операцию Refresh для Access токена можно выполнить только с помощью Refresh токена, выданного вместе с ним.

## Установка

Чтобы установить и запустить проект, выполните следующие шаги:

1. Убедитесь, что у вас установлен Go на вашем компьютере.
2. Клонируйте репозиторий проекта:

    ```bash
    git clone <repository_url>
    ```
3. Создайте файл `.envrc` в корне проекта со следующим содержимым:

    ```bash
    export MONGODB_URI=<your_mongodb_connection_string>
    ```

    Замените `<your_mongodb_connection_string>` на вашу строку подключения к MongoDB.

4. Установите зависимости проекта:

    ```bash
    go mod tidy
    ```

6. Соберите и запустите проект с помощью команды из Makefile:

    ```bash
    make run/auth
    ```

## Использование

### Маршрут 1: Выдача токенов

- **Метод**: POST
- **Путь**: `/auth/token`
- **Описание**: Выдаёт пару токенов Access и Refresh для пользователя с указанным GUID.
- **Параметры запроса**: GUID пользователя в формате JSON: `{"guid": string}`
- **Вывод**: При успешном выполнении запроса маршрут возвращает токены в формате JSON: `{"tokens": {"access": string, "refresh": string}}`

### Маршрут 2: Обновление токенов

- **Метод**: POST
- **Путь**: `/auth/token/refresh`
- **Описание**: Выполняет операцию обновления пары токенов Access и Refresh.
- **Параметры запроса**: Refresh токен в формате JSON: `{"refresh_token": string}`
- **Вывод**: При успешном выполнении запроса маршрут возвращает токены в формате JSON: `{"tokens": {"access": string, "refresh": string}}`

## Примечание

В рамках проекта были реализованы следующие особенности:

1. **Проверка запросов пользователя**: Все входящие HTTP запросы проходят валидацию. Валидация включает в себя проверку формата входящих данных, а также обработку ошибок запросов, если они происходят.

2. **Отлов и отправка ошибок в JSON формате**: При возникновении ошибок сервер возвращает информацию об ошибке в формате JSON

3. **Graceful Shutdown**: Сервер корректно завершает работу, обеспечивая корректное завершение всех полученных запросов.

4. **Соблюдение принципов SOLID**: При разработке проекта я старался соблюдать принципы SOLID и выстраивал архитектуру проекта с учетом этих принципов.