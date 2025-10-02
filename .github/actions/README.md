# GitHub Actions для dSheltie

Этот каталог содержит универсальный composite action для оптимизированной работы с Go проектами.

## 🚀 Единый Universal Action

### `setup-go-project`

Универсальная настройка Go проекта с оптимизированным кэшированием для всех сценариев использования.

**Использование:**

```yaml
- name: Setup Go project
  uses: ./.github/actions/setup-go-project
  with:
    go-version: "1.24.1"
    cache-key-suffix: "my-suffix"
    skip-generate-networks: "false"
    skip-download-deps: "false"
    checkout-repo: "true"
    fetch-depth: "1"
```

**Параметры:**

- `go-version` - версия Go (по умолчанию: "1.24.1")
- `cache-key-suffix` - дополнительный суффикс для ключа кэша
- `skip-generate-networks` - пропустить генерацию сетей (по умолчанию: "false")
- `skip-download-deps` - пропустить загрузку зависимостей при попадании в кэш (по умолчанию: "false")
- `checkout-repo` - выполнить checkout репозитория (по умолчанию: "true")
- `fetch-depth` - глубина checkout (по умолчанию: "1")

**Выходы:**

- `cache-hit` - был ли кэш найден
- `go-cache-key` - сгенерированный ключ кэша
- `go-version` - фактическая версия Go

**Особенности:**

- ✅ Автоматическое кэширование с умными ключами
- ✅ Условная загрузка зависимостей (только при cache miss)
- ✅ Верификация модулей
- ✅ Детальное логирование
- ✅ Поддержка всех сценариев использования

## 🎯 Стратегия кэширования

### Ключи кэша

Формат: `{OS}-go-{GO_VERSION}-{HASH(go.sum)}-{SUFFIX}`

Примеры:

- `Linux-go-1.24.1-abc123-test`
- `macOS-go-1.23.x-def456-build-macos-latest`

### Fallback ключи

1. `{OS}-go-{GO_VERSION}-`
2. `{OS}-go-`

### Кэшируемые пути

- `~/.cache/go-build` - Go build cache
- `~/go/pkg/mod` - Go modules cache

### Логика кэширования

```
1. Генерация ключа кэша на основе go.sum + суффикса
2. Попытка восстановления кэша
3. При cache hit: пропуск загрузки зависимостей
4. При cache miss: загрузка и верификация зависимостей
5. Сохранение кэша для следующих запусков
```

## 📊 Оптимизации

### Производительность

- ✅ Shallow clone (`fetch-depth: 1`)
- ✅ Встроенное кэширование Go (`cache: true`)
- ✅ Многоуровневое кэширование модулей
- ✅ Умные ключи кэша с fallback
- ✅ Условная загрузка зависимостей

### Надежность

- ✅ Верификация модулей (`go mod verify`)
- ✅ Детальное логирование
- ✅ Статус кэша в выводе
- ✅ Информация о Go окружении

## 🔄 Workflow интеграция

### Простое использование

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: ./.github/actions/setup-go-project
        with:
          cache-key-suffix: "test"
      - run: make test
```

### С матрицей

```yaml
jobs:
  test:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest]
        go-version: ["1.24.1", "1.23.x"]
    steps:
      - uses: ./.github/actions/setup-go-project
        with:
          go-version: ${{ matrix.go-version }}
          cache-key-suffix: "test-${{ matrix.os }}"
      - run: make test
```

### Параллельные jobs

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: ./.github/actions/setup-go-project
        with:
          cache-key-suffix: "test"
      - run: make test

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: ./.github/actions/setup-go-project
        with:
          cache-key-suffix: "lint"
      - uses: golangci/golangci-lint-action@v8
        with:
          version: v2.1

  build:
    needs: [test, lint]
    runs-on: ubuntu-latest
    steps:
      - uses: ./.github/actions/setup-go-project
        with:
          cache-key-suffix: "build"
      - run: make build
```

## 🎯 Преимущества единого action

### ✅ Простота

- Один action для всех сценариев
- Нет необходимости в prepare/restore логике
- Автоматическое управление кэшем

### ✅ Производительность

- Кэш переиспользуется между jobs с одинаковыми зависимостями
- Умная логика загрузки зависимостей
- Оптимальные fallback стратегии

### ✅ Поддержка

- Единая точка изменений
- Консистентное поведение
- Простая отладка

## 🛠️ Поддержка

Action поддерживает:

- Множественные версии Go
- Кроссплатформенность (Linux, macOS, Windows)
- Детальное логирование
- Гибкую конфигурацию
- Оптимальное кэширование
- Все типы Go проектов
