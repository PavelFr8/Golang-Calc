# Golang-Calc

## Функциональность

- Поддержка арифметических операций: сложение, вычитание, умножение, деление.
- Обработка ошибок:
  - Неверное выражение (422 Unprocessable Entity).
  - Ошибка выполнения (500 Internal Server Error).
  - Деление на ноль (500 Internal Server Error).
- Логирование HTTP-запросов и ошибок с использованием библиотеки [zap](https://github.com/uber-go/zap).

## Примеры использования

### Успешный запрос

```bash
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+2*2\"}"
```

#### Ответ

```json
{
  "result": "6.000000"
}
```

### Ошибка 422: Неверное выражение

```bash
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"2+\"}"
```

#### Ответ

```json
{
  "error": "Expression is not valid"
}
```

### Ошибка 500: Деление на ноль

```bash
curl --location "http://localhost:8080/api/v1/calculate" --header "Content-Type: application/json" --data "{\"expression\": \"999/0\"}"
```

#### Ответ

```json
{
  "error": "Internal server error"
}
```

## Инструкция по запуску

### Требования

- Go версии 1.23.0 или выше

### Установка и запуск

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/PavelFr8/Golang-Calc.git
   cd Golang-Calc
   ```

2. Запустите сервис:

   ```bash
   go run ./cmd/main.go
   ```

3. Сервис будет запущен на `http://localhost:8080`.

### Переменные окружения

- `PORT` — порт для запуска сервера. Если переменная не указана, используется значение по умолчанию: `8080`.

## Тестирование

Тесты для http-сервера находятся в файле `internal/application/application_test.go`.

Тесты для функции вычисления математических выражений находятся в файле `pkg/calculation/calculation_test.go`.

Для запуска тестов выполните:

```bash
go test ./...
```