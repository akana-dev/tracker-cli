# Tracker CLI

Консольный трекер времени задач на Go с поддержкой нескольких серверов, комментариев, тегов, шаблонов, bulk-операций и автообновления.

## Возможности

- **Множественные серверы** — работа с несколькими трекерами одновременно
- **Управление задачами** — создание, редактирование, закрытие, пауза, возобновление
- **Кастомные времена** — указание произвольного времени для pause/resume
- **Комментарии** — Markdown, @mentions, watch через polling
- **Серверные теги** — классификация задач с цветовой маркировкой (True Color)
- **Серверные шаблоны** — предзаполненные наборы для быстрого создания задач
- **Bulk-операции** — массовое закрытие, назначение, удаление задач
- **Поиск** — по названию, комментариям задач и содержимому комментариев
- **Экспорт** — CSV, JSON, XLSX с пресетами и относительными датами
- **Автообновление** — через GitHub Releases с проверкой checksums
- **Алиасы** — короткие команды для частых операций
- **Относительные даты** — `today`, `last 7 days`, `this month`, `last monday`

## Установка

### Через curl/wget

```bash
# Linux/macOS
curl -fsSL https://github.com/akana-dev/tracker-cli/releases/latest/download/tracker-$(uname -s)-$(uname -m) -o tracker
chmod +x tracker
sudo mv tracker /usr/local/bin/

# Проверка версии
tracker --version
```

### Через Go

```bash
go install github.com/akana-dev/tracker-cli/cmd/tracker@latest
```

### Автообновление

```bash
tracker update              # Обновить до последней версии
tracker update --check      # Только проверить наличие обновления
tracker update --pre-release # Включая pre-release версии
```

## Быстрый старт

```bash
# 1. Настройка сервера
tracker configure
tracker server add work https://tracker.example.com

# 2. Авторизация
tracker login --username ivanov

# 3. Создание задачи
tracker task add "Исправить баг в авторизации" --company ACME --tag bug --tag urgent

# 4. Просмотр задач
tracker task list --today
tracker task list --tag bug

# 5. Детали задачи
tracker task view ACME-15

# 6. Пауза/возобновление с кастомным временем
tracker pause ACME-15 --at "14:30"
tracker resume ACME-15 --start "15:00"

# 7. Закрытие задачи
tracker task close ACME-15 --solution "Исправлено"
```

## Конфигурация

Файлы хранятся в `~/.tracker/`:

| Файл | Описание |
|------|----------|
| `servers.json` | Серверы и токены авторизации |
| `aliases.json` | Алиасы команд |
| `update-check.json` | Кэш проверок обновлений |
| `export-presets.yaml` | Пресеты экспорта |

### Управление серверами

```bash
tracker server add work https://tracker.example.com
tracker server list
tracker server use work
tracker server remove work
```

### Дефолтная компания

```bash
tracker config default-company ACME
tracker config show
```

## Команды

### Задачи

```bash
# Создание
tracker task add "Название задачи" [флаги]
  -s, --start string      Начало (по умолчанию: now)
  -e, --end string        Конец
  -q, --company string    Компания (по умолчанию — из конфига)
  -a, --assignee string   Исполнитель
  -S, --solution string   Статус
  -C, --comment string    Комментарий
  -T, --tag strings       Теги (можно указать несколько)

# Из шаблона
tracker task from <имя_шаблона> [флаги]
  -t, --title string      Переопределить название
  -s, --start string      Начало
  -q, --company string    Переопределить компанию
  -T, --tag strings       Дополнительные теги

# Список
tracker task list [флаги]
  -t, --today             Только сегодня
  -w, --week              За неделю
  -m, --month             За месяц
  -q, --company string    Фильтр по компании
  -S, --solution string   Фильтр по статусу
  -a, --assignee string   Фильтр по исполнителю
  -s, --search string     Поиск
  -C, --search-comments   Искать также в комментариях
  -T, --tag strings       Фильтр по тегам
  -A, --all               Показать все задачи
  -p, --page int          Номер страницы
  -l, --limit int         Количество задач на странице
  -o, --offset int        Смещение от начала

# Просмотр
tracker task view <тикет>
  -N, --no-comments       Не показывать комментарии

# Редактирование
tracker task edit <тикет> [флаги]
  -t, --title string      Новое название
  -s, --start string      Новое время начала
  -e, --end string        Новое время окончания
  -q, --company string    Новая компания
  -a, --assignee string   Новый исполнитель
  -S, --solution string   Новый статус
  -C, --comment string    Новый комментарий
  -T, --tag strings       Новые теги (полная замена)

# Жизненный цикл
tracker task close <тикет> [-s, --solution string]
tracker task pause <тикет> [-t, --at string]
tracker task resume <тикет> [-s, --start string]
tracker task assign <тикет> <исполнитель>
tracker task delete <тикет>

# Bulk-операции
tracker task bulk close <тикет1> <тикет2> ...
tracker task bulk assign <исполнитель> <тикет1> <тикет2> ...
tracker task bulk delete <тикет1> <тикет2> ... [-f, --force]

# Экспорт
tracker task export [флаги]
  -f, --format string     Формат: csv, json, xlsx
  -o, --output string     Имя выходного файла
  -p, --period string     Относительный период
  -D, --date-from string  Дата начала
  -T, --date-to string    Дата конца
  -i, --interactive       Интерактивный режим
  -r, --preset string     Имя пресета
```

