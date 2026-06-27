# Tracker CLI

Консольный трекер времени задач с поддержкой нескольких серверов, комментариев, тегов и шаблонов.

[![Release](https://img.shields.io/github/v/release/akana-dev/tracker-cli)](https://github.com/akana-dev/tracker-cli/releases)
[![Platforms](https://img.shields.io/badge/platforms-linux%20%7C%20macos%20%7C%20windows-blue)](https://github.com/akana-dev/tracker-cli/releases)
[![License](https://img.shields.io/github/license/akana-dev/tracker-cli)](LICENSE)

## Возможности

- **Управление задачами** — создание, редактирование, пауза, возобновление, закрытие
- **Комментарии** — с поддержкой Markdown и упоминаниями `@username`
- **Теги** — гибкая классификация задач
- **Шаблоны** — быстрое создание типовых задач из YAML-файлов
- **Экспорт** — в CSV, JSON, XLSX с фильтрацией и пресетами
- **Несколько серверов** — переключение между разными API
- **Авторизация** — по паролю или через Active Directory
- **Алиасы** — короткие команды для частых действий
- **Автообновление** — проверка и установка новых версий одной командой
- **Часовые пояса** — корректная работа с DST и локальным временем

## Установка

### Через GitHub Releases (рекомендуется)

Скачайте бинарник для вашей платформы со страницы [Releases](https://github.com/akana-dev/tracker-cli/releases):

```bash
# Linux amd64
curl -L https://github.com/akana-dev/tracker-cli/releases/latest/download/tracker-linux-amd64 -o /usr/local/bin/tracker
chmod +x /usr/local/bin/tracker

# macOS arm64 (Apple Silicon)
curl -L https://github.com/akana-dev/tracker-cli/releases/latest/download/tracker-darwin-arm64 -o /usr/local/bin/tracker
chmod +x /usr/local/bin/tracker

# Windows
# Скачайте tracker-windows-amd64.exe и добавьте в PATH
```

### Автоматическое обновление

Если tracker уже установлен:

```bash
tracker update
```

### Сборка из исходников

Требования: Go 1.26+

```bash
git clone https://github.com/akana-dev/tracker-cli.git
cd tracker-cli
make build
sudo make install
```

## Быстрый старт

```bash
# 1. Настройка подключения к серверу
tracker configure

# 2. Авторизация
tracker login

# 3. Создание задачи
tracker task add "Исправить баг в авторизации" --company ACME

# 4. Просмотр списка задач
tracker task list --today

# 5. Добавление комментария
tracker task comment add ACME-1 "Нашёл причину, исправляю"

# 6. Закрытие задачи
tracker task close ACME-1 --solution "Исправлено"
```

## Основные команды

### Авторизация

| Команда | Описание |
|---------|----------|
| `tracker login` | Вход в систему (password или AD) |
| `tracker logout` | Выход из системы |
| `tracker me` | Информация о текущем пользователе |
| `tracker register` | Регистрация нового пользователя |

### Задачи

| Команда | Описание |
|---------|----------|
| `tracker task add [название]` | Создать задачу |
| `tracker task list` | Список задач с фильтрацией |
| `tracker task view [тикет]` | Подробная информация о задаче |
| `tracker task edit [тикет]` | Редактировать задачу |
| `tracker task close [тикет]` | Закрыть задачу |
| `tracker task pause [тикет]` | Поставить на паузу |
| `tracker task resume [тикет]` | Возобновить задачу |
| `tracker task assign [тикет] [user]` | Назначить исполнителя |
| `tracker task delete [тикет]` | Удалить задачу |

**Примеры:**

```bash
# Создать задачу с указанием компании и исполнителя
tracker task add "Рефакторинг модуля" --company ACME --assignee ivanov

# Список задач за сегодня
tracker task list --today

# Список задач за неделю с фильтром по компании
tracker task list --week --company ACME

# Пагинация
tracker task list --page 2 --limit 20

# Фильтрация по тегам
tracker task list --tag backend --tag urgent

# Показать все задачи без пагинации
tracker task list --all
```

### Комментарии

| Команда | Описание |
|---------|----------|
| `tracker task comment list [тикет]` | Список комментариев |
| `tracker task comment add [тикет] [текст]` | Добавить комментарий |
| `tracker task comment edit [тикет] [id]` | Редактировать комментарий |
| `tracker task comment delete [тикет] [id]` | Удалить комментарий |
| `tracker task comment watch [тикет]` | Следить за новыми комментариями |

**Примеры:**

```bash
# Короткий комментарий
tracker task comment add ACME-1 "Задача выполнена"

# Через редактор (для длинных комментариев с Markdown)
tracker task comment add ACME-1 --editor

# Из файла
tracker task comment add ACME-1 --file report.md

# Интерактивный режим
tracker task comment add ACME-1 --interactive

# Слежение за новыми комментариями (Ctrl+C для остановки)
tracker task comment watch ACME-1 --interval 10s
```

**Поддерживаемый Markdown:**

- `**жирный**`, `*курсив*`, `` `код` ``
- Списки, заголовки, ссылки
- Блоки кода с подсветкой языка
- Упоминания `@username`

### Теги

| Команда | Описание |
|---------|----------|
| `tracker tag add [тикет] [теги...]` | Добавить теги к задаче |
| `tracker tag remove [тикет] [теги...]` | Удалить теги |
| `tracker tag list` | Список всех тегов |
| `tracker tag show [тикет]` | Теги конкретной задачи |

**Примеры:**

```bash
tracker tag add ACME-1 backend refactoring urgent
tracker tag remove ACME-1 urgent
tracker tag list
tracker task list --tag backend
```

### Шаблоны

| Команда | Описание |
|---------|----------|
| `tracker template add [имя]` | Создать шаблон (откроется редактор) |
| `tracker template list` | Список шаблонов |
| `tracker template show [имя]` | Содержимое шаблона |
| `tracker template edit [имя]` | Редактировать шаблон |
| `tracker template remove [имя]` | Удалить шаблон |
| `tracker task from [шаблон]` | Создать задачу из шаблона |

**Пример YAML-шаблона** (`~/.tracker/templates/daily.yaml`):

```yaml
name: daily-standup
title: "Ежедневная планёрка"
company: ACME
comment: "Регулярная встреча команды"
tags:
  - meeting
  - daily
```

Использование:

```bash
tracker task from daily-standup
tracker task from daily-standup --title "Планёрка по проекту X"
```

### Экспорт

| Команда | Описание |
|---------|----------|
| `tracker task export` | Экспорт задач в файл |
| `tracker export preset save [имя]` | Сохранить пресет экспорта |
| `tracker export preset list` | Список пресетов |
| `tracker export preset remove [имя]` | Удалить пресет |

**Примеры:**

```bash
# Базовый экспорт
tracker task export --format csv --output tasks.csv

# С пресетом
tracker task export --preset monthly

# С относительным периодом
tracker task export --format xlsx --period "last month"

# Preview перед экспортом
tracker task export --format csv --preview --today

# Интерактивный режим
tracker task export --interactive

# С выбором колонок
tracker task export --format csv --fields ticket,title,hours --today

# Сложный фильтр
tracker task export --format xlsx \
  --period "last 7 days" \
  --company ACME \
  --open-only \
  --all-users \
  --timezone Europe/Moscow
```

**Поддерживаемые периоды:**

- `today`, `yesterday`, `tomorrow`
- `this week`, `last week`, `next week`
- `this month`, `last month`, `next month`
- `this quarter`, `last quarter`
- `this year`, `last year`
- `last 7 days`, `last 30 days`
- `last monday`, `next friday`

### Алиасы

| Команда | Описание |
|---------|----------|
| `tracker alias add [имя] [команда]` | Создать алиас |
| `tracker alias list` | Список алиасов |
| `tracker alias remove [имя]` | Удалить алиас |

**Примеры:**

```bash
tracker alias add ll "task list --today"
tracker alias add w "task list --week"
tracker alias add st "task list"

# Использование
tracker ll
tracker w
```

### Компании

| Команда | Описание |
|---------|----------|
| `tracker company list` | Список компаний |
| `tracker company add [название]` | Добавить компанию (admin) |
| `tracker company delete [название]` | Удалить компанию (admin) |

### Конфигурация

| Команда | Описание |
|---------|----------|
| `tracker configure` | Настроить подключение к API |
| `tracker server add [имя] [url]` | Добавить сервер |
| `tracker server list` | Список серверов |
| `tracker server use [имя]` | Переключиться на сервер |
| `tracker server remove [имя]` | Удалить сервер |
| `tracker config default-company [название]` | Компания по умолчанию |
| `tracker config show` | Показать конфигурацию |

### Обновление

| Команда | Описание |
|---------|----------|
| `tracker update` | Проверить и установить обновление |
| `tracker update --check` | Только проверить наличие обновления |
| `tracker update --pre-release` | Включая pre-release версии |
| `tracker --version` | Показать текущую версию |

Уведомление о новой версии показывается автоматически при запуске (не чаще раза в сутки). Отключить:

```bash
export TRACKER_NO_UPDATE_CHECK=1
```

## Глобальные флаги

| Флаг | Описание |
|------|----------|
| `--help`, `-h` | Справка |
| `--version`, `-v` | Версия |
| `--install` | Установить в систему (добавить в PATH) |
| `--uninstall` | Удалить из системы |

## Структура файлов

```
~/.tracker/
├── servers.json          # Конфигурация серверов и токены
├── aliases.json          # Алиасы команд
├── tags.json             # Теги задач (локально, по серверам)
├── update-check.json     # Кэш проверки обновлений
── export-presets.yaml   # Пресеты экспорта
└── templates/            # Шаблоны задач (YAML)
    ├── daily.yaml
    └── meeting.yaml
```

## Разработка

### Сборка

```bash
# Локальная сборка
make build

# Кроссплатформенная сборка
make build-all

# Запуск тестов
make test

# Очистка
make clean
```

### Создание релиза

```bash
git tag -a v1.5.0 -m "Release v1.5.0"
git push origin v1.5.0
```

GitHub Actions автоматически соберёт бинарники для всех платформ и создаст релиз.

### Структура проекта

```
tracker-cli/
├── cmd/tracker/main.go        # Точка входа
├── internal/
│   ├── cli/                   # CLI-команды (cobra)
│   │   ├── task/              # Подкоманды task
│   │   │   ├── comment/       # Подкоманды comment
│   │   │   └── ...
│   │   ── ...
│   ├── client/                # HTTP-клиент для API
│   ├── config/                # Работа с конфигом
│   ├── models/                # Модели данных
│   ├── service/               # Бизнес-логика
│   ├── ui/                    # Цветной вывод в терминал
│   ├── updater/               # Автообновление
│   ├── version/               # Версия (встраивается при сборке)
│   ── installer/             # Установка в систему
├── pkg/                       # Переиспользуемые пакеты
│   ├── table/                 # Таблицы
│   └── timeparse/             # Парсинг времени
── Makefile
└── README.md
```

## Лицензия

MIT License. См. файл [LICENSE](LICENSE) для подробностей.