# Веб-калькулятор на Go

Этот проект представляет собой веб-калькулятор для подсчёта арифметических выражений. Сервис принимает выражение через HTTP-запрос и возвращает результат вычисления.

## Используемые технологии

- **Go** — основной язык программирования.
- **mux** — для маршрутизации HTTP-запросов.
- **logrus** — для логирования.

---

## Основные возможности

- **HTTP API**: Принимает математическое выражение и возвращает результат вычислений.
- **Обработка ошибок**: Возвращает сообщения об ошибках с соответствующими HTTP-кодами.
- **Логирование**: Логирование запросов и ошибок осуществляется в консоль с использованием библиотеки `logrus`.

---

## Установка и запуск

### Требования

- Go версии 1.23.4.

### Инструкция

1. Склонируйте репозиторий:

```bash
   git clone https://github.com/dimasmir03/Calc_GO
   cd Calc_GO
```

2. Установите зависимости:

```bash
   go mod tidy
```

3. Запустите сервер:

## C указанием порта через переменную окружения

### Linux, Mac, Windows(Git Bash, WSL)

```bash
   export PORT=8787 && go run ./cmd/calc/...
```

### Windows

#### PowerShell

```ps
   $env:PORT=8787; go run ./cmd/calc/...
```

#### Command Promt

- Сервер будет запущен на порту 8080

```cmd
   go run ./cmd/calc/...
```

После запуска сервер будет доступен по адресу `http://localhost:<PORT>`, где `<PORT>` — это порт, который вы указали в переменной окружения.

---

Тестирование

Для запуска тестов выполните следующую команду:

```bash
go test ./...
```

Тесты проверяют корректность работы API, обработку различных ошибок и корректность выполнения математических операций.

Пример результата тестирования:

ok github.com/dimasmir03/Calc_GO/internal/application 0.272s
ok github.com/dimasmir03/Calc_GO/pkg/calculation 0.181s

---

## API

### Единственный эндпоинт

#### `POST /api/v1/calculate`

Эндпоинт принимает JSON-запрос с математическим выражением и возвращает результат вычислений или сообщение об ошибке.

**Пример запроса:**

```json
{
	"expression": "2+2*2"
}
```

**Пример ответа:**

- Успешный результат:

  ```json
  {
  	"result": "6"
  }
  ```

- Ошибка (Некорректное выражение):

  ```json
  {
  	"error": "Expression is not valid"
  }
  ```

  **HTTP-код:** 422

- Ошибка (Внутренняя ошибка сервера):
  ```json
  {
  	"error": "Internal server error"
  }
  ```
  **HTTP-код:** 500

---

## Примеры использования

#### Успешные вычисления

**Запрос:**

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2+2*2"}'
```

**Ответ:**

```json
{
	"result": "6"
}
```

- Примечание: Для отображения заголовка ответа, включая код статуса, используйте параметр -i:

```bash
curl -i --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2+2*2"}'
```

#### Ошибки

1. **Некорректное выражение:**

   **Запрос:**

   ```bash
   curl -i --location 'localhost:8080/api/v1/calculate' \
   --header 'Content-Type: application/json' \
   --data '{"expression": "2+2**2"}'
   ```

   **Ответ:**

   ```json
   {
   	"error": "invalid expression"
   }
   ```

   **HTTP-код:** 422

2. **Деление на ноль:**

   **Запрос:**

   ```bash
   curl --location 'localhost:8080/api/v1/calculate' \
   --header 'Content-Type: application/json' \
   --data '{"expression": "10/0"}'
   ```

   **Ответ:**

   ```json
   {
   	"error": "division by zero"
   }
   ```

   **HTTP-код:** 422

3. **Неправильный символ в выражении:**

**Запрос:**

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2+a"}'
```

**Ответ:**

```json
{
	"error": "invalid character"
}
```

**HTTP-код:** 422

4.  **Несовпадающие круглых скобок:**

**Запрос:**

```bash
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2*(2+2"}'
```

**Ответ:**

```json
{
	"error": "mismatched parentheses"
}
```

**HTTP-код:** 422

---

## Логирование

Все запросы и ошибки логируются в консоль с использованием библиотеки `logrus`. Логи содержат информацию о типе запроса, времени его обработки и ошибках.

Пример логов:

```
INFO[2024-12-20 10:00:00] HTTP Запрос POST /api/v1/calculate 200ms
ERROR[2024-12-20 10:01:00] HTTP Запрос POST /api/v1/calculate 500ms
INFO[2024-12-20 10:02:00] HTTP Запрос POST /api/v1/calculate 100ms
ERROR[2024-12-20 10:03:00] Internal server error
ERROR[2024-12-20 10:04:00] invalid expression
ERROR[2024-12-20 10:05:00] mismatched parentheses
```