### Комментарии

```bash
tracker task comment list <тикет> [флаги]
  -a, --all               Показать все комментарии
  -p, --page int          Номер страницы
  -l, --limit int         Количество на странице

tracker task comment add <тикет> [флаги]
  --editor                Открыть в редакторе
  --file string           Прочитать из файла
  -i, --interactive       Интерактивный режим

tracker task comment edit <id> <новый текст>
tracker task comment delete <id>
tracker task comment watch <тикет>   # Polling новых комментариев
```

### Теги

```bash
tracker tag add <имя> [-c, --color string]   # Цвет опционален (#RRGGBB)
tracker tag list [-s, --search string]
tracker tag update <id> [-n, --name string] [-c, --color string]
tracker tag delete <id> [-f, --force]
```

Примеры:
```bash
tracker tag add bug --color "#FF5733"
tracker tag add urgent
tracker tag add golang -c "#00add8"
tracker task add "Исправить баг" --tag bug --tag urgent
tracker task list --tag bug --tag urgent
```

### Шаблоны

```bash
tracker template add <имя> [флаги]
  -t, --title string      Заголовок задачи (обязательный)
  -d, --description string Описание
  -c, --company string    Компания
  -s, --solution string   Статус по умолчанию
  -p, --public            Публичный шаблон

tracker template list [-a, --all]           # --all только для admin
tracker template show <id>
tracker template use <id>                   # Создать задачу из шаблона
tracker template delete <id> [-f, --force]
```

### Авторизация

```bash
tracker login [-u, --username string] [-m, --method string]
  -u, --username string   Имя пользователя (пропустить ввод)
  -m, --method string     Метод: password, ad

tracker logout
tracker me
tracker register
```

### Компании

```bash
tracker company add <название> [-d, --description string]
tracker company list [-a, --all] [-p, --page int] [-l, --limit int]
tracker company delete <название>
```

### Алиасы

```bash
tracker alias add <имя> <команда>
tracker alias list
tracker alias remove <имя>
```

### Пресеты экспорта

```bash
tracker export preset save <имя> [флаги]
tracker export preset list
tracker export preset show <имя>
tracker export preset remove <имя>
```

### Обновление

```bash
tracker update                  # Обновить до последней версии
tracker update --check          # Только проверить
tracker update --pre-release    # Включая pre-release
```

## Примеры использования

### Работа с тегами

```bash
# Создание тегов с цветами
tracker tag add bug --color "#FF5733"
tracker tag add urgent -c "#C70039"
tracker tag add golang -c "#00add8"

# Привязка тегов к задаче
tracker task add "Оптимизация запроса" --tag golang --tag urgent

# Фильтрация по тегам (OR-логика)
tracker task list --tag bug --tag urgent

# Обновление тегов задачи (полная замена)
tracker task edit ACME-15 --tag critical-bug --tag urgent

# Очистка всех тегов
tracker task edit ACME-15 --tag
```

### Bulk-операции

```bash
# Массовое закрытие
tracker task bulk close ACME-15 ACME-16 ACME-17

# Массовое назначение
tracker task bulk assign ivanov ACME-20 ACME-21 ACME-22

# Массовое удаление с подтверждением
tracker task bulk delete ACME-99 ACME-100
tracker task bulk delete --force ACME-99 ACME-100
```

### Поиск

```bash
# Поиск по названию и комментарию задачи
tracker task list --search "авторизации"

# Поиск также в комментариях
tracker task list --search "авторизации" --search-comments
tracker task list -s "авторизации" -C
```

### Кастомные времена

```bash
# Пауза с указанием времени
tracker pause ACME-15 --at "14:30"
tracker pause ACME-15 -t "2026-07-01T14:30:00"

# Возобновление с указанием времени
tracker resume ACME-15 --start "15:00"
tracker resume ACME-15 -s "2026-07-01T15:00:00"
```

### Экспорт

```bash
# Экспорт за сегодня в CSV
tracker task export --format csv --period today

# Экспорт за последнюю неделю в XLSX
tracker task export -f xlsx -p "last 7 days"

# Экспорт с пресетом
tracker export preset save weekly --format xlsx --period "this week"
tracker task export --preset weekly

# Интерактивный режим
tracker task export --interactive
```

## Разработка

### Сборка

```bash
make build              # Локальная сборка
make build-all          # Кроссплатформенная + checksums
```

### Тесты

```bash
make test               # go test ./...
```

### Релиз

```bash
git tag -a v1.5.0 -m "Release v1.5.0"
git push origin v1.5.0  # Запустит GitHub Actions
```

### Локальная проверка обновления

```bash
./build/tracker update --check
./build/tracker update --pre-release
```

### Отключение проверки обновлений

```bash
export TRACKER_NO_UPDATE_CHECK=1
```

## Зависимости

- `github.com/fatih/color` — цветовая разметка
- `github.com/jedib0t/go-pretty/v6` — таблицы
- `github.com/spf13/cobra` — CLI-фреймворк
- `github.com/charmbracelet/glamour` — рендеринг Markdown
- `golang.org/x/term` — работа с терминалом
- `gopkg.in/yaml.v3` — YAML-парсинг

## Лицензия

MIT
