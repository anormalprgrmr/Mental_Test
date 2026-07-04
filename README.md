# Mental Test Telegram Bot

A clean, scalable Go Telegram bot built with [`telebot`](https://gopkg.in/telebot.v3), SQLite, and data-driven test definitions.

## Features

- `/start` and `/tests` show available tests.
- Users can take exams through Telegram inline buttons.
- Results are calculated, saved in SQLite, and shown to the user.
- Admins receive a Telegram notification after each completed test.
- Initial tests:
  - Gardner Multiple Intelligences Test
  - Clifton Strengths Style Test
- Domain models/tests are separated from handlers and database code for easier scaling.

## Configuration

Copy `.env.example` to `.env` and set values, or export variables in your shell:

```env
BOT_TOKEN=123456:your-token
ADMINS=111111111,222222222
DB_PATH=mental_test.db
```

`BIT_TOKEN` is also accepted as a fallback alias if `BOT_TOKEN` is empty.

> `ADMINS` must contain numeric Telegram user IDs, not usernames.

## Run

```bash
go mod tidy
go run ./cmd/bot
```

If you keep variables in `.env`, load them before running, for example:

```bash
set -a
. ./.env
set +a
go run ./cmd/bot
```

## Project structure

```text
cmd/bot/                 Application entrypoint and dependency wiring
internal/config/         Environment configuration
internal/db/             SQLite connection, migrations, and queries
internal/handlers/       Telegram handlers and user interaction flow
internal/models/         Domain models, test definitions, catalog, and scoring engine
```

## Adding a new test

1. Create a new file in `internal/models`, for example `disc.go`.
2. Return a `TestDefinition` with dimensions and questions.
3. Register it in `DefaultCatalog()`:

```go
func DefaultCatalog() (*Catalog, error) {
    return NewCatalog(GardnerTest(), CliftonTest(), DiscTest())
}
```

The bot UI, scoring, persistence, and admin notification flow will work automatically for the new test.
