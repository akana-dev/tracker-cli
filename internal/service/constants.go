package service

import "time"

const (
	// WarningHoursWorked — порог, после которого сессия подсвечивается жёлтым.
	WarningHoursWorked = 4.0

	// ErrorHoursWorked — порог, после которого сессия подсвечивается красным.
	ErrorHoursWorked = 8.0
)

// DefaultPageSize — дефолтный размер страницы для list-команд.
const DefaultPageSize = 20

// WrapWidth — ширина для переноса текста в PrintIndented.
const WrapWidth = 70

const (
	// HTTPTimeout — таймаут на весь HTTP-запрос (включая retry).
	HTTPTimeout = 30 * time.Second

	// MaxRetries — максимальное количество повторных попыток при временных ошибках.
	MaxRetries = 3

	// InitialBackoff — начальная задержка между retry.
	InitialBackoff = 100 * time.Millisecond

	// MaxBackoff — максимальная задержка между retry.
	MaxBackoff = 2 * time.Second
)

const (
	// MaxTitleLength — максимальная длина названия задачи.
	MaxTitleLength = 200

	// MaxCommentLength — максимальная длина комментария.
	MaxCommentLength = 5000

	// MaxSolutionLength — максимальная длина решения.
	MaxSolutionLength = 2000

	// MaxUsernameLength — максимальная длина имени пользователя.
	MaxUsernameLength = 50

	// MaxEmailLength — максимальная длина email.
	MaxEmailLength = 254

	// MaxCompanyNameLength — максимальная длина названия компании.
	MaxCompanyNameLength = 100

	// MaxCompanyDescriptionLength — максимальная длина описания компании.
	MaxCompanyDescriptionLength = 500
)
